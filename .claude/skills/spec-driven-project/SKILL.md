---
name: spec-driven-project
description: Write or update Specification Driven Development artifacts for portfolio repositories, including OpenSpec-style artifact graph, component pack selection, program fit, proficiency signal, decision_brain, engineering principles, project.yaml, architecture decision, stack decision, API style, cloud local-first, messaging decision, language/framework profile, design-system expectations, spec.md, benchmark-plan.md, ADRs, release checklist, acceptance criteria, and publication readiness.
---

# Spec Driven Project

Drive implementation from explicit project evidence. In generated projects, prefer `.portfolio/` standards when present.

1. Convert the project idea into one measurable claim.
2. Select the component pack from `component-packs/manifest.yaml` or `.portfolio/component-packs/manifest.yaml`.
3. Follow `decision-brain/agentic-spec-governance.yaml` or `.portfolio/decision-brain/agentic-spec-governance.yaml` for the artifact graph.
4. If `openspec/config.yaml` exists, keep it coherent with `project.yaml` and `sdd/`; if OpenSpec is unavailable, use the same artifacts manually.
5. Place the repo into a program from `catalog/programs.yaml` or `.portfolio/catalog/programs.yaml` and explain how it strengthens that program narrative.
6. Define the public proficiency signal from `catalog/proficiency.yaml` or `.portfolio/catalog/proficiency.yaml`.
7. Define in-scope behavior, out-of-scope behavior, and default Docker path.
8. Choose the software architecture from `architecture/decision-matrix.yaml` or `.portfolio/architecture/decision-matrix.yaml` using problem forces.
9. Define coupling, SOLID/LSP, KISS/YAGNI/DRY, and testability rules from `decision-brain/engineering-principles.yaml` or `.portfolio/decision-brain/engineering-principles.yaml`.
10. Choose concrete stack profile from `decision-brain/stack-matrix.yaml` or `.portfolio/decision-brain/stack-matrix.yaml`.
11. Choose API style from `decision-brain/api-style-matrix.yaml` or `.portfolio/decision-brain/api-style-matrix.yaml`.
12. Choose messaging mode from `decision-brain/messaging-matrix.yaml` or `.portfolio/decision-brain/messaging-matrix.yaml`.
13. Choose cloud mode from `decision-brain/cloud-matrix.yaml` or `.portfolio/decision-brain/cloud-matrix.yaml`; use Kumo for AWS-like local-first cloud.
14. Choose core libraries from `decision-brain/library-selection.yaml` or `.portfolio/decision-brain/library-selection.yaml`.
15. Choose the primary language/framework profile from `language-profiles/` or `.portfolio/language-profiles/` after the stack decision.
16. Apply `design-system/tokens.yaml` or `.portfolio/design-system/tokens.yaml` to README structure, architecture summary, diagrams, and benchmark tables.
17. Record program, decision_brain, architecture, language profile, design system, and benchmark in `project.yaml`.
18. Create `sdd/architecture-decision.md` and `sdd/technical-decision.md`.
19. Specify metric name, unit, target, benchmark command, fixture, seed, and JSON result path.
20. Identify reusable assets from the kit: programs, proficiency map, decision brain, component packs, OpenSpec schema, architecture catalog, language profiles, design system, templates, contracts, harness, CI, skills, and references.
21. Add release criteria that include Docker, tests, benchmark JSON, references, README result, and no-secret default demo.

Use `catalog/`, `decision-brain/`, `architecture/`, `language-profiles/`, `design-system/`, `sdd/templates/`, and `contracts/` when present.
