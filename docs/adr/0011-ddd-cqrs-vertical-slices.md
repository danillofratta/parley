# ADR-0011: DDD + CQRS + Vertical Slices in every service

- Status: Accepted
- Date: 2026-09-30

## Context

Parley has four services in three stacks and is a public reference
architecture. Readers must find the same structure in every service, and the
structure must make business rules easy to locate and hard to bypass.
The domain has real invariants (a conversation's state machine, reply
approval, handoff), especially in the Conversations and Agent contexts.

## Options

### Layered / Clean Architecture (horizontal layers)
Familiar and well documented. A single use case is spread across controllers,
services, repositories and DTO folders; changes touch many layers and
generic "service" classes grow without bounds.

### Transaction scripts per endpoint
Simple and fast to write. Business rules end up duplicated across scripts and
nothing protects invariants.

### DDD + CQRS + Vertical Slices (chosen)
One folder per use case keeps a change local. Aggregates protect invariants.
Commands and queries evolve independently. More concepts to learn, and a risk
of ceremony where the domain is thin.

### Full CQRS with a separate read store or event sourcing
Maximum scalability and auditability, at the cost of eventual consistency
inside a service and much more infrastructure. Not justified by current needs.

## Decision

Every service uses `domain/` + `features/<use_case>/` + `infrastructure/`,
logical CQRS on a single database, and tactical DDD whose depth follows the
subdomain type (full in Conversations and Agent, light in Channels).
Each slice keeps one file per role (request, response, command/query,
result, validator, handler, endpoint/consumer). There is no mediator: the
entry point calls its handler directly. Dependency injection is ADR 0012.

## Consequences

- Positive: consistent structure across Go, .NET and Python; use cases are
  easy to find, review and test; invariants live in one place.
- Without a mediator there is no generic pipeline; cross-cutting concerns use
  middleware, endpoint filters, decorators and the unit of work, which keeps
  the call path explicit and easy to debug.
- Negative: more files per feature and explicit mapping between transport and
  application models; developers must learn aggregate design;
  some duplication between slices is accepted on purpose.
- Enforcement: architecture tests in every service.
- Simplifications compared with production: no separate read store and no
  event sourcing; introduced only by a new ADR if a query requires it.
