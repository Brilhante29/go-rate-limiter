---
name: fastapi-backend
description: Apply the FastAPI backend stack profile, including ASGI runtime, Pydantic v2 contracts, async decisions, SQLAlchemy/Alembic persistence, observability, tests, Docker, and API benchmarks.
---

# FastAPI Backend

Use this skill when `decision_brain.stack_profile` is `fastapi-backend`.

1. Read `language-profiles/fastapi-backend.yaml` or `.portfolio/language-profiles/fastapi-backend.yaml`.
2. Use FastAPI for typed API contracts, OpenAPI, ML/data integration, and Python ecosystem leverage.
3. Use async only when database/client/broker libraries are async end to end.
4. Do not share one SQLAlchemy `AsyncSession` across concurrent tasks.
5. Use Pydantic v2 models for request/response contracts; keep domain models separate when domain is non-trivial.
6. Use SQLAlchemy 2 plus Alembic for non-trivial persistence; use SQLModel only for small CRUD or teaching clarity.
7. Use structlog or equivalent structured logging and expose metrics when benchmark/observability matters.
8. Add k6 or similar load benchmark that writes JSON under `benchmarks/results`.

Rule: FastAPI repos must prove backend engineering, not just a thin ML script behind an endpoint.
