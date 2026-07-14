---
name: go-backend
description: Apply the Go backend stack profile, including net/http/chi default, Gin/Fiber decision rules, pgx/sqlc persistence, context cancellation, pprof, benchmarks, and low-overhead service organization.
---

# Go Backend

Use this skill when `decision_brain.stack_profile` is `go-backend`.

1. Read `language-profiles/go-backend.yaml` or `.portfolio/language-profiles/go-backend.yaml`.
2. Start with `net/http` plus `chi` for idiomatic, composable, standard-library-compatible APIs.
3. Use Gin only when ergonomics/ecosystem matter and the dependency is justified.
4. Use Fiber only when fasthttp compatibility is acceptable and raw throughput is the explicit benchmark claim.
5. Use `pgx` for PostgreSQL; use `sqlc` when explicit SQL contracts strengthen the repo.
6. Respect `context.Context` cancellation in HTTP, database, and worker paths.
7. Include pprof or allocation benchmarks when performance is part of the claim.
8. Use k6 for HTTP and ghz for gRPC/protocol benchmarks.

Rule: Go repos must prove low-overhead engineering with measurements, not just minimal code.
