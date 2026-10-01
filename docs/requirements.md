# Parley — Release 1 Requirements

Sep 30, 2026 · @Danillo Fratta

## Product vision

Release 1 delivers one end-to-end flow: a customer sends a Telegram message and gets an AI-generated reply, with reliable Kafka messaging between Go, .NET and Python.

Parley is an omnichannel customer service platform: a multi-tenant system where companies serve customers across channels (Telegram, webchat, WhatsApp, email). AI agents answer from company knowledge, run simple actions and hand the conversation to a human when confidence is low.

Release 1 proves the architecture, not the feature set: duplicate webhooks never produce duplicate replies, messages in a conversation are processed in order, and no reply is lost when a service goes down. It becomes the first article and the first tag (v0.1.0). Release 1.1 (v0.2.0) runs the same system on AWS and Azure.

## Release 1 scope

Release 1 includes only what the Telegram → AI → Telegram flow needs to work with delivery guarantees, running locally. Release 1.1 deploys that same scope to AWS and Azure.

| Item | Release 1 (local) | Release 1.1 (cloud) | Later |
| --- | --- | --- | --- |
| Channels | Telegram | Telegram | Webchat (R2), WhatsApp and email (R3+) |
| Channel Gateway (Go) | Webhook, validation, Inbox, publishing | Public HTTPS endpoint | Adapters for other channels |
| Outbound Dispatcher (Go) | Sending with retry topics and Inbox | Same | Per-channel rate limits, delivery receipts |
| Conversations (.NET) | Contacts, conversations, messages, Outbox/Inbox | Same | Light CRM, search, full history |
| Agent Orchestrator (Python) | Simple LangGraph graph: intent + reply | Same | RAG, tools, handoff with summary |
| Human handoff | Conversation flagged as `awaiting_human` | Same | Real-time agent console (Realtime Hub in Go) |
| Multi-tenancy | One fixed tenant, but `tenantId` on every event and table | Same | Tenant onboarding, plans, billing |
| Identity | API key on the admin API | Keys in AWS Secrets Manager / Azure Key Vault | Keycloak, agent and supervisor roles |
| Messaging | Kafka (KRaft) in Docker Compose | Amazon MSK / Azure Event Hubs (Kafka endpoint) | Schema Registry |
| Data | PostgreSQL in Docker Compose | Amazon RDS / Azure Database for PostgreSQL | pgvector, Redis |
| Runtime | Docker Compose | EKS and AKS with shared Helm charts, provisioned by Terraform | Autoscaling, multi-region |
| Observability | OpenTelemetry with end-to-end trace | OTel Collector in each cluster | Dashboards, alerts, Langfuse |

Out of scope for Release 1 and 1.1: agent UI, RAG, WhatsApp, SLAs, billing and the LLM Gateway.

## Personas

Four personas cover Release 1; the last one is also whoever reads and reuses the repository.

| Persona | Who | What they want |
| --- | --- | --- |
| End customer | Person contacting the company on Telegram | A fast, correct answer without repeating themselves |
| Support agent | Employee who takes over hard conversations | To see which conversations need them, with context |
| Tenant admin | Owner of the support operation | To configure the channel and review what the AI answered |
| Developer / operator | Maintains the platform (and reads the repo) | One-command startup; trace a message end to end; deploy to any cloud |

## Epics and user stories

There are 9 epics and 25 stories; each story becomes a GitHub Issue with its ID in the title (e.g. `OC-102`). Acceptance criteria use Given / When / Then so they map directly to tests.

| Epic | Stack | Stories | Release |
| --- | --- | --- | --- |
| E0 Repository foundation | All | OC-001 to OC-003 | 1 |
| E1 Inbound messages | Go | OC-101 to OC-103 | 1 |
| E2 Conversations | .NET | OC-201 to OC-203 | 1 |
| E3 AI agent | Python | OC-301 to OC-302 | 1 |
| E4 Outbound replies | Go | OC-401 to OC-402 | 1 |
| E5 Minimal human handoff | .NET + Python | OC-501 | 1 |
| E6 Observability and operations | All | OC-601 to OC-602 | 1 |
| E7 Cloud deployment (AWS and Azure) | Terraform, Helm, GitHub Actions | OC-701 to OC-705 | 1.1 |
| E8 MCP integration | .NET (server), Python (client), Go (server) | OC-801 to OC-804 | 1.2 |

