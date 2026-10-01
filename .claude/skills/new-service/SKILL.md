---
name: new-service
description: Use when creating a new Parley microservice in Go, .NET or Python. Scaffolds the DDD + CQRS + vertical slice layout, config, health check, database schema and tests.
---

# New service

1. Confirm the bounded context in `docs/domain/context-map.md`. A new context
   needs an ADR.
2. Create the layout from `.claude/rules/architecture.md`:
   `domain/`, `features/`, `infrastructure/`, `bootstrap/` (stack naming from
   the stack rule).
3. Create the composition root: Go `internal/bootstrap` with manual wiring;
   .NET `Program.cs` + `Features/FeatureRegistration.cs` with built-in DI;
   Python `bootstrap/container.py` with `dependency-injector`.
4. Add configuration from environment variables with safe local defaults,
   except secrets (fail closed).
5. Add `GET /healthz` checking the service's own dependencies.
6. Create the service's PostgreSQL user and schema in
   `deploy/compose/postgres/init/` and its migrations in the service folder.
7. Add `outbox_messages` / `inbox_messages` if it publishes or consumes (skill `outbox-inbox`).
8. Add an architecture test enforcing the dependency rule.
9. Add a `README.md` for the service: purpose, slices, events in/out, how to run.
10. Add the service to the table in `CLAUDE.md` and to `docker-compose.yml`.
