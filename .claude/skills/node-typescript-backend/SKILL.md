---
name: node-typescript-backend
description: Apply the Node.js/TypeScript backend profile using NestJS/Rocketseat-style modular Clean Architecture, Prisma, REST vs GraphQL decisions, Express vs Fastify tradeoffs, validation, tests, and benchmarks.
---

# Node TypeScript Backend

Use this skill when `decision_brain.stack_profile` is `node-typescript-backend`.

1. Read `language-profiles/node-typescript-backend.yaml` or `.portfolio/language-profiles/node-typescript-backend.yaml`.
2. Prefer NestJS for structured backend projects that should show modular architecture, dependency injection, controllers/resolvers, providers, guards, pipes, and testability.
3. Use Rocketseat-style organization as a reference for clean boundaries: `core`, `domain`, `application/use-cases`, `infra`, and `test`.
4. Choose REST/HTTP by default for CRUD and command APIs.
5. Choose GraphQL only when flexible client graph reads, BFF aggregation, or schema-driven client workflow is part of the proof.
6. Use Prisma by default for type-safe persistence and migrations; choose Drizzle when SQL visibility is more important.
7. Use Express by default; use Fastify when p95/p99 performance is part of the benchmark and middleware compatibility is checked.
8. Keep controllers/resolvers thin and route everything through use cases.
9. Add Vitest/Jest and Supertest for HTTP tests; add GraphQL operation tests when GraphQL is selected.
10. Export OpenAPI or GraphQL schema when relevant.

Rule: TypeScript backend repos must prove structure and typed contracts, not just Node endpoints.
