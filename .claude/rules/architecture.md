# Architecture: DDD + CQRS + Vertical Slices

These rules apply to every service, in every stack.

## Service layout

```
<service>/
  domain/            aggregates, entities, value objects, domain events,
                     domain errors, repository ports (interfaces)
  features/
    <use_case>/      one folder per use case (a "slice"), one file per role
  infrastructure/    technical adapters shared by slices: database, Kafka,
                     outbox/inbox, HTTP server, config, telemetry, LLM client,
                     CQRS handler interfaces
  bootstrap/         composition root: builds the dependency graph and routes
```

Naming per stack: Go uses lowercase packages (`features/receivewebhook`);
.NET uses `Features/<Context>/<UseCase>`; Python uses snake_case
(`features/propose_reply`). In .NET the composition root is `Program.cs`
plus `Features/FeatureRegistration.cs`.

## Dependency rule

```
bootstrap      ──▶ everything (wiring only, no logic)
features       ──▶ domain, infrastructure
infrastructure ──▶ domain   (implements domain ports)
domain         ──▶ nothing  (standard library only)
```

- A slice never imports another slice. Shared behaviour belongs in `domain/`
  (business) or `infrastructure/` (technical).
- Architecture tests enforce this rule in every service.

## Slice anatomy: one file per role

| Role | File (Go / .NET / Python) | When | Responsibility |
| --- | --- | --- | --- |
| Request | `request.go` / `<UseCase>Request.cs` / `request.py` | HTTP entry point with input | Transport model: the exact shape received over HTTP |
| Response | `response.go` / `<UseCase>Response.cs` / `response.py` | HTTP entry point with a body | Transport model: the exact shape returned over HTTP |
| Message | `message.go` / `<UseCase>Message.cs` / `message.py` | Kafka entry point | Payload of the consumed integration event |
| Command or Query | `command.go`·`query.go` / `<UseCase>Command.cs`·`<UseCase>Query.cs` / `command.py`·`query.py` | Always | Application input, independent of transport |
| Result | `result.go` / `<UseCase>Result.cs` / `result.py` | Always | Application output, independent of transport |
| Validator | `validator.go` / `<UseCase>Validator.cs` / `validator.py` | Input has rules beyond types | Boundary validation of the command or query |
| Handler | `handler.go` / `<UseCase>Handler.cs` / `handler.py` | Always | Orchestrates the use case |
| Endpoint | `endpoint.go` / `<UseCase>Endpoint.cs` / `endpoint.py` | HTTP entry point | Maps Request → Command, calls the handler, maps Result → Response |
| Consumer | `consumer.go` / `<UseCase>Consumer.cs` / `consumer.py` | Kafka entry point | Inbox check, maps Message → Command, calls the handler |
| MCP tool | `tool.go` / `<UseCase>Tool.cs` / `tool.py` | Use case exposed to MCP clients | Maps tool arguments → Command/Query, calls the handler, maps Result → tool output |
| Ports | `ports.go` / — / `ports.py` | Slice-specific dependencies | Interfaces only this slice needs (a clock, an external API) |
| Registration | — / `<UseCase>Registration.cs` / — | .NET | Registers the slice's services in DI and maps its endpoint |

Why Request/Response and Command/Result are separate: the transport can
change (HTTP today, Kafka or gRPC tomorrow) without touching the use case,
and several entry points (an endpoint, a consumer, an MCP tool) can feed the
same command or query.
The cost is more files and explicit mapping; we accept it for clarity.

## No mediator

- The endpoint, consumer or MCP tool receives its handler through constructor or
  parameter injection and calls it directly. There is no `Send()`/dispatcher.
- Handlers implement small interfaces (`ICommandHandler<TCommand, TResult>`,
  `IQueryHandler<TQuery, TResult>` in .NET; `Handle(ctx, cmd)` methods in Go;
  `async def handle(...)` in Python) so they can be decorated and faked in tests.
- Cross-cutting concerns (logging, tracing, transactions) use middleware,
  endpoint filters, decorators or the unit of work, never hidden pipelines.

## Dependency injection

- Constructor injection only. A class/struct declares everything it needs.
- One composition root per service (`bootstrap/` or `Program.cs`). Only the
  composition root knows concrete implementations.
- Slices and the domain never import the container and never resolve services
  at runtime (no service locator).
- Lifetimes (.NET, Python): handlers and repositories are scoped per request or
  per consumed message; clients such as the Kafka producer, HTTP clients and
  the LLM model are singletons.

## CQRS (logical, same database)

- **Commands** change state. The handler validates, loads an aggregate through
  a repository port, calls a domain method, saves it, and writes the resulting
  integration events to the Outbox in the same transaction. It returns a Result
  with at most the ids it created or flags about what happened.
- **Queries** never change state and never load aggregates. The handler reads
  projections straight from the database into a Result shaped for the caller.
- Simplification: commands and queries share one database; there is no
  separate read store and no event sourcing. A separate read model is only
  introduced by an ADR when a query needs it.
- Consumers contain no business logic: Inbox check, map, call the handler.

## Tactical DDD

- Aggregates are small, protect their invariants and expose behaviour, not
  setters. Create them through factory methods.
- Modify one aggregate per transaction. Reference other aggregates by id.
- Value objects are immutable and validate themselves on creation
  (`ConversationKey`, `TenantId`, `Confidence`).
- Aggregates record domain events (past tense: `ReplyApproved`). Handlers or
  repositories map them to integration events (`contracts/events/`). Domain
  events never leave the service as-is.
- Repository ports live in `domain/`, one per aggregate root; implementations
  in `infrastructure/`.
- Validators check the input at the boundary (required, length, format);
  business invariants live in value objects and aggregates.
- Domain errors are typed; the entry point maps them to HTTP status codes or
  retry/dead-letter decisions.

## Depth per bounded context

See `docs/domain/context-map.md`. Conversations and Agent use full tactical
DDD. Channels keeps a small domain; do not invent aggregates without
invariants to protect.
