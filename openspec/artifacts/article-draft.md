# #12 go-rate-limiter: total_rps = 11511.03 requests_per_second

Two HTTP nodes enforce one global token-bucket quota through atomic Redis state.

This repository belongs to the Backend Reliability and Architecture Platform program. Its job is narrow: prove the measurable claim through the selected component pack before adding unrelated infrastructure or features.

The benchmark is the proof. total_rps = 11511.03 requests_per_second. P95 latency: 11.98 ms. The result is stored in `benchmarks/results/rate-limiter-baseline.json` and can be reproduced from the Docker/local path.

The important architecture decision is hexagonal. The hard boundary is substituting local memory and remote atomic Redis behavior without coupling policy or HTTP to either adapter.

The default path stays local-first. The project uses go-backend, exposes rest-http, uses messaging mode `none`, and stores data with `redis`. The dependency rule is explicit: HTTP and store adapters depend inward on limiter ports and types; limiter code imports no framework or Redis package.

The rejected work matters as much as the implemented work. Anything that does not improve the benchmark stays out of the first version.

Post angle: start with the number, show the architecture boundary, then explain which future adapter can be added without changing the core use cases.
