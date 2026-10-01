---
name: new-integration-event
description: Use when creating or changing an event that crosses service boundaries. Keeps contracts in contracts/events versioned and consistent across Go, .NET and Python.
---

# New integration event

1. Name it in past tense with the ubiquitous language: `ReplyApproved.v1`.
2. Create `contracts/events/<Name>.v1.schema.json` (JSON Schema draft 2020-12)
   with the common envelope and a strict `payload` (`additionalProperties: false`).
3. Add an example: `contracts/events/examples/<Name>.v1.json`.
4. Update the events table in `docs/requirements.md` (topic, producer, consumers).
5. Breaking change (removed/renamed/retyped field) = new version `.v2`;
   producers publish both until every consumer migrates.
6. Adding an optional field is non-breaking; consumers must ignore unknown fields.
7. Producer and consumers get a contract test that validates against the schema.
8. Never put personal data in event names or keys; payload fields holding
   personal data are listed in the schema `description`.
