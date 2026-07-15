# Technical Decision: Go, REST, Redis, Lua, And k6

## Status

Accepted

## Selected Stack

- Go 1.26.5
- chi 5.2.5 over `net/http`
- go-redis 9.20.0
- Redis 8.8.0
- Redis Lua script for atomic token-bucket transition
- built-in Go load generator for committed JSON
- k6 2.1.0 for an independent Compose load profile
- Docker multi-stage build and GitHub Actions

## Decision Reasons

- Go exposes concurrency, allocations, cancellation, connection pooling, and low HTTP overhead directly.
- chi adds routing and recovery middleware without replacing `net/http` contracts.
- go-redis is the official Go client and supports script execution and context cancellation.
- Redis provides a shared atomic execution point. Its `TIME` command removes application-node clock skew from the token calculation.
- The Go benchmark owns the portfolio result schema and correctness invariant; k6 independently exercises thresholds and expected 429 responses.

## API Style

Selected: REST/HTTP.

The operation is a synchronous command with one key and one cost. GraphQL selection, gRPC schema/tooling, WebSocket state, and asynchronous events do not improve this contract. OpenAPI lives at `api/openapi.yaml`.

## Messaging

Selected: none.

RabbitMQ, Kafka, Redis Streams, and NATS are rejected because a caller needs an immediate atomic decision. A queue would increase latency and create a second consistency problem.

## Cloud And Local-First

Selected cloud mode: none.

The complete proof runs in Docker with local Redis. Kumo is reserved for repositories that emulate AWS APIs; this project has no AWS behavior to emulate. A real managed Redis-compatible endpoint changes `REDIS_ADDR` and credentials at the adapter boundary, but is intentionally not part of the default path.

## Runtime And Security

- The service image is a static non-root binary.
- The self-contained demo image adds Redis only to satisfy the one-command portfolio proof; Compose keeps app and store in separate containers.
- Keys are restricted to 128 safe characters and prefixed before reaching Redis.
- Request bodies are limited to 1 KiB and reject unknown fields.
- HTTP timeouts, Redis timeouts, context cancellation, and graceful shutdown are explicit.
- Store errors fail closed; logs do not include keys or request bodies.

## Benchmark Impact

The primary performance number is `total_rps`. It is interpreted with `accepted_rps`, `rejected_rps`, p95, node count, and `global_limit_preserved`. Throughput alone would hide whether overload was rejected atomically or whether the benchmark reached both nodes.

## Revisit When

- Redis becomes a measured bottleneck;
- policies must vary dynamically by tenant;
- cluster/failover or multi-region semantics become part of the claim;
- a gateway integration needs a standard external rate-limit protocol.
