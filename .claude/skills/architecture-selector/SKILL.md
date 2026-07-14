---
name: architecture-selector
description: Choose and document the best software architecture for a portfolio project from the problem forces, including MVC, layered, modular monolith, Clean Architecture, Hexagonal, MVVM, pipeline, event-driven, CQRS/event sourcing, serverless, or microservices.
---

# Architecture Selector

Choose architecture by problem shape, not by stack.

1. Read `architecture/decision-matrix.yaml` or `.portfolio/architecture/decision-matrix.yaml`, plus `docs/architecture-decision-guide.md` when available.
2. Identify the project forces: domain complexity, integration pressure, UI state complexity, data/ML reproducibility, auditability, async/throughput pressure, and independent deployability need.
3. Pick the simplest architecture that still proves the project claim and benchmark.
4. Update `project.yaml` with `architecture.style`, `architecture.reason`, `architecture.dependency_rule`, `architecture.boundaries`, and `architecture.rejected_alternatives`.
5. Create an ADR from `sdd/templates/architecture-decision.md` when the choice affects folder layout, tests, dependencies, or benchmark methodology.
6. Organize folders so the dependency rule is visible from the repository tree.
7. Do not copy another repository's implementation. Reuse organization ideas, documentation patterns, and benchmark approach with attribution in `REFERENCES.md`.

Default decisions:

- Use MVC or layered for simple CRUD and admin APIs.
- Use modular monolith when multiple business contexts belong in one deployable app.
- Use Clean Architecture when use cases and domain rules are the main evidence.
- Use Hexagonal when external adapters, providers, queues, or databases are central.
- Use MVVM for rich UI state.
- Use pipeline for ML, CV, data, ETL, or reproducible batch workflows.
- Use event-driven for queues, streams, lag, and async consistency.
- Use CQRS/event sourcing when audit history or replay is part of the claim.
- Use microservices only when independent deployability is part of the evidence.

A valid architecture decision must explain rejected alternatives.


