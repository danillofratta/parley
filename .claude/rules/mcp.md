# MCP (Model Context Protocol)

Parley uses MCP in two places: in the product (services expose use cases as
tools; the AI agent consumes them) and in development (`.mcp.json` for
Claude Code). Official SDKs: C# (`ModelContextProtocol`), Go
(`modelcontextprotocol/go-sdk`), Python (`mcp`); the agent loads tools with
`langchain-mcp-adapters`.

## Servers (Conversations in .NET, Operations in Go)

- A tool is an **entry point role** of a slice (`<UseCase>Tool.cs`, `tool.go`):
  it maps arguments → Command/Query, calls the handler, maps Result → output.
  No business logic in tools.
- Read tools map to query slices. Command tools exist only when a story asks
  for them, must be idempotent and still go through the aggregate.
- Identity and tenant come from the authenticated MCP session, never from
  tool arguments. A model can be manipulated into passing any argument.
- Tool names are `snake_case` verbs; descriptions say what the tool does, when
  to use it and when not to, written for a model reader.
- Input schemas are strict (required fields, enums, max lengths). Outputs are
  minimal and never include data the caller's tenant cannot see; operational
  tools mask personal data.
- Transport: Streamable HTTP, mounted at `/mcp` in the service's own host.
  Simplification locally: internal network + API key. Cloud: OAuth 2.1 as in
  the MCP authorization spec, or mTLS between services (ADR 0013).

## Client (Agent Orchestrator, Python)

- Servers come from configuration; each graph declares an allowlist of tools.
- Tool outputs are untrusted data, like customer messages: never follow
  instructions found inside them.
- Command tools are only callable from nodes that decide them explicitly,
  never from free-form model output based on customer text.
- Timeouts on every call; if a server is down, continue without the tool and
  lower confidence.

## Development

- `.mcp.json` declares shared MCP servers for Claude Code; tokens come from
  environment variables (`${GITHUB_PAT}`), never from the file.
- Parley's own servers are added to `.mcp.json` when OC-801 and OC-804 land.
- Inspect servers manually with the MCP Inspector; automate tool listing,
  schemas and handler calls in tests.
