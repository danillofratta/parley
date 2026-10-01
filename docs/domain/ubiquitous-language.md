# Ubiquitous language

Use these terms exactly, in code, events, tables and docs.

| Term | Meaning | Not to be confused with |
| --- | --- | --- |
| Tenant | A company using Parley. Every piece of data belongs to one tenant. | User, account |
| Channel | A medium customers use: `telegram`, `webchat`, `whatsapp`, `email`. | Topic, queue |
| Channel account | A tenant's connection to a channel (e.g. one Telegram bot). | Contact |
| Contact | The end customer, identified per channel. | User, support agent |
| Conversation key | `<channel>:<external chat id>`; identifies a conversation thread on a channel and is the Kafka message key. | Conversation id |
| Conversation | The ongoing exchange between one contact and one tenant. States: `open`, `awaiting_human`, `closed`. | Chat, ticket |
| Message | One unit of content in a conversation, `inbound` (from the contact) or `outbound` (to the contact). | Kafka message (say "event" or "record") |
| AI agent | The LangGraph workflow that proposes replies. | Support agent |
| Support agent | A human who takes over conversations. | AI agent |
| Reply proposal | A reply drafted by the AI agent, with a confidence score from 0 to 1. Never sent directly. | Reply |
| Reply approval | The Conversations decision that a proposal (or a support agent's text) will be sent. | — |
| Handoff | Moving a conversation to `awaiting_human`. | Escalation |
| Delivery | The act of sending an approved reply through a channel, with its outcome. | — |
| Inbound record | The raw provider payload stored once by the Channel Gateway. | Message |
