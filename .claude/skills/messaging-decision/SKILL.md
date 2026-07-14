---
name: messaging-decision
description: Decide whether a portfolio repository should use no broker, transactional outbox, RabbitMQ, Kafka/Redpanda, Redis Streams, or NATS based on delivery semantics, ordering, replay, retry/DLQ, throughput, and benchmark evidence.
---

# Messaging Decision

Use a broker only when the problem needs broker semantics.

1. Read `decision-brain/messaging-matrix.yaml` or `.portfolio/decision-brain/messaging-matrix.yaml`.
2. Identify the actual messaging need: none, async reliability, work queue, routing, retry/DLQ, event log, replay, stream processing, fanout, or lightweight pub/sub.
3. Choose one messaging mode and reject the alternatives.
4. Update `project.yaml` with `decision_brain.messaging`.
5. If a broker is selected, document topology and failure case in `sdd/technical-decision.md`.
6. Add a benchmark that proves the broker choice: throughput, lag, retry success, lost messages under failure, or replay time.

Defaults:

- Use no broker when synchronous flow is enough.
- Use transactional outbox when reliability is needed but a broker is not central to the default proof.
- Use RabbitMQ for work queues, routing, acknowledgements, retries, DLQ, and task processing.
- Use Kafka or Redpanda for event streaming, durable logs, replay, partitions, CQRS projections, CDC, and high-throughput pipelines.
- Use Redis Streams only when Redis is already in the stack and lightweight stream semantics are enough.
- Use NATS only when lightweight pub/sub or request/reply is the core proof.

Rule: every broker choice must include rejected alternatives and a failure benchmark.
