# Architecture Decision: Hexagonal Rate Limiter

## Status

Accepted

## Context

The project must prove one quota across multiple HTTP nodes. The core policy is small, but integration pressure is high because correctness depends on an atomic remote state transition. The backing store must be replaceable in tests without pretending that an in-memory limiter is distributed.

## Decision

Use a narrow Hexagonal architecture:

- `limiter.Service` validates key/cost and calls a `BucketStore` output port;
- `redisstore.Store` implements the port with one atomic Lua script and Redis server time;
- `MemoryStore` implements the same port for deterministic unit tests and local single-node use;
- `httpapi` is an input adapter built on chi and `net/http`;
- benchmark and demo orchestration stay outside the policy boundary.

Dependency rule:

```text
httpapi -> limiter <- redisstore
demo/loadbench -> public behavior, never internal policy mutation
```

## SOLID And Simplicity

- SRP: policy, Redis atomicity, HTTP translation, load generation, and process orchestration have separate packages.
- OCP: another atomic store can implement `BucketStore` without changing service or HTTP behavior.
- LSP: memory and Redis adapters preserve allowed/remaining/retry semantics; only Redis is used for the distributed claim.
- ISP: the application port exposes only `Take`; health, reset, and close remain adapter lifecycle methods.
- DIP: the service depends on `BucketStore`, not go-redis.
- KISS: one use case, one endpoint, one script, no policy control plane.
- YAGNI: no broker, cluster abstraction, repository layer, event stream, or service mesh.

## Rejected Alternatives

| Alternative | Why rejected |
|---|---|
| In-memory-only limiter | Concurrent-safe inside one process but cannot enforce one quota across nodes. |
| Layered controller/service/repository | Adds familiar folders but obscures the behavioral store port and adapter substitution. |
| Microservices control plane | Dynamic policy administration is not required to prove the quota invariant. |
| Client-side Redis transaction | Read-modify-write and node clock differences make correctness harder than one server-side script. |
| Fixed-window counter | Simpler, but produces boundary bursts and does not prove refill behavior. |

## Consequences

Positive:

- The distributed proof and local unit path use the same application contract.
- Infrastructure failures fail closed and are visible as `503`.
- Redis-specific code is isolated and integration-tested by the Docker demo.

Tradeoffs:

- The Lua script is a critical adapter artifact and needs contract-level integration coverage.
- One Redis instance is a coordination dependency and throughput ceiling.
- Cluster slotting and failover require additional design before production use.

Migration path: preserve `BucketStore`, add a separately benchmarked Redis Cluster or managed Redis adapter, and record changed consistency/failure semantics before claiming equivalence.
