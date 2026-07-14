---
name: cloud-local-first
description: Apply the portfolio cloud local-first rule: use sivchari/kumo as the default local AWS emulator, keep cloud providers behind ports/adapters, and make real AWS pluggable through configuration and parity tests.
---

# Cloud Local First

Use this skill whenever a repository touches cloud-like services: object storage, queues, pub/sub, event bus, secrets, serverless functions, API gateway, NoSQL, streams, monitoring, or AWS-compatible SDKs.

1. Read `decision-brain/cloud-matrix.yaml` or `.portfolio/decision-brain/cloud-matrix.yaml`.
2. Use Kumo as the default local provider for AWS-like services.
3. Add a local Docker path for Kumo on port `4566`.
4. Keep local credentials non-secret: `test/test` or equivalent.
5. Add ports/interfaces for every cloud capability used by the application.
6. Put Kumo and real cloud SDK calls in adapters only; do not call cloud SDKs from domain or use cases.
7. Make provider switching configurable, usually `CLOUD_PROVIDER=kumo|aws` plus endpoint override.
8. Add parity tests that run against Kumo for the default path.
9. Document unsupported Kumo behaviors when real cloud behavior cannot be fully emulated.
10. Update `project.yaml` with `decision_brain.cloud`.

Default command:

```powershell
docker run -p 4566:4566 ghcr.io/sivchari/kumo:latest
```

Persistent local path:

```powershell
docker run -p 4566:4566 -e KUMO_DATA_DIR=/data -v kumo-data:/data ghcr.io/sivchari/kumo:latest
```

Rule: cloud real is never the default demo path. The portfolio proof must run locally first.
