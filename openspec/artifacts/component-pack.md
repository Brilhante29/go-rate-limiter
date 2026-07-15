# Component Pack: go-rate-limiter

## Selected Pack

- Pack id: `backend-reliability-platform`
- Pack name: Backend Reliability and Architecture Platform
- Problem: Prove backend architecture, reliability, latency, consistency, and failure handling.

## Benchmark Focus

- p95_latency_ms
- p99_latency_ms
- throughput_rps
- consistency_rate
- event_throughput
- zero_loss_under_failure

## Preferred Artifacts

- ADRs
- API contract
- load test
- failure scenario
- consistency proof

## Rejection Rules

- Reject framework-first architecture with domain depending on infra.
- Reject RabbitMQ, Kafka, or microservices unless the benchmark needs async scale or fault isolation.
- Reject "clean architecture" claims without dependency direction tests.

## Reuse Priority

1. Use repo-local `.codex/skills/` and `.claude/skills/`.
2. Use `.portfolio/` and upstream `portfolio-reuse-kit`.
3. Use external repositories as references for organization, workflow, schemas, tests, benchmarks, and docs.
4. Use external code only with license compatibility, attribution, and a decision record.
