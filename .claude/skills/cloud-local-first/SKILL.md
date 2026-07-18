---
name: cloud-local-first
description: Apply the portfolio cloud local-first rule with Kumo pinned by digest, narrow ports, guarded real AWS switching, scoped parity tests, and numeric compatibility diagnostics.
---

# Cloud Local First

Use this skill whenever a repository touches AWS-compatible storage, queues, pub/sub, events, secrets, functions, API gateway, NoSQL, streams, monitoring, or SDKs.

1. Read `decision-brain/cloud-matrix.yaml` or `.portfolio/decision-brain/cloud-matrix.yaml`.
2. Use the reviewed Kumo release as the default local provider.
3. Pin the Kumo image by readable version tag and immutable digest; reject committed mutable references.
4. Keep local credentials non-secret: `test/test` or equivalent.
5. Add narrow ports/interfaces for every claimed cloud capability.
6. Keep AWS SDK calls in adapters; domain and use cases depend only on ports.
7. Use one adapter configuration path: Kumo adds endpoint override, real AWS uses default endpoints.
8. Require an explicit destructive-action guard and unique run ID before real AWS execution.
9. Add scoped parity/conformance checks and list every unsupported behavior.
10. Include provider version, digest, SDK version, region, operation counts, and compatibility diagnostics in benchmark JSON.
11. Update `project.yaml` and `sdd/reuse-improvement-review.md` after the benchmark.

Reviewed command:

```powershell
docker run -p 4566:4566 ghcr.io/sivchari/kumo:0.25.3@sha256:7ea090ae0b6d1d34615e8b7bd04a2f1cd864ec640a6826a91e90f40e975e196b
```

Persistent local path:

```powershell
docker run -p 4566:4566 -e KUMO_DATA_DIR=/data -v kumo-data:/data ghcr.io/sivchari/kumo:0.25.3@sha256:7ea090ae0b6d1d34615e8b7bd04a2f1cd864ec640a6826a91e90f40e975e196b
```

Do not claim full AWS conformance from a scoped emulator suite. Real cloud is never the default demo path.
