# Context map

| Bounded context | Services | Subdomain type | Owns |
| --- | --- | --- | --- |
| Channels | Channel Gateway, Outbound Dispatcher | Supporting | Provider protocols, inbound records, deliveries |
| Conversations | Conversations | Core | Contacts, conversations, messages, approvals, handoffs |
| Agent | Agent Orchestrator | Core | Intent classification, reply proposals, LLM usage |

Relationships:

- **Telegram → Channels: Anticorruption Layer.** Provider payloads are translated
  into Parley's language at the edge; no provider type leaves the Channels context.
- **Channels → Conversations: Published Language.** The contract is the versioned
  events in `contracts/events/`.
- **Conversations → Agent: Customer/Supplier.** Conversations publishes
  `ConversationUpdated`; the Agent consumes it.
- **Agent → Conversations: proposals only.** The Agent publishes `ReplyProposed`;
  only Conversations decides approval or handoff. The Agent never owns conversation state.
- **Conversations as an Open Host Service over MCP.** Conversations exposes
  read tools (history, awaiting-human list) and one command tool
  (`request_handoff`) through an MCP server. The Agent consumes them as an MCP
  client. Conversations still enforces every invariant.
- **Operations MCP (Channels).** The Channel Gateway exposes read-only
  operational tools for support engineers through an MCP server.

Tactical DDD depth follows the subdomain: Conversations and Agent get full
aggregates and domain events; Channels has a small domain (value objects and a
few entities) because its logic is mostly protocol translation. The folder
layout and CQRS rules are identical everywhere.
