// Package main provides a lightweight health check binary for init containers.
// This replaces curl-based health checks since the container image doesn't include curl.
//
// Usage:
//
//	healthcheck <url> [--timeout=<seconds>] [--interval=<seconds>] [--max-retries=<n>]
//
// Examples:
//
//	healthcheck http://coordination-engine:8080/health
//	healthcheck http://prometheus:9090/-/ready --timeout=10 --interval=5
//	healthcheck http://service:8080/healthz --max-retries=30
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	// Parse flags
	timeout := flag.Duration("timeout", 5*time.Second, "HTTP request timeout")
	interval := flag.Duration("interval", 10*time.Second, "Retry interval between health checks")
	maxRetries := flag.Int("max-retries", 0, "Maximum number of retries (0 = unlimited)")
	flag.Parse()

	// Validate args
	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Usage: healthcheck <url> [--timeout=<duration>] [--interval=<duration>] [--max-retries=<n>]")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Examples:")
		fmt.Fprintln(os.Stderr, "  healthcheck http://coordination-engine:8080/health")
		fmt.Fprintln(os.Stderr, "  healthcheck http://prometheus:9090/-/ready --timeout=10s --interval=5s")
		fmt.Fprintln(os.Stderr, "  healthcheck http://service:8080/healthz --max-retries=30")
		os.Exit(1)
	}

	url := args[0]
	client := &http.Client{Timeout: *timeout}

	retries := 0
	for {
		resp, err := client.Get(url)
		if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			_ = resp.Body.Close()
			fmt.Printf("Service at %s is ready! (HTTP %d)\n", url, resp.StatusCode)
			os.Exit(0)
		}

		// Log the failure reason
		if err != nil {
			fmt.Printf("Service at %s not ready: %v\n", url, err)
		} else {
			fmt.Printf("Service at %s not ready: HTTP %d\n", url, resp.StatusCode)
			_ = resp.Body.Close()
		}

		retries++
		if *maxRetries > 0 && retries >= *maxRetries {
			fmt.Fprintf(os.Stderr, "Max retries (%d) exceeded, giving up\n", *maxRetries)
			os.Exit(1)
		}

		fmt.Printf("Retrying in %v... (attempt %d", *interval, retries)
		if *maxRetries > 0 {
			fmt.Printf("/%d", *maxRetries)
		}
		fmt.Println(")")
		time.Sleep(*interval)
	}
}
