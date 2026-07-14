---
name: api-style-decision
description: Decide whether a project should expose REST/HTTP, GraphQL, gRPC, WebSocket, SSE, or CLI based on client access patterns, schema needs, real-time behavior, protocol benchmarks, and operational complexity.
---

# API Style Decision

Use this skill whenever a repository exposes an API or command interface.

1. Read `decision-brain/api-style-matrix.yaml` or `.portfolio/decision-brain/api-style-matrix.yaml`.
2. Identify the main client: public user, frontend dashboard, backend service, benchmark tool, worker, or CLI user.
3. Choose one primary API style: REST/HTTP, GraphQL, gRPC, WebSocket, SSE, or CLI.
4. Reject at least one plausible alternative.
5. Update `project.yaml` with `decision_brain.api_style`.
6. Record the decision in `sdd/technical-decision.md` when it changes contracts, benchmark, or architecture.
7. Add the correct contract artifact: OpenAPI, GraphQL schema, protobuf, event contract, or CLI output schema.
8. Add a benchmark that matches the style.

Defaults:

- Use REST/HTTP for simple resource/command APIs.
- Use GraphQL when clients need flexible nested reads or a BFF graph.
- Use gRPC for service-to-service typed protocol proof.
- Use WebSocket for bidirectional real-time behavior.
- Use SSE for server-to-client updates.
- Use CLI for harnesses, generators, and local tools.

Rule: GraphQL requires complexity/depth control and N+1 prevention.
