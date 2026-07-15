---
name: agentic-spec-governance
description: Apply portfolio component packs and OpenSpec-style artifact governance for reusable project planning. Use when choosing agent/skill packs, creating or updating SDD/OpenSpec artifacts, aligning a project with the portfolio planning store, deciding whether a reusable component should be added to portfolio-reuse-kit, or validating that a repository is not just isolated code.
---

# Agentic Spec Governance

Use the local `.portfolio/` snapshot first. If it is missing, use the upstream `portfolio-reuse-kit`.

## Workflow

1. Read `.portfolio/decision-brain/agentic-spec-governance.yaml` or `decision-brain/agentic-spec-governance.yaml`.
2. Select one component pack from `.portfolio/component-packs/manifest.yaml` or `component-packs/manifest.yaml`.
3. Prioritize user-owned skills and the local kit. Use external repositories only as reference material for organization, workflow, contracts, tests, benchmarks, or proven patterns.
4. Use the selected pack to decide which skills, decision sources, SDD files, benchmark assets, and validation gates apply.
5. If `openspec/config.yaml` exists, treat it as the project planning config. If OpenSpec is not installed, still follow the same artifact graph manually in `sdd/`.
6. Keep planning artifacts iterative: update intent, portfolio impact, architecture, component pack, reuse delta, benchmark proof, tasks, and verification as implementation learns new facts.
7. Never install external components, MCPs, hooks, or settings automatically. Record the recommendation first and require explicit user approval for anything that changes machine-level tooling.
8. Every reusable discovery must end in one of three states: patch the kit now, record backlog, or reject with reason.

## Required Evidence

Before implementation:

- Project has one measurable claim.
- Program and component pack are selected.
- Architecture and stack decisions name rejected alternatives.
- SOLID, LSP, KISS, dependency inversion, and testability constraints are explicit.
- Benchmark metric, command, fixture, and JSON output are defined.

Before publication:

- Reuse delta is complete.
- Benchmark proof exists.
- Validation passes.
