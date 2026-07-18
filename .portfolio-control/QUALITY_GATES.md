# Quality Gates: #12 go-rate-limiter

Completion requires evidence, not intent.

- [ ] README opens with `#12 <name>` and reports the current benchmark number.
- [ ] `project.yaml` names the problem, architecture, stack, primary metric, and result path.
- [ ] SDD and OpenSpec artifacts agree with the implementation.
- [ ] Domain logic is isolated from transport, persistence, broker, provider, and vendor details.
- [ ] SOLID, DRY, KISS, YAGNI, and Law of Demeter review has no unexplained exception.
- [ ] Tests cover the contract and the failure paths that affect the claim.
- [ ] Docker runs the documented default path from a clean checkout.
- [ ] CI runs the same meaningful checks without mutable dependencies or secrets.
- [ ] Benchmark writes valid JSON under `benchmarks/results/` and can be repeated.
- [ ] README, benchmark JSON, and `project.yaml` report the same primary metric.
- [ ] Reuse review records every kit improvement, backlog item, or rejected duplication.
- [ ] Independent review found no blocker and publication has not happened before this gate.
