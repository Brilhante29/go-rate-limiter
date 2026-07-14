---
name: spring-kotlin-backend
description: Apply the Spring Kotlin backend stack profile, including MVC vs WebFlux, persistence choice, hexagonal/clean/modular organization, Testcontainers, observability, and k6 benchmark rules.
---

# Spring Kotlin Backend

Use this skill when `decision_brain.stack_profile` is `spring-kotlin-backend`.

1. Read `language-profiles/spring-kotlin.yaml` or `.portfolio/language-profiles/spring-kotlin.yaml`.
2. Choose Spring MVC by default for transactional REST APIs with blocking persistence.
3. Choose WebFlux only when the app is non-blocking end to end and the benchmark needs that behavior.
4. Choose Spring Data JDBC by default for explicit persistence; use JPA only when ORM mapping is part of the proof; use R2DBC only with WebFlux end to end.
5. Organize packages so domain/application do not depend on adapters or framework code.
6. Add Testcontainers for PostgreSQL, Redis, RabbitMQ, Kafka/Redpanda, or other external services used by tests.
7. Add Actuator and Micrometer when observability or production readiness is part of the claim.
8. Add k6 benchmark for HTTP p95/p99 and failure benchmark for outbox/saga/broker projects.

Rule: Spring Kotlin repos must show Kotlin fluency plus backend architecture, not just generated Spring code.