### E0 Repository foundation

**OC-001 Start the stack with one command.** As a developer, I want `docker compose up` to start Kafka, databases and services, so I can contribute within minutes.

- Given a clean clone and a `.env` copied from `.env.example`, when I run `docker compose up`, then every service is healthy within 3 minutes.
- No real secret exists in the repository.

**OC-002 Per-service CI pipeline.** As a developer, I want each PR to build, lint and test only the services it changes, so feedback is fast.

- Given a PR touching only `services/go/*`, when CI runs, then only the Go jobs execute.
- Merging to `main` requires green CI and one human approval.

**OC-003 AI context in the repository.** As a developer, I want `CLAUDE.md`, rules, skills and subagents versioned, so Claude Code and automated review follow project conventions.

- The skills `new-service`, `new-feature-slice`, `new-aggregate`, `outbox-inbox`, `new-integration-event`, `new-mcp-tool` and `write-adr` exist.
- `.mcp.json` declares the MCP servers used during development; no token is stored in it.
- Automated PR review runs read-only and is triggered by the `claude-review` label.

### E1 Inbound messages (Channel Gateway, Go)

**OC-101 Receive Telegram webhooks.** As an end customer, I want my message received even under traffic spikes, so I am never left without an answer.

- Given a valid webhook, when it arrives, then the gateway stores it and returns HTTP 200 in under 200 ms (p95).
- Given a webhook without the correct secret token, when it arrives, then it is rejected with 401 and nothing is stored.

**OC-102 Deduplicate webhooks (Inbox).** As an operator, I want a redelivered webhook to create no duplicate, so the customer never gets two replies.

- Given a webhook already received with the same provider id, when it arrives again, then the gateway returns 200 and publishes no new event.
- Deduplication relies on a database unique constraint, not memory.

**OC-103 Publish a normalized message.** As a developer, I want every channel's message to become the same `MessageReceived` event, so downstream services are channel-agnostic.

- The event carries `messageId`, `tenantId`, `conversationKey`, `channel`, `text`, `receivedAt` and `traceparent`.
- The Kafka message key is `conversationKey`.
- Publishing uses an Outbox: if Kafka is down, the event is sent once it recovers.

### E2 Conversations (Conversations Service, .NET)

**OC-201 Register contact and conversation.** As a support agent, I want each message linked to the right contact and conversation, so I have the full history.

- Given a `MessageReceived` from a new contact, when consumed, then the contact, the conversation (status `open`) and the message are created.
- Given an existing contact with an open conversation, then the message joins that conversation.

**OC-202 Process in order per conversation.** As an end customer, I want my messages handled in the order I sent them, so the AI understands the context.

- Given three messages from one conversation, when consumed by a consumer group with several instances, then they are stored in the original order.
- Different conversations are processed in parallel across partitions.
- A failed message does not let later messages of the same conversation overtake it.

**OC-203 Consume with Inbox, publish with Outbox.** As an operator, I want no duplicated messages and no lost events, so state stays consistent.

- Given the same event delivered twice, when consumed, then state changes exactly once.
- Saving the message and writing `ConversationUpdated` to the Outbox happen in one transaction.

### E3 AI agent (Agent Orchestrator, Python)

**OC-301 Reply with a LangGraph graph.** As an end customer, I want a useful answer to my question, so I solve my problem without waiting for a human.

- Given a `ConversationUpdated` with a customer message, when the graph runs, then it classifies intent, drafts a reply and publishes `ReplyProposed` with a confidence score (0 to 1).
- The graph receives the last 10 messages of the conversation as context.
- Inbox and Outbox are implemented by hand on PostgreSQL.

**OC-302 Fail safely.** As an operator, I want an LLM error or timeout not to block the topic, so other conversations keep flowing.

- Given an LLM provider timeout, when 3 attempts fail, then the message goes to the dead-letter topic and `HandoffRequested` is published with reason `ai_error`.

### E4 Outbound replies (Outbound Dispatcher, Go)

**OC-401 Send the reply to the channel.** As an end customer, I want the answer on the same channel I used.

