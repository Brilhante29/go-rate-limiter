---
name: spec-driven-project
description: Write or update Specification Driven Development artifacts for portfolio repositories, including program fit, proficiency signal, decision_brain, engineering principles, project.yaml, architecture decision, stack decision, API style, cloud local-first, messaging decision, language/framework profile, design-system expectations, spec.md, benchmark-plan.md, ADRs, release checklist, acceptance criteria, and publication readiness.
---

# Spec Driven Project

Drive implementation from explicit project evidence. In generated projects, prefer `.portfolio/` standards when present.

1. Convert the project idea into one measurable claim.
2. Place the repo into a program from `catalog/programs.yaml` or `.portfolio/catalog/programs.yaml` and explain how it strengthens that program narrative.
3. Define the public proficiency signal from `catalog/proficiency.yaml` or `.portfolio/catalog/proficiency.yaml`.
4. Define in-scope behavior, out-of-scope behavior, and default Docker path.
5. Choose the software architecture from `architecture/decision-matrix.yaml` or `.portfolio/architecture/decision-matrix.yaml` using problem forces.
6. Define coupling, SOLID/LSP, KISS/YAGNI/DRY, and testability rules from `decision-brain/engineering-principles.yaml` or `.portfolio/decision-brain/engineering-principles.yaml`.
7. Choose concrete stack profile from `decision-brain/stack-matrix.yaml` or `.portfolio/decision-brain/stack-matrix.yaml`.
8. Choose API style from `decision-brain/api-style-matrix.yaml` or `.portfolio/decision-brain/api-style-matrix.yaml`.
9. Choose messaging mode from `decision-brain/messaging-matrix.yaml` or `.portfolio/decision-brain/messaging-matrix.yaml`.
10. Choose cloud mode from `decision-brain/cloud-matrix.yaml` or `.portfolio/decision-brain/cloud-matrix.yaml`; use Kumo for AWS-like local-first cloud.
11. Choose core libraries from `decision-brain/library-selection.yaml` or `.portfolio/decision-brain/library-selection.yaml`.
12. Choose the primary language/framework profile from `language-profiles/` or `.portfolio/language-profiles/` after the stack decision.
13. Apply `design-system/tokens.yaml` or `.portfolio/design-system/tokens.yaml` to README structure, architecture summary, diagrams, and benchmark tables.
14. Record program, decision_brain, architecture, language profile, design system, and benchmark in `project.yaml`.
15. Create `sdd/architecture-decision.md` and `sdd/technical-decision.md`.
16. Specify metric name, unit, target, benchmark command, fixture, seed, and JSON result path.
17. Identify reusable assets from the kit: programs, proficiency map, decision brain, architecture catalog, language profiles, design system, templates, contracts, harness, CI, skills, and references.
18. Add release criteria that include Docker, tests, benchmark JSON, references, README result, and no-secret default demo.

Use `catalog/`, `decision-brain/`, `architecture/`, `language-profiles/`, `design-system/`, `sdd/templates/`, and `contracts/` when present.
