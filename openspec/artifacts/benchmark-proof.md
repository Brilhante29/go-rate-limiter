# Benchmark Proof: go-rate-limiter

## Primary Metric

- Metric: `total_rps`
- Unit: `requests_per_second`
- Result: total_rps = 5543.34 requests_per_second (median of three runs)
- Result path: `benchmarks/results/rate-limiter-baseline.json`

## Command

    pwsh -NoProfile -File tools/benchmark.ps1

## Evidence

Median p95 latency: 28.290 ms. Both nodes were observed with zero unexpected errors and zero global-limit violations.

The README/post number must come from the committed benchmark JSON, not from manual text.
