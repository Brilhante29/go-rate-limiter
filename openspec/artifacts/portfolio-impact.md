# Portfolio Impact: go-rate-limiter

## Program

- Program id: `backend-reliability-platform`
- Program name: Backend Reliability and Architecture Platform
- Component pack: `backend-reliability-platform`

## System Story

Provides the shared traffic-control primitive used by gateways and multi-tenant services in the backend reliability program.

This repository is not a standalone demo. It is one part of the Backend Reliability and Architecture Platform system and should produce reusable fixtures, benchmark patterns, and decisions for later repositories.

## Proficiency Signal

- Primary profile: `go-backend`
- Stack profile: `go-backend`
- Stack:
- go-1.26
- chi-v5
- redis-8.8
- go-redis-v9
- lua
- k6
- docker

## Post Angle

Open with total_rps = 11511.03 requests_per_second, then explain why the architecture and local-first path make the result reproducible.
