# Decision Brain

This directory is the decision brain for every portfolio repository.

Agents must consult these files before choosing orchestration flow, architecture, stack, API style, broker, cloud provider, database, framework, library, benchmark style, or public proof angle. The goal is to make decisions explicit, reusable, and reviewable instead of hidden inside prompts.

## Decision Order

1. `decision-brain/agentic-spec-governance.yaml`: which artifact graph, planning store, and component-pack gates govern the repo.
2. `component-packs/manifest.yaml`: which skills, templates, benchmark assets, and gates apply to the portfolio program.
3. `decision-brain/agent-graph.yaml`: how the principal agent, specialist subagents, local-first runtime, validation, and publication gates coordinate.
4. `catalog/programs.yaml`: where the repo fits in the portfolio.
5. `catalog/proficiency.yaml`: what public proficiency it should prove.
6. `architecture/decision-matrix.yaml`: which architecture fits the problem forces.
7. `decision-brain/engineering-principles.yaml`: how coupling, SOLID, LSP, KISS, YAGNI, DRY, and testability are enforced.
8. `decision-brain/stack-matrix.yaml`: which concrete stack profile fits the problem.
9. `decision-brain/api-style-matrix.yaml`: whether to expose REST/HTTP, GraphQL, gRPC, WebSocket, SSE, or CLI.
10. `decision-brain/messaging-matrix.yaml`: whether to use no broker, outbox, RabbitMQ, Kafka, Redis Streams, or NATS.
11. `decision-brain/cloud-matrix.yaml`: whether to use Kumo local-first, no cloud, an adapter fake, or real cloud.
12. `decision-brain/library-selection.yaml`: which libraries are preferred and when they are rejected.
13. `language-profiles/*.yaml`: how the chosen stack is organized, tested, benchmarked, and documented.

## Rule

Every non-trivial repo must record its decision in `project.yaml` and, when the choice affects orchestration, architecture, coupling, SOLID boundaries, cloud, API contract, or benchmark, in `sdd/technical-decision.md` and `sdd/agent-handoff.md`.

When a change affects multiple repositories or improves the reusable layer, create or update the matching OpenSpec-style artifacts under `openspec/` or `sdd/` before implementation.
