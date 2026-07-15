# Benchmark Proof: go-rate-limiter

## Primary Metric

- Metric: `total_rps`
- Unit: `requests_per_second`
- Result: total_rps = 11511.03 requests_per_second
- Result path: `benchmarks/results/rate-limiter-baseline.json`

## Command

    pwsh -NoProfile -File tools/benchmark.ps1

## Evidence

P95 latency: 11.98 ms.

The README/post number must come from the committed benchmark JSON, not from manual text.
