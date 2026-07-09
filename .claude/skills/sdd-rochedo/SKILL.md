---
name: sdd-rochedo
description: Write or update Specification Driven Development artifacts for portfolio projects, including spec.md, benchmark-plan.md, ADRs, release checklist, and post angle. Use before implementation or when the user asks to define scope, architecture, metrics, benchmark, or acceptance criteria.
---

# SDD Rochedo

Before implementation:

1. Write the claim in one sentence.
2. Define in-scope and out-of-scope behavior.
3. Specify the exact benchmark command, metric, unit, target, and output file.
4. Define the Docker run path.
5. Add dataset or fixture details, including license and deterministic seed.
6. Add a definition of done that includes README, Docker, benchmark JSON, references, and tests.

Use `sdd/templates/spec.md` and `sdd/templates/benchmark-plan.md` if present.

When a decision affects architecture, dependency, storage, protocol, model, or benchmark design, add an ADR using `sdd/templates/adr.md`.
