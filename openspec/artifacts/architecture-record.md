# Architecture Record: go-rate-limiter

## Decision

- Architecture: `hexagonal`
- Stack profile: `go-backend`
- API style: `rest-http`
- Messaging: `none`
- Database/runtime: `redis` / `Go 1.26 static binary plus Redis 8.8`

## Reason

The hard boundary is substituting local memory and remote atomic Redis behavior without coupling policy or HTTP to either adapter.

## Dependency Direction

HTTP and store adapters depend inward on limiter ports and types; limiter code imports no framework or Redis package.

## Boundaries

- limiter policy and application service
- bucket-store output port
- Redis and memory output adapters
- chi HTTP input adapter
- load benchmark and demo orchestration

## Library Policy

chi preserves net/http compatibility; go-redis is the official Redis client; the atomic Lua script is project-owned.

## Principle Check

- SRP: keep benchmark, API, use cases, and adapters separate.
- OCP: new providers must be adapters, not domain rewrites.
- LSP: replacement providers must preserve observable behavior.
- ISP: ports stay narrow.
- DIP: application depends on behavior, not infrastructure.
- KISS/YAGNI: leave out anything that does not improve the benchmark.
