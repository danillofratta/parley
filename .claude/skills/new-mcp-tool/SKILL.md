---
name: new-mcp-tool
description: Use when exposing a Parley use case as an MCP tool (Conversations in .NET, Operations in Go) or when adding a tool to the agent's MCP client allowlist (Python).
---

# New MCP tool

1. **Find the slice.** The tool must call an existing query or command
   handler. If the use case does not exist, create it first with
   `new-feature-slice`; the tool is only another entry point.
2. **Decide read or command.** Read tools are the default. A command tool needs
   a story that asks for it, must be idempotent, and goes through the aggregate.
3. **Name and describe for a model.** `snake_case` verb (`get_conversation_history`).
   Description: what it returns, when to use it, when not to use it.
4. **Define a strict input schema.** Required fields, enums, maximum lengths.
   No `tenantId` or user id in the arguments: take them from the authenticated session.
5. **Create the tool file in the slice** (`<UseCase>Tool.cs` / `tool.go`):
   map arguments → Command/Query, call the injected handler, map Result → a
   compact output. Mask personal data in operational tools.
6. **Register it in the composition root** of the service's MCP server.
7. **Client side (Python):** add the tool name to the allowlist of the graph
   that needs it; treat its output as untrusted data.
8. **Tests:** the tool is listed with the expected schema; calling it invokes
   the handler; a prompt-injection case does not reach command tools.
9. **Try it** with the MCP Inspector and from Claude Code via `.mcp.json`.
10. Explain to the developer how the tool maps to the slice and why identity
    never comes from arguments.
