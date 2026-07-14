---
name: agent-orchestration
description: Use the portfolio agent graph to coordinate principal-agent work, specialist subagent responsibilities, local-first execution, benchmark evidence, validation, reuse improvement, and publication gates for generated portfolio repositories.
---

# Agent Orchestration

Use this skill when creating, upgrading, reviewing, or publishing a portfolio repository.

Read `decision-brain/agent-graph.yaml` in this kit, or `.portfolio/decision-brain/agent-graph.yaml` inside a generated project, before implementation.

## Operating Rule

The principal agent owns the user objective, sequencing, conflict resolution, validation, reuse improvement, and publication decision.

Specialist subagents own narrow decisions:

1. `program-planner`
2. `architecture-selector`
3. `engineering-principles-reviewer`
4. `stack-decision-agent`
5. `api-style-agent`
6. `cloud-local-first-agent`
7. `messaging-agent`
8. `language-profile-agent`
9. `benchmark-harness-agent`
10. `design-system-agent`
11. `security-reuse-reviewer`
12. `reuse-improvement-reviewer`
13. `release-ci-publisher`

If a runtime cannot spawn subagents, execute the same roles sequentially and record the outputs in `sdd/agent-handoff.md` and `sdd/reuse-improvement-review.md`.

## Required Decision Order

1. Program and portfolio fit.
2. Architecture by problem forces.
3. Decoupling, SOLID, LSP, KISS, YAGNI, DRY, and testability rules.
4. Concrete stack profile.
5. API style.
6. Cloud local-first mode.
7. Messaging mode.
8. Language/framework profile.
9. Benchmark plan.
10. Design-system presentation.
11. Security and reuse review.
12. Reuse improvement review.
13. Release validation.

## Local-First Rule

Default demos must run without paid credentials. Use Docker for execution and Kumo for AWS-like local cloud behavior. Real cloud providers are adapters behind ports and must not be imported by domain or use-case code.

## Publication Gate

Do not publish or call a project finished until these exist:

- `project.yaml`
- `sdd/spec.md`
- `sdd/architecture-decision.md`
- `sdd/technical-decision.md`
- `sdd/benchmark-plan.md`
- `sdd/agent-handoff.md`
- `sdd/reuse-improvement-review.md`
- Docker run path
- benchmark JSON
- README number, claim, and result
- `REFERENCES.md`
- passing validation
