---
name: portfolio-rochedo
description: Build or review one of the 30 portfolio repositories with the required reusable standard: Docker, README opening with the project number and benchmark result, SDD, REFERENCES.md, and reproducible benchmark. Use when the user asks to create, improve, scaffold, validate, or prepare a portfolio project/repo/rochedo.
---

# Portfolio Rochedo

Apply the portfolio standard before coding:

1. Identify the project number, name, claim, stack, and benchmark metric from `catalog/projects.yaml` when available.
2. Create or update SDD first: `sdd/spec.md` and `sdd/benchmark-plan.md`.
3. Keep the default path runnable without paid credentials.
4. Ensure README starts with `# #<id> <project-name>`, then "Proves" and "Benchmark".
5. Include `REFERENCES.md` and state what was reused: dependency, architecture, or benchmark idea.
6. Add Docker support. Prefer one documented `docker build` plus one `docker run`.
7. Add a benchmark command that writes JSON under `benchmarks/results/`.
8. Validate with tests or syntax checks before finalizing.

If implementation details are open, choose the stack already defined in the catalog. Do not create empty placeholder repos.
