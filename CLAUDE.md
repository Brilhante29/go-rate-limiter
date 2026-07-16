# Claude Code Project Context

This repository is governed by `portfolio-reuse-kit`.

Read `AGENTS.md` first. Then read `.aitmpl/context-card.md` when it exists,
`project.yaml`, the selected files under `.portfolio/`, and the relevant SDD
and OpenSpec artifacts before changing code.

Use the same execution graph as Codex:

1. Confirm the portfolio program and measurable claim.
2. Confirm architecture and dependency direction.
3. Confirm the stack, API style, messaging, cloud, database, and library decisions.
4. Keep the default path local-first and free of paid secrets.
5. Run the benchmark and project validation.
6. Review whether the reuse kit should improve because of the work.

Do not install external AITmpl components or other machine-level tooling unless
the user explicitly asks for it. Local skills and `.portfolio/` remain the
primary source of truth.