- Given a `ReplyApproved`, when consumed, then the message is sent through the Telegram API and `MessageSent` is published.
- Given the same `ReplyApproved` delivered twice, then the message is sent only once.

**OC-402 Retry with backoff.** As an operator, I want transient Telegram failures retried, so replies are not lost.

- Given a 429 or 5xx, then sending is retried through retry topics with exponential backoff, up to 5 times, honoring the provider's retry-after.
- After the last attempt, the message goes to the dead-letter topic and `MessageSendFailed` is published.

### E5 Minimal human handoff

**OC-501 Flag a conversation for a human.** As a support agent, I want to see which conversations the AI could not resolve, so I can take them over.

- Given a `ReplyProposed` below the configured confidence threshold (default 0.7), then the reply is not sent and the conversation moves to `awaiting_human`.
- The admin API lists `awaiting_human` conversations and accepts a manual reply.

### E6 Observability and operations

**OC-601 Trace a message end to end.** As an operator, I want to follow a message from webhook to delivery in one trace, to debug across Go, .NET and Python.

- Given a `traceparent` Kafka header, when the message crosses all four services, then the full trace appears in Jaeger (or Grafana Tempo).

**OC-602 Test failure scenarios.** As a developer, I want automated tests for duplicates, ordering and broker outages, to prove the architecture's guarantees.

- Integration tests (Testcontainers) cover: duplicate webhook, out-of-order delivery, Kafka unavailable during publishing, consumer restarted mid-processing and partition rebalance.

### E7 Cloud deployment (AWS and Azure)

**OC-701 Infrastructure as code.** As an operator, I want each cloud environment defined in Terraform, so it can be created and destroyed on demand.

- Given `terraform apply` for `envs/aws` or `envs/azure`, then network, Kubernetes cluster, managed Kafka, PostgreSQL, secret store and container registry are created.
- Given `terraform destroy`, then no billable resource remains.

**OC-702 Deploy to AWS.** As an operator, I want the platform running on AWS, to validate a real deployment.

- Services run on EKS using the shared Helm charts; Kafka is Amazon MSK; databases are Amazon RDS for PostgreSQL; secrets come from AWS Secrets Manager.
- The Telegram webhook reaches the gateway over HTTPS and the end-to-end flow works.

**OC-703 Deploy to Azure.** As an operator, I want the same platform on Azure with no application code changes, to prove portability.

- Services run on AKS with the same Helm charts and only different values files; Kafka clients use the Azure Event Hubs Kafka endpoint; databases are Azure Database for PostgreSQL; secrets come from Azure Key Vault.
- The Kafka compatibility tests (ordering, consumer groups, retry topics) pass on Event Hubs, or the gaps are documented in an ADR.

**OC-704 Continuous delivery without long-lived keys.** As a developer, I want GitHub Actions to deploy to both clouds using OIDC federation, so no cloud credential is stored in the repository.

- Given a release tag, when the CD workflow runs, then images are pushed and deployed to the target cloud after a manual approval (GitHub Environments).

**OC-705 Cost guardrails.** As the project owner, I want budgets and automatic teardown, so test environments do not generate surprise bills.

- A budget alert exists in each cloud account.
- A scheduled workflow destroys test environments that are idle beyond a configured window.

### E8 MCP integration (Release 1.2)

MCP tools are one more entry point to existing slices, exactly like an HTTP endpoint or a Kafka consumer; they add no business logic of their own.

**OC-801 Conversations MCP server (.NET).** As the AI agent, I want Conversations' read use cases exposed as MCP tools, so any MCP client can fetch conversation context through a standard protocol.

- Tools `get_conversation_history` and `list_awaiting_human_conversations` call the existing query handlers.
- Streamable HTTP transport; the tenant comes from the authenticated caller, never from tool arguments.
- Tool names, descriptions and input schemas are covered by tests.

**OC-802 Request handoff through MCP.** As the AI agent, I want a `request_handoff` tool, so the graph can hand a conversation to a human when it decides to.

- The tool calls the `RequestHandoff` command handler; Conversations still enforces its invariants (a closed conversation rejects it).
- Calling it twice for the same conversation has the same effect as calling it once.

**OC-803 Agent as MCP client (Python).** As a developer, I want the LangGraph agent to load its tools from MCP servers, so tools can be added without changing the agent code.

