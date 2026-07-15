# Spec: go-rate-limiter

## Number

#12

## Measurable Claim

Two independently addressed HTTP nodes enforce one global token-bucket quota through atomic Redis state.

## Acceptance Criteria

- A request accepted by either node decrements the same bucket.
- Concurrent decisions cannot admit more than `burst + rate * elapsed` tokens, allowing one token for timing precision.
- Excess traffic returns HTTP 429 with remaining tokens and retry delay.
- Redis failures fail closed with HTTP 503.
- The benchmark reaches both node identities, records no unexpected response, and produces portfolio-compatible JSON.
- `docker run --rm go-rate-limiter` executes the complete proof without a paid secret.

## Scope

In:

- token-bucket policy with configurable rate, burst, key, and request cost;
- Redis-backed atomic adapter and deterministic memory adapter;
- REST API and OpenAPI contract;
- two-node Docker demo, Compose topology, Go benchmark, and k6 profile;
- unit, race, vet, allocation, integration, and portfolio validation gates.

Out:

- multi-region quota semantics;
- Redis Cluster or Sentinel failover claims;
- dynamic policy administration or control plane;
- authentication, billing, or tenant provisioning;
- message brokers, GraphQL, gRPC, or cloud SDKs.

## Architecture

```text
HTTP/chi -> limiter service -> BucketStore port <- Redis Lua adapter
                                ^
                                `- memory adapter for unit/local use

Go benchmark -> node-a + node-b -> shared Redis bucket -> result JSON
```

## Failure Behavior

- Invalid key/cost/body: `400`.
- Exhausted bucket: `429`; this is an expected decision, not an infrastructure error.
- Redis unavailable or script error: `503`; no request is admitted speculatively.
- Missing node, unexpected status, benchmark error, or quota violation: benchmark exits non-zero.

## Definition Of Done

- [x] Docker command is self-contained.
- [x] Core behavior has deterministic and concurrent tests.
- [x] API contract is versioned.
- [x] Benchmark writes JSON and checks the global invariant.
- [x] Architecture and technical alternatives are recorded.
- [x] Reuse improvements are patched, backlogged, or rejected.
- [x] README and SDD contain the measured baseline.
- [x] GitHub CI is green on the published commit.
