---
name: benchmark-harness
description: Add, run, or validate reproducible benchmarks for portfolio repositories. Use when the user asks for benchmark scripts, k6 tests, latency/throughput/recall/cost metrics, benchmark JSON, p95/p99 comparison, or README benchmark tables.
---

# Benchmark Harness

Use benchmarks as product evidence.

1. Prefer one primary metric and at most two secondary metrics.
2. Make the benchmark deterministic: fixed seed, fixed fixture, documented environment.
3. Write machine-readable JSON to `benchmarks/results/`.
4. Put the headline metric in the first screen of README.
5. Use `harness/bench.py` for generic command latency when no domain-specific runner exists.
6. Use k6 for HTTP services and include p95/p99 thresholds.
7. Compare old vs new results with `harness/compare_results.py`.

Do not report a benchmark without the command needed to reproduce it.
