# Parley

Omnichannel customer service platform with AI agents. Public reference
architecture and learning project: every decision must be explainable.

## Architecture in one paragraph

Microservices in three stacks, all built with the same three patterns:
**DDD** (bounded contexts, aggregates, value objects, domain events),
**CQRS** (commands change state through aggregates; queries read projections)
and **Vertical Slice Architecture** (one folder per use case, from entry point
to persistence). Services talk through versioned Kafka events, with
Outbox on the producer side and Inbox on the consumer side. Use cases can also
be exposed as **MCP tools** (Conversations and Operations servers; the Agent
is an MCP client).

| Service | Stack | Bounded context | Path |
| --- | --- | --- | --- |
| Channel Gateway | Go | Channels (ingress) | `services/go/channel-gateway` |
| Outbound Dispatcher | Go | Channels (egress) | `services/go/outbound-dispatcher` |
| Conversations | .NET 10 | Conversations (core) | `services/dotnet/conversations` |
| Agent Orchestrator | Python 3.13 + LangGraph | Agent (core) | `services/python/agent-orchestrator` |

Detailed rules live in `.claude/rules/` and load automatically.
Domain language and context map:

@docs/domain/ubiquitous-language.md
@docs/domain/context-map.md

## Non-negotiables

- Every service uses the layout `domain/` + `features/<use-case>/` + `infrastructure/`.
- Each slice keeps one file per role: request, response, command/query, result,
  validator, handler, endpoint or consumer (see `.claude/rules/architecture.md`).
- No mediator: the endpoint, consumer or MCP tool calls its own handler directly.
- MCP tools are entry points to existing slices: no business logic in a tool,
  tenant taken from the authenticated caller, never from tool arguments.
- Dependencies are injected through constructors and wired in one composition
  root per service (.NET built-in DI, Python `dependency-injector`, Go manual wiring).
  Slices never resolve dependencies themselves (no service locator).
- `domain/` imports nothing from features, infrastructure or frameworks.
- Every entity and aggregate follows the domain model standard
  (`.claude/rules/domain-model.md`): `Id` (UUID v7), audit fields
  (`CreatedAt`, `CreatedBy`, `UpdatedAt`, `UpdatedBy`), `Version`, `TenantId` on
  aggregate roots, business rules checked with `CheckRule`.
- A slice never imports another slice.
- Commands go through an aggregate; queries never load aggregates and never write.
- One aggregate changed per transaction; cross-aggregate effects happen through events.
- Events leave a service only through its Outbox, in the same transaction as the state change.
- Every consumer deduplicates through its Inbox before handling.
- Kafka message key = `conversationKey`.
- `tenantId` on every event, table and query.
- No personal data in logs. No secrets in the repository.
- LLM provider only via `LLM_PROVIDER` / `LLM_MODEL` (default Groq); never import a provider SDK inside a slice.

## Commands

| Task | Command |
| --- | --- |
| Start infrastructure | `docker compose up -d` |
| Go tests (per service) | `go test ./...` |
| .NET tests | `dotnet test` |
| Python tests | `uv run pytest` |
| Python lint/format | `uv run ruff check . && uv run ruff format --check .` |

The developer works on Windows (cmd/PowerShell). Prefer cross-platform
commands; when shell syntax matters, show PowerShell and bash.

## How to work in this repo

- This is a learning project: explain the reasoning and trade-offs behind
  every change, and flag simplifications compared with production.
- Use the skills: `new-service`, `new-feature-slice`, `new-aggregate`,
  `outbox-inbox`, `new-integration-event`, `new-mcp-tool`, `write-adr`.
- MCP servers for development are declared in `.mcp.json` (tokens come from
  environment variables, e.g. `GITHUB_PAT`).
- Architecture decisions become ADRs in `docs/adr/`.
- Requirements and story IDs (OC-xxx): `docs/requirements.md`.
- Conventional Commits; small PRs that reference the Issue (`Closes #N`).
