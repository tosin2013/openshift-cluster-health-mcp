# GitHub Issue: JSON Field Mismatch Causes anomaly_score to Always Be 0

## Issue Title

`Bug: JSON field mismatch causes anomaly_score to always be 0 when parsed from Coordination Engine`

## Labels

- `bug`
- `high-priority`
- `analyze-anomalies`

---

## Description

The `analyze-anomalies` tool always returns `anomaly_score: 0` and `confidence: 0` even when the Coordination Engine detects critical anomalies with `anomaly_score: 1, confidence: 0.87`.

This causes OpenShift Lightspeed to report "No anomalies detected" or "info severity" when critical anomalies actually exist.

## Root Cause

JSON field name mismatch in `pkg/clients/coordination_engine.go`:

**Current Code (line 127-136):**
```go
type AnomalyPattern struct {
    Metric      string  `json:"metric"`
    Type        string  `json:"type"`
    Severity    string  `json:"severity"`
    Score       float64 `json:"score"`      // ❌ Expects "score"
    Timestamp   string  `json:"timestamp"`
    Value       float64 `json:"value"`
    ExpectedMin float64 `json:"expected_min"`
    ExpectedMax float64 `json:"expected_max"`
    Model       string  `json:"model"`
}
```

**But Coordination Engine returns:**
```json
{
  "anomalies": [{
    "timestamp": "2026-01-28T18:01:31Z",
    "severity": "critical",
    "anomaly_score": 1,           // ✅ Uses "anomaly_score" not "score"
    "confidence": 0.87,           // ✅ Has "confidence" field (not mapped)
    "metrics": {
      "container_restart_count": 9,
      "node_cpu_utilization": 0.21
    },
    "explanation": "Container restarts detected (9)",
    "recommended_action": "restart_pod"
  }]
}
```

## Steps to Reproduce

### 1. Direct API Call (Working)
```bash
oc exec -n self-healing-platform deployment/coordination-engine -- curl -s -X POST \
  http://localhost:8080/api/v1/anomalies/analyze \
  -H "Content-Type: application/json" \
  -d '{"scope": "cluster"}' | jq '{anomalies_detected, anomalies}'
```

**Result:** 
```json
{
  "anomalies_detected": 1,
  "anomalies": [{
    "severity": "critical",
    "anomaly_score": 1,
    "confidence": 0.87
  }]
}
```

### 2. Via MCP Server / Lightspeed (Broken)

Ask Lightspeed: "Check for anomalies cluster-wide"

**Result from Lightspeed:**
```
Metrics checked: cpu_usage, memory_usage, pod_restarts
Result: 1 anomaly reported for each metric, but all have 
anomaly_score = 0, confidence = 0 and severity = "info"
```

## Expected Behavior

MCP server should correctly parse and return the anomaly scores from Coordination Engine:
- `anomaly_score: 1` (not 0)
- `confidence: 0.87` (not 0)  
- `severity: critical` (not info)

## Proposed Fix

Update `pkg/clients/coordination_engine.go`:

```go
type AnomalyPattern struct {
    Metric            string             `json:"metric"`
    Type              string             `json:"type"`
    Severity          string             `json:"severity"`
    Score             float64            `json:"anomaly_score"`    // FIX: was "score"
    Confidence        float64            `json:"confidence"`       // ADD: missing field
    Timestamp         string             `json:"timestamp"`
    Value             float64            `json:"value"`
    ExpectedMin       float64            `json:"expected_min"`
    ExpectedMax       float64            `json:"expected_max"`
    Model             string             `json:"model"`
    Metrics           map[string]float64 `json:"metrics"`          // ADD: for detailed metrics
    Explanation       string             `json:"explanation"`      // ADD: human-readable explanation
    RecommendedAction string             `json:"recommended_action"` // ADD: suggested action
}
```

Also update `internal/tools/analyze_anomalies.go` line 214 to use `Confidence` instead of `Score`:

```go
anomaly := AnomalyResult{
    Timestamp:    pattern.Timestamp,
    MetricName:   pattern.Metric,
    Value:        pattern.Value,
    AnomalyScore: pattern.Score,
    Confidence:   pattern.Confidence,  // FIX: was pattern.Score
    Severity:     pattern.Severity,
    Explanation:  generateExplanation(pattern.Metric, pattern.Score, pattern.Confidence),
}
```

## Environment

| Component | Version |
|-----------|---------|
| MCP Server | `quay.io/takinosh/openshift-cluster-health-mcp:4.18-latest` |
| Coordination Engine | `ocp-4.18-814cb25` |
| OpenShift | 4.18.21 |
| OpenShift Lightspeed | 1.0 |

## Impact

- **High**: Users relying on Lightspeed for anomaly detection will miss critical alerts
- **Workaround**: Use direct API calls to Coordination Engine (`/api/v1/anomalies/analyze`)

## Additional Context

### Coordination Engine Response Structure

The full response from `/api/v1/anomalies/analyze`:

```json
{
  "status": "success",
  "time_range": "1h",
  "scope": {
    "target_description": "cluster-wide"
  },
  "model_used": "anomaly-detector",
  "anomalies_detected": 1,
  "anomalies": [
    {
      "timestamp": "2026-01-28T18:01:31Z",
      "severity": "critical",
      "anomaly_score": 1,
      "confidence": 0.87,
      "metrics": {
        "container_restart_count": 9,
        "node_cpu_utilization": 0.067,
        "node_memory_utilization": 0.184,
        "pod_cpu_usage": 0.0006,
        "pod_memory_usage": 0.337
      },
      "explanation": "Container restarts detected (9)",
      "recommended_action": "restart_pod"
    }
  ],
  "summary": {
    "max_score": 1,
    "average_score": 1,
    "metrics_analyzed": 5,
    "features_generated": 45
  },
  "recommendation": "CRITICAL: Immediate investigation recommended."
}
```

### MCP Server Code Path

1. User asks Lightspeed about anomalies
2. Lightspeed calls MCP tool `analyze-anomalies`
3. `internal/tools/analyze_anomalies.go:Execute()` calls `coordinationEngine.AnalyzeAnomalies()`
4. `pkg/clients/coordination_engine.go:AnalyzeAnomalies()` makes HTTP POST to `/api/v1/anomalies/analyze`
5. Response is unmarshaled into `AnalyzeAnomaliesResponse` with `Patterns []AnomalyPattern`
6. **BUG**: `AnomalyPattern.Score` expects `json:"score"` but API returns `"anomaly_score"`
7. Go unmarshals `anomaly_score: 1` → `Score: 0` (field name mismatch, uses zero value)
8. MCP returns `anomaly_score: 0` to Lightspeed

---

**Discovered:** 2026-01-28
**Reporter:** AI Ops Platform Team
