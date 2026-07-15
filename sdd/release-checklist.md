# Release Checklist

- [x] `go test -race ./...` passes.
- [x] `go vet ./...` passes.
- [x] Docker build includes tests and produces the static service and self-contained demo targets.
- [x] `docker run --rm go-rate-limiter` reaches two nodes and preserves the global quota.
- [x] Benchmark result is stored under `benchmarks/results/`.
- [x] README opens with project number and measured result.
- [x] OpenAPI, architecture, technical decision, and agent handoff exist.
- [x] `REFERENCES.md` records licenses and confirms no copied code.
- [x] Default path requires no API key or paid service.
- [x] Project-specific code remains outside the reuse kit.
- [x] OpenSpec artifacts are generated from the final manifest and result.
- [x] Portfolio validation passes.
- [x] GitHub Actions is green on the published commit.
