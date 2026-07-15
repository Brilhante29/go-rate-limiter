# Intent: go-rate-limiter

## Measurable Claim

Two HTTP nodes enforce one global token-bucket quota through atomic Redis state.

## Problem

Provides the shared traffic-control primitive used by gateways and multi-tenant services in the backend reliability program.

## In Scope

- Use the selected component pack: `backend-reliability-platform`.
- Keep the project under the Backend Reliability and Architecture Platform program.
- Preserve the benchmark contract: `total_rps` in `benchmarks/results/rate-limiter-baseline.json`.
- Keep the default path local-first and reproducible.

## Out Of Scope

- Paid credentials for the default demo.
- External infrastructure that is not required by the benchmark.
- Replacing local portfolio skills with external components silently.

## Default Demo Path

- Status: benchmarked
- Runtime: Go 1.26 static binary plus Redis 8.8
- Benchmark command: `pwsh -NoProfile -File tools/benchmark.ps1`

## Public Proof

- Benchmark: total_rps = 11511.03 requests_per_second
- Result path: `benchmarks/results/rate-limiter-baseline.json`
