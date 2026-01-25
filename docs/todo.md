# ADR Documentation Enhancement Todo List

**Generated:** 2026-01-25
**Priority:** Low (Non-blocking)
**Status:** Backlog

## Overview

This document tracks minor documentation enhancements identified during the ADR synchronization review. All ADRs are fully implemented with 9/10 compliance scores. These items improve documentation clarity but are not blocking issues.

---

## ADR-001: Go Language Selection for MCP Server

**Current Compliance:** 9/10
**Gap:** Missing Kubernetes references in documentation

### Tasks
- [ ] Add explicit Kubernetes client-go references to Context section
  - Location: `docs/adrs/001-go-language-selection.md` - Context section
  - Detail: Reference `pkg/clients/kubernetes.go` as evidence of Go + Kubernetes integration
  - Effort: 15 minutes

- [ ] Document Go version compatibility with Kubernetes API versions
  - Location: `docs/adrs/001-go-language-selection.md` - Technical Considerations section
  - Detail: Add table showing Go 1.24+ compatibility with Kubernetes 1.24-1.26 API versions
  - Effort: 30 minutes

---

## ADR-005: Stateless Design (No Database)

**Current Compliance:** 9/10
**Gap:** Missing Kubernetes stateless pattern references

### Tasks
- [ ] Enhance Consequences section with Kubernetes StatefulSet comparison
  - Location: `docs/adrs/005-stateless-design.md` - Consequences section
  - Detail: Add explicit note about why Deployment is used instead of StatefulSet
  - Effort: 20 minutes

- [ ] Document caching strategy as stateless pattern implementation
  - Location: `docs/adrs/005-stateless-design.md` - Implementation section
  - Detail: Reference `pkg/cache/memory_cache.go` as example of stateless caching
  - Effort: 15 minutes

---

## ADR-007: RBAC-Based Security Model

**Current Compliance:** 9/10
**Gap:** Missing Kubernetes RBAC documentation cross-references

### Tasks
- [ ] Cross-reference Kubernetes RBAC documentation
  - Location: `docs/adrs/007-rbac-based-security-model.md` - Decision section
  - Detail: Add links to official Kubernetes RBAC docs and OpenShift RBAC best practices
  - Effort: 10 minutes

- [ ] Add examples of ClusterRole and ServiceAccount YAML from charts/
  - Location: `docs/adrs/007-rbac-based-security-model.md` - Implementation section
  - Detail: Include snippets from `charts/openshift-cluster-health-mcp/templates/clusterrole.yaml`
  - Effort: 25 minutes

---

## ADR-009: Architecture Evolution Roadmap

**Current Compliance:** 9/10
**Gap:** Missing PostgreSQL and Kubernetes future planning details

### Tasks
- [ ] Update Phase 3 PostgreSQL planning section with current timeline
  - Location: `docs/adrs/009-architecture-evolution-roadmap.md` - Roadmap section
  - Detail: Clarify Phase 3 timeline or mark as "deferred" if no active plans
  - Effort: 20 minutes

- [ ] Document decision criteria for when to implement persistent storage
  - Location: `docs/adrs/009-architecture-evolution-roadmap.md` - Decision Criteria section
  - Detail: Add specific triggers (e.g., "when incident history >1000 items" or "user request")
  - Effort: 30 minutes

---

## ADR-010: Version Compatibility and Upgrade Roadmap

**Current Compliance:** 9/10
**Gap:** Missing Kubernetes version compatibility matrix

### Tasks
- [ ] Add Kubernetes API version compatibility matrix
  - Location: `docs/adrs/010-version-compatibility-upgrade-roadmap.md` - Compatibility section
  - Detail: Create table showing MCP server versions vs. Kubernetes/OpenShift versions
  - Effort: 45 minutes

- [ ] Document tested Kubernetes versions (1.24, 1.25, 1.26, etc.)
  - Location: `docs/adrs/010-version-compatibility-upgrade-roadmap.md` - Testing section
  - Detail: List versions tested in CI/CD and production environments
  - Effort: 20 minutes

---

## Security Notice Review

**Current Status:** Informational (not a gap)
**Confidence:** 60% (likely false positive)

### Tasks
- [ ] Review GitHub workflow token usage in `.github/workflows/ci.yml:86`
  - Detail: Verify this is standard `secrets.GITHUB_TOKEN` reference
  - Expected: False positive from tree-sitter analysis (standard GitHub Actions pattern)
  - Effort: 5 minutes

- [ ] Review GitHub workflow token usage in `.github/workflows/ci.yml:98`
  - Detail: Verify this is standard `secrets.GITHUB_TOKEN` reference
  - Expected: False positive from tree-sitter analysis (standard GitHub Actions pattern)
  - Effort: 5 minutes

---

## Implementation Notes

### Process
1. Create GitHub Issues for each ADR enhancement (optional)
2. Tag issues with `documentation`, `adr`, and `enhancement` labels
3. Assign to appropriate team member during sprint planning
4. Target completion: Q2 2026 (no urgency)

### Acceptance Criteria
- ADR documents updated with new content
- Cross-references to code files include line numbers
- All links to external documentation are valid
- Changes reviewed by at least one team member
- ADR compliance score improves to 9.5+/10 after changes

### Estimated Total Effort
**Total:** ~4 hours (distributed across multiple team members)

---

## Related Documents
- [ADR Sync Report](./ADR_SYNC_REPORT.md)
- [ADR Directory](./adrs/)
- [Branch Protection](./BRANCH_PROTECTION.md)
- [CLAUDE.md](../CLAUDE.md)
