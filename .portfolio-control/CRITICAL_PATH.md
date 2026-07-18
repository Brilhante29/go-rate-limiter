# Critical Path: #12 go-rate-limiter

```text
intent
  -> component pack
  -> agent graph
  -> SDD/OpenSpec artifacts
  -> architecture decision
  -> stack/API/messaging/cloud decisions
  -> implementation
  -> tests and Docker
  -> benchmark JSON
  -> independent review
  -> README/post evidence
```

## Agent Handoff Contract

1. `reuse-architect` inspects the kit and proposes reusable inputs.
2. `repository-auditor` checks the repository against the gates.
3. `sdd-spec-agent` freezes the problem, constraints, architecture, and benchmark.
4. `implementation-agent` owns domain and application behavior.
5. `docker-ci-agent` owns reproducible runtime and CI.
6. `benchmark-agent` owns the measured evidence.
7. `documentation-agent` aligns README, SDD, and benchmark numbers.
8. `independent-reviewer` can reject completion when evidence is missing.

No agent silently changes another agent's responsibility. Each handoff records inputs, outputs, risks, and the next command.
