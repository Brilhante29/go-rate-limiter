# #12 go-rate-limiter

**Status:** scaffold

**Proves:** rate limiter distribuido.

**Benchmark target:** accepted_rps, rejected_rps, p95_latency_ms.

**Stack:** go, redis, chi, k6, docker.

## Next milestone

Implement the smallest Docker-runnable version and produce the first JSON benchmark under enchmarks/results/.

## Run

`ash
docker build -t go-rate-limiter .
docker run --rm go-rate-limiter
`

## Benchmark

`ash
docker run --rm go-rate-limiter benchmark
`

| Metric | Value | Unit |
|---|---:|---|
| accepted_rps, rejected_rps, p95_latency_ms | pending | pending |

## Architecture

Defined in sdd/spec.md before implementation.

## References

See REFERENCES.md.