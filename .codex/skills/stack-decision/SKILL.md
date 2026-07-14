---
name: stack-decision
description: Choose the concrete stack profile for a portfolio repository from the decision brain, including Spring Kotlin backend, FastAPI backend, Go backend, Node TypeScript backend, Angular, Next.js, Python ML, or Terraform, based on problem forces, public signal, architecture, API style, benchmark, and rejected alternatives.
---

# Stack Decision

This skill must run before implementation when a repository needs a stack choice.

1. Read `decision-brain/stack-matrix.yaml` or `.portfolio/decision-brain/stack-matrix.yaml`.
2. Read `catalog/proficiency.yaml` or `.portfolio/catalog/proficiency.yaml`.
3. Read the architecture decision from `project.yaml` or `sdd/architecture-decision.md`.
4. Select the concrete stack profile that best proves the problem and benchmark.
5. Reject at least one plausible alternative with a technical reason.
6. Update `project.yaml` with `decision_brain.stack_profile` and `language_profile.primary`.
7. Create or update `sdd/technical-decision.md` from `sdd/templates/technical-decision.md` when the stack changes libraries, runtime, Docker services, API style, or benchmark methodology.

Default choices:

- Spring Kotlin backend for transactional enterprise domain systems, payments, consistency, and Testcontainers-heavy backend proof.
- FastAPI backend for Python APIs around ML, RAG, CV, evaluation, data, and OpenAPI-first services.
- Go backend for gateways, rate limiters, emulators, protocol benchmarks, low overhead services, and networking-heavy proof.
- Node TypeScript backend for NestJS/Rocketseat-style modular architecture, Prisma, REST/GraphQL contracts, BFFs, and product APIs.
- Angular for enterprise admin/dashboard UI with typed forms and modular frontend structure.
- Next.js for public product UI, dashboard demos, React proficiency, and shareable portfolio surfaces.

Rule: never write "use Python/Java/Go/Node" as the decision. Use a concrete profile.
