---
name: engineering-principles
description: Enforce decoupling, SOLID, LSP, KISS, YAGNI, DRY, dependency inversion, ports/adapters, and testability evidence across portfolio repositories before implementation and review.
---

# Engineering Principles

Use this skill before implementing or reviewing architecture-sensitive code.

1. Read `decision-brain/engineering-principles.yaml` or `.portfolio/decision-brain/engineering-principles.yaml`.
2. Treat `LISP` in architecture discussion as `LSP`: Liskov Substitution Principle.
3. Verify dependency direction: domain and use cases must not depend on HTTP, GraphQL, ORM, database driver, cloud SDK, broker SDK, framework runtime, env vars, or UI rendering.
4. Ensure controllers, resolvers, CLI handlers, consumers, and HTTP handlers are thin and delegate to use cases/application services.
5. Put external systems behind small capability-specific ports.
6. Use adapters for database, cloud, broker, HTTP clients, files, and local fakes/emulators.
7. Apply SOLID pragmatically: SRP, OCP, LSP, ISP, DIP.
8. Apply KISS and YAGNI: reject infrastructure or abstractions that do not improve the claim, benchmark, testability, or replaceability.
9. Apply DRY only to real duplicated knowledge; avoid premature generic abstractions.
10. Require tests that prove use cases run without transport and infrastructure.
11. Update `sdd/technical-decision.md` with coupling boundaries and rejected overengineering when architecture is non-trivial.

Review checklist:

- No framework imports in domain.
- No cloud/broker/database SDK calls in use cases.
- Ports are small and named by business capability.
- Multiple adapters obey the same contract and error semantics.
- Kumo/local adapter and real cloud adapter are substitutable when cloud is used.
- Default Docker path stays simple.