- Tools are loaded with `langchain-mcp-adapters` from servers listed in configuration, filtered by an allowlist per graph.
- If an MCP server is unavailable, the graph continues without tools and lowers its confidence, which leads to a handoff.
- Prompt-injection test: a customer message telling the agent to call a tool does not trigger a command tool.

**OC-804 Operations MCP server (Go).** As an operator, I want to inspect inbound records and dead-lettered messages from an MCP client (Claude Desktop or Claude Code), to debug without writing SQL.

- Read-only tools `find_inbound_records` (by conversation key) and `list_dead_letters`; personal data masked in the output.
- Enabled only outside production or behind admin authentication.

## Non-functional requirements

These guarantees are the core of Release 1; each one has at least one automated test that proves it.

| Area | Requirement | How to verify |
| --- | --- | --- |
| Delivery | At-least-once between services; exactly-once effect through the Inbox | Duplicate event test (OC-602) |
| Idempotency | Every consumer deduplicates by `messageId` with a database unique constraint | Redelivery test |
| Ordering | Kafka key = `conversationKey`, so order holds per conversation; parallelism across partitions | Three-message sequence test with several consumer instances |
| Offsets | Offsets committed only after the Inbox transaction commits | Consumer restart test |
| Atomicity | State + event in one transaction (Outbox) in every producer | Kafka-unavailable test |
| Resilience | Retry topics with backoff, one dead-letter topic per consumer, timeouts on every external call | Telegram and LLM failure tests |
| Latency | Webhook answered in < 200 ms (p95); customer reply in < 10 s (p95) | Load test with k6 |
| Security | Webhook secret token validated; secrets only in environment variables or the cloud secret store | Invalid webhook test; secret scanning in CI |
| LGPD / privacy | No personal data in logs; phone and name masked in traces | Code review + log test |
| Multi-tenancy | `tenantId` required on every event, table and query | Contract test |
| Observability | W3C trace context propagated in Kafka headers across Go, .NET and Python | Full trace in Jaeger (OC-601) |
| Contracts | Events versioned (`v1`) and validated with JSON Schema by every consumer | Contract tests in CI |
| Portability | Same images and Helm charts on EKS and AKS; cloud differences only in Terraform and values files | OC-703 deploys with no code change |

## Events and topics

Seven events cover Release 1; all share one envelope and use `conversationKey` as the Kafka message key.

| Event | Topic | Producer | Consumers | When |
| --- | --- | --- | --- | --- |
| `MessageReceived.v1` | `messages.inbound` | Channel Gateway (Go) | Conversations (.NET) | Valid, new webhook received |
| `ConversationUpdated.v1` | `conversations.events` | Conversations (.NET) | Agent Orchestrator (Python) | Customer message stored |
| `ReplyProposed.v1` | `agent.replies` | Agent Orchestrator (Python) | Conversations (.NET) | Graph produced a reply with confidence |
| `ReplyApproved.v1` | `messages.outbound` | Conversations (.NET) | Outbound Dispatcher (Go) | Confidence above threshold, or manual agent reply |
| `HandoffRequested.v1` | `conversations.handoff` | Agent Orchestrator (Python) or Conversations (.NET) | Conversations (.NET) | Low confidence or AI error |
| `MessageSent.v1` | `messages.delivery` | Outbound Dispatcher (Go) | Conversations (.NET) | Channel confirmed delivery |
| `MessageSendFailed.v1` | `messages.delivery` | Outbound Dispatcher (Go) | Conversations (.NET) | Attempts exhausted |

Each consumer also owns `<topic>.<consumer>.retry` and `<topic>.<consumer>.dlt` topics. The approve-or-handoff decision lives in .NET, which owns the conversation; Python only proposes. Business rules stay in one place and the agent stays replaceable.

Common envelope (the payload travels as the Kafka value; `traceparent` also travels as a Kafka header):

```json
{
  "messageId": "uuid",
  "type": "MessageReceived.v1",
  "tenantId": "uuid",
  "conversationKey": "telegram:123456",
  "occurredAt": "2026-10-01T12:00:00Z",
  "traceparent": "00-...",
  "payload": { }
}
```

