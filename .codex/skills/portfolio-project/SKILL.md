---
name: portfolio-project
description: Build, review, or upgrade one repository in the 30-project portfolio using the shared reuse layer: agent graph, reuse-improvement loop, program catalog, proficiency map, decision brain, engineering principles, architecture decision, stack profile, API style, cloud local-first, messaging, language/framework profile, design system, SDD, Docker path, benchmark contract, references, validation, and publication workflow.
---

# Portfolio Project

Use the portfolio reuse layer as the source of truth. In this kit, read files from the repository root. In generated projects, read the same standards from `.portfolio/` when present.

Decision order:

1. Read `decision-brain/agent-graph.yaml` or `.portfolio/decision-brain/agent-graph.yaml` and follow the principal-agent/subagent handoff flow.
2. Read `decision-brain/reuse-improvement-loop.yaml` or `.portfolio/decision-brain/reuse-improvement-loop.yaml` and keep a running review of kit improvements.
3. Read `catalog/programs.yaml` or `.portfolio/catalog/programs.yaml` and place the repo inside one portfolio program. Do not create isolated projects.
4. Read `catalog/projects.yaml` or `.portfolio/catalog/projects.yaml` to identify project id, claim, benchmark target, and reference ideas.
5. Read `catalog/proficiency.yaml` or `.portfolio/catalog/proficiency.yaml` to choose what the repo should prove publicly on GitHub/LinkedIn.
6. Read `architecture/decision-matrix.yaml` or `.portfolio/architecture/decision-matrix.yaml` and choose the architecture that fits the problem forces, not the stack.
7. Read `decision-brain/engineering-principles.yaml` or `.portfolio/decision-brain/engineering-principles.yaml` and define coupling, SOLID/LSP, KISS/YAGNI/DRY, and testability rules.
8. Read `decision-brain/stack-matrix.yaml`, `api-style-matrix.yaml`, `messaging-matrix.yaml`, `cloud-matrix.yaml`, and `library-selection.yaml` or their `.portfolio/decision-brain/` copies.
9. Choose concrete stack profile, API style, messaging mode, cloud mode, database/runtime, and library policy. Reject plausible alternatives.
10. Use Kumo as the default local-first provider for AWS-like cloud capabilities; real cloud must be behind ports/adapters.
11. Read the matching file in `language-profiles/` or `.portfolio/language-profiles/` and apply stack-specific repository standards.
12. Read `design-system/tokens.yaml` or `.portfolio/design-system/tokens.yaml` and apply shared README, diagram, badge, benchmark, and dashboard conventions.
13. Create or update `project.yaml` from `templates/project.yaml`, including program, decision_brain, language profile, architecture, design system, benchmark, release, and reuse sections.
14. Create or update `sdd/spec.md`, `sdd/benchmark-plan.md`, `sdd/architecture-decision.md`, and `sdd/technical-decision.md` before implementation.
15. Create or update `sdd/agent-handoff.md` with subagent decisions, local-first runtime, benchmark handoff, risks, and release gates.
16. Create or update `sdd/reuse-improvement-review.md`, remove template placeholders, classify findings as patch_now/backlog/reject, and mark the final gate complete before publication.
17. Keep the default path runnable without paid credentials.
18. Ensure `README.md` starts with the project number, claim, and benchmark status/result.
19. Include `REFERENCES.md` and state whether reuse is dependency, architecture, organization pattern, benchmark idea, or documentation pattern.
20. Add Docker support and one documented command path.
21. Add a benchmark command that writes JSON compatible with `contracts/benchmark-result.schema.json` or `.portfolio/contracts/benchmark-result.schema.json` under `benchmarks/results/`.
22. Run local validation before commit or publication.

Do not present scaffold-only repositories as finished portfolio evidence.
