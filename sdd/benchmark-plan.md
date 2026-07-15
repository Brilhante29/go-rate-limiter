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
| duration | 10 s |
| concurrency | 64 workers |
| benchmark repetitions | 1 committed baseline; rerun for comparison |

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
| accepted_rps | 1,083.62 |
| rejected_rps | 10,427.40 |
| total_rps | 11,511.03 |
| p95_latency_ms | 11.985 |
| nodes_observed | 2 |
| errors | 0 |
| global_limit_preserved | true |

The committed run completed 115,129 decisions in 10.002 seconds on Docker Linux/amd64 with 16 CPUs, Go 1.26.5, and Redis 8.8.0. A second non-writing confirmation reached 11,146.16 req/s at 12.096 ms p95, a 3.2% throughput difference. The independent k6 topology also passed all thresholds at 6,897.03 HTTP req/s and 8.44 ms p95.

Measured on 2026-07-15.
