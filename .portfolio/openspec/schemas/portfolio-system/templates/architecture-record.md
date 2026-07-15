# Architecture Record

## Problem Forces

- Domain complexity:
- Integration pressure:
- UI state complexity:
- Data or ML reproducibility:
- Auditability:
- Throughput or async pressure:
- Independent deployability:

## Decision

- Architecture:
- Stack profile:
- API style:
- Messaging:
- Cloud mode:
- Database/runtime:
- Library policy:

## Dependency Direction

State which layers may depend on which layers. Domain/use-case code must not depend on framework, database, broker, cloud SDK, transport, or UI.

## Principles Evidence

- SRP:
- OCP:
- LSP:
- ISP:
- DIP:
- KISS/YAGNI:

## Rejected Alternatives
