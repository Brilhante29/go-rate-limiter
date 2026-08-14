# Agent Handoff

Project: `12 - go-rate-limiter`

## Principal Agent Summary

- Objective: prove atomic global rate limiting across two application nodes.
- Portfolio program: Backend Reliability and Architecture Platform.
- Public proof: accepted/rejected/total requests per second, p95 latency, two observed nodes, and preserved global quota.
- Default runnable path: `docker build -t go-rate-limiter .` then `docker run --rm go-rate-limiter`.

## Subagent Decisions

| Role | Decision | Evidence | Status |
|---|---|---|---|
| program-planner | Shared traffic-control primitive for later gateway and multi-tenant projects | `project.yaml` | complete |
| architecture-selector | Narrow Hexagonal boundary around `BucketStore` | `sdd/architecture-decision.md` | complete |
| engineering-principles-reviewer | SOLID plus explicit KISS/YAGNI rejections | architecture and technical decisions | complete |
| stack-decision-agent | Go 1.26, chi, go-redis, Redis Lua, Docker | `sdd/technical-decision.md` | complete |
| api-style-agent | REST/HTTP; reject GraphQL and async messaging | `api/openapi.yaml` | complete |
| cloud-local-first-agent | Docker-local Redis; cloud mode none; Kumo not applicable | technical decision | complete |
| messaging-agent | No broker on synchronous decision path | technical decision | complete |
| language-profile-agent | `cmd`, `internal`, colocated tests, race/vet/bench gates | repository layout and CI | complete |
| benchmark-harness-agent | Two-node correctness-aware Go runner plus k6 profile | `internal/loadbench`, `benchmarks/k6.js` | complete |
| security-reuse-reviewer | bounded input, fail closed, no secrets, attributed references | API, code, `REFERENCES.md` | complete |
| release-ci-publisher | Final V2 evidence is ready; exact-head GitHub Actions remains the publication gate | workflow and release checklist | pending remote CI |

## Architecture Boundaries

- Core: policy, decision, service, and `BucketStore` port.
- Input adapter: chi HTTP handler.
- Output adapters: Redis and memory.
- Operations: server lifecycle, self-contained demo, Compose, and benchmark.
- Dependency direction: adapters import core; core imports no adapter or framework.

## Benchmark Handoff

- Primary metric: total_rps.
- Companion metrics: accepted_rps, rejected_rps, p95_latency_ms, nodes_observed, unexpected_errors, and global_limit_violations.
- Correctness gate: `nodes_observed=2`, `unexpected_errors=0`, and `global_limit_violations=0` in every repetition.
- Result: `benchmarks/results/rate-limiter-baseline.json`.
- Fixture: one shared key, two nodes, one Redis instance, cost one.

## Open Risks

- Localhost benchmark numbers vary by host; compare only with recorded environment and inputs.
- Redis Cluster, failover, and partitions are explicitly untested.
- The demo co-locates processes for reproducibility; Compose is the separated topology.

## Publication Gates

- [x] Docker proof works from the current tree.
- [x] Benchmark JSON and README numbers agree.
- [x] Architecture, API, references, and reuse review are complete.
- [x] No secret or paid service is required.
- [x] Project validation passes.
- [ ] GitHub CI passes on the exact final published commit.
