---
name: language-standards
description: Select and apply repository conventions for the chosen stack profile or implementation framework, including Spring Kotlin, FastAPI, Go backend, Node TypeScript backend, Angular, Next.js, Python, Java, TypeScript, Terraform, package layout, testing, linting, Docker, benchmark rules, and README expectations.
---

# Language Standards

Use this skill after the project program, proficiency signal, architecture, and decision brain stack profile are known.

1. Identify `decision_brain.stack_profile` from `project.yaml`.
2. Read `decision-brain/stack-matrix.yaml` or `.portfolio/decision-brain/stack-matrix.yaml` to confirm the stack was selected by problem forces.
3. Read the matching file in `language-profiles/` or `.portfolio/language-profiles/`.
4. Apply language/framework conventions without overriding the architecture, API style, cloud, and messaging decisions.
5. Ensure tests, linting, Docker, benchmark commands, README scripts, and observability match the selected profile.
6. If a project uses multiple languages, choose one primary profile and list secondary profiles in `project.yaml`.

Available concrete profiles:

- `node-typescript-backend.yaml` for NestJS/Rocketseat-style modular backends, REST/GraphQL decisions, Prisma, typed contracts, and API benchmarks.
- `spring-kotlin.yaml` for Kotlin + Spring backend, MVC/WebFlux decision, persistence, Testcontainers, observability, and k6.
- `fastapi-backend.yaml` for FastAPI, ASGI, Pydantic v2, SQLAlchemy/Alembic, async rules, and API benchmarks.
- `go-backend.yaml` for net/http/chi default, Gin/Fiber decision, pgx/sqlc, pprof, and low-overhead service benchmarks.
- `angular.yaml` for enterprise frontend, admin tools, typed forms, guards, and modular UI.
- `nextjs.yaml` for React product frontend, public demos, dashboards, and SSR/ISR tradeoffs.

Available generic profiles:

- `python.yaml` for AI, ML, data, RAG, FastAPI-adjacent tooling, MLflow, Airflow, and benchmarks.
- `java.yaml` for Spring/Java backend architecture, reliability, Kafka, and event sourcing.
- `go.yaml` for Go services where the concrete backend profile is not needed.
- `typescript.yaml` for Node, NestJS, tooling, and k6 JavaScript.
- `terraform.yaml` for infrastructure modules and AWS baselines.

Rule: language standards shape implementation hygiene; the decision brain chooses the stack.