JSON Schemas live in `contracts/events/` and are the single source of truth for Go, .NET and Python.

## Definition of Done

A story closes only when every item below is checked in the PR that implements it.

- [ ] Acceptance criteria covered by automated tests
- [ ] Green CI (build, lint, tests, secret scanning)
- [ ] AI review done and comments addressed or justified
- [ ] Approved by a human
- [ ] Event contract updated in `contracts/` when it changes
- [ ] Logs free of personal data and trace context propagated
- [ ] Helm chart and values updated when configuration changes
- [ ] Service README or docs updated
- [ ] ADR written if an architecture decision was made
- [ ] PR references the Issue (`Closes #N`) and follows Conventional Commits

## ADRs and open decisions

Thirteen ADRs cover Releases 1 to 1.2; ADRs 0011 to 0013 are already accepted, and each ADR is also an article topic.

| ADR | Decision | Proposal or outcome |
| --- | --- | --- |
| 0001 | Polyglot architecture: role of Go, .NET and Python | Go at the edge (I/O), .NET for the domain, Python for AI |
| 0002 | Message broker | Kafka (decided): KRaft locally, managed service in each cloud |
| 0003 | Outbox relay and Inbox per stack | Compare polling relay vs Debezium CDC; .NET library (Wolverine or MassTransit, check license) vs hand-written in Go and Python |
| 0004 | Per-conversation ordering and retries | Key = `conversationKey`; retry topics without letting a conversation's later messages overtake a failed one |
| 0005 | Event contracts | Versioned JSON Schema in `contracts/`; Schema Registry evaluated later |
| 0006 | Database per service | PostgreSQL, one schema per service locally, one database per service in the cloud |
| 0007 | LLM provider | Groq by default; provider and model set by `LLM_PROVIDER` and `LLM_MODEL` through LangChain's chat model abstraction |
| 0008 | Monorepo | One repository, CI by changed path |
| 0009 | Multi-cloud strategy | Kubernetes + shared Helm charts; Terraform modules per cloud; only values differ |
| 0010 | Managed Kafka per cloud | Amazon MSK on AWS; Azure Event Hubs Kafka endpoint on Azure, validated against OC-703 tests; fallback: Strimzi on AKS |
| 0011 | Architecture style | Accepted: DDD + CQRS + vertical slices, one file per role, no mediator |
| 0012 | Dependency injection | Accepted: .NET built-in DI, Python `dependency-injector`, Go manual composition root |
| 0013 | MCP | Accepted: MCP tools are slice entry points; Conversations (C#) and Operations (Go) servers; Agent (Python) client; MCP servers for development in `.mcp.json` |

Decided on 2026-09-30: repository and documentation in English; Kafka as the broker; real deployments on AWS and Azure; Groq as the default LLM provider, kept swappable; project name Parley.

Still open:

- [ ] Monthly LLM cost cap for the public demo
- [ ] Cloud budget per month for test environments

## From this document to GitHub

Each epic becomes a label, each release a Milestone, each story an Issue, and a public GitHub Project board shows progress.

1. Create labels: `epic:E0` to `epic:E8`, `stack:go`, `stack:dotnet`, `stack:python`, `stack:infra`, `type:story`, `claude-review`.
2. Create Milestones `v0.1.0 – Release 1 (local)` and `v0.2.0 – Release 1.1 (AWS + Azure)`.
3. Add an Issue template at `.github/ISSUE_TEMPLATE/user-story.md` with: story, acceptance criteria and the Definition of Done checklist.
4. Create one Issue per story with its ID in the title (e.g. `OC-102 Deduplicate webhooks (Inbox)`).
5. Create the GitHub Project with columns Backlog, Ready, In progress, In review and Done, plus an Epic field.
6. Suggested order: OC-001 → OC-002 → OC-003 → OC-101 → OC-102 → OC-103 → OC-201 → OC-203 → OC-202 → OC-301 → OC-401 → OC-501 → OC-302 → OC-402 → OC-601 → OC-602, then OC-701 → OC-704 → OC-702 → OC-703 → OC-705, then OC-801 → OC-802 → OC-803 → OC-804.

Issues can be created in bulk with the GitHub CLI (`gh issue create`) or by asking Claude Code to create them from this document.
