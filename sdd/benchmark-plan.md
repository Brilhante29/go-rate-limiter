# Benchmark Plan: go-rate-limiter

## Hypothesis

Under load greater than the configured quota, two HTTP nodes backed by one Redis bucket will:

- both participate in decisions;
- admit no more than the global token-bucket maximum;
- return accepted requests/s, rejected requests/s, total requests/s, and p95 latency;
- complete with zero unexpected HTTP or transport errors.

## Reproduction

```powershell
pwsh -NoProfile -File tools/benchmark.ps1
```

Minimal output-only path:

```powershell
docker build -t go-rate-limiter .
docker run --rm go-rate-limiter
```

## Controlled Inputs

| Input | Baseline |
|---|---:|
| nodes | 2 |
| store | Redis 8.8 |
| shared keys | 1 |
| cost per request | 1 token |
| rate | 1,000 tokens/s |
| burst | 1,000 tokens |
| duration | 5 s per repetition |
| concurrency | 64 workers |
| benchmark repetitions | 3 after one warm-up |

## Metrics

| Metric | Unit | Direction | Source |
|---|---|---|---|
| total_rps | requests/s | higher; primary throughput metric | all completed decisions / elapsed |
| accepted_rps | requests/s | constrained by policy; correctness companion | HTTP 200 count / elapsed |
| rejected_rps | requests/s | evidence of overload | HTTP 429 count / elapsed |
| p95_latency_ms | ms | lower | sorted end-to-end request samples |
| nodes_observed | count | exactly 2 | `X-Rate-Limiter-Node` headers |
| global_limit_preserved | boolean | true | accepted <= ceil(burst + rate * elapsed) + 1 |
| errors | count | 0 | transport and unexpected statuses |

## Methodology Boundaries

- The baseline runs all processes on one Docker host to make reproduction cheap and deterministic.
- It proves shared state and atomic coordination, not network partition tolerance or multi-region consistency.
- Redis server time is used in the script so node clock skew does not affect refill calculations.
- p95 includes HTTP, routing, client pooling, Redis round-trip, script execution, and JSON response time.

## Result

File: `benchmarks/results/rate-limiter-baseline.json`.

The README number must be copied from that file after reproduction; manually invented numbers are not accepted.
## Measured Baseline

| Metric | Result |
|---|---:|
| accepted_rps | 1,197.66 median |
| rejected_rps | 4,347.93 median |
| total_rps | 5,543.34 median |
| p95_latency_ms | 28.290 median |
| nodes_observed | 2 |
| errors | 0 |
| global_limit_violations | 0 |

The three measured throughput samples were 5,404.27, 6,473.98, and 5,543.34 req/s. Every repetition observed both HTTP nodes and recorded zero unexpected errors and zero global-limit violations.

Measured on 2026-08-14 with Docker Linux/amd64, 6 CPUs, Go 1.26.5, and Redis 8.8.0.
