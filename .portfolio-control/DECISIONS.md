# Decision Register: #12 go-rate-limiter

Decide from the problem, constraints, and proof target. The language is a consequence of the decision, not the decision itself.

| Decision | Selected option | Evidence or reason | Revisit trigger |
|---|---|---|---|
| Architecture | _pending_ | problem shape and change boundaries | coupling or test cost rises |
| API style | _pending_ | REST/GraphQL/gRPC/event decision matrix | client or latency requirement changes |
| Messaging | none until justified | Kafka/RabbitMQ adds operational cost | async delivery or fan-out is proven necessary |
| Storage | _pending_ | access pattern and consistency | query or scale assumptions change |
| Local-first/cloud | Kumo adapter where applicable | local reproducibility with cloud port | parity gap is measured |
| Libraries | _pending_ | maintenance, fit, benchmark, license | dependency becomes a bottleneck |

## Design Principles

- **SRP:** one reason to change per module and adapter.
- **OCP:** extend through ports, policies, and composition instead of editing stable domain rules.
- **LSP:** every implementation of a port preserves its contract and failure semantics.
- **ISP:** small interfaces expose only the use case that needs them.
- **DIP:** application policy owns abstractions; infrastructure depends on them.
- **DRY:** remove duplicated knowledge, not merely repeated syntax.
- **KISS/YAGNI:** choose the smallest design that proves the claim; do not prebuild distributed machinery.
- **Law of Demeter:** collaborators talk to their direct boundary, not an object graph.
