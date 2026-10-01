# Messaging (Kafka, Outbox, Inbox)

- Kafka message key = `conversationKey`, always. It guarantees order per conversation.
- Every event uses the common envelope (`messageId`, `type`, `tenantId`,
  `conversationKey`, `occurredAt`, `traceparent`, `payload`) and is validated
  against its JSON Schema in `contracts/events/`.
- Event types are versioned: `MessageReceived.v1`. Breaking changes create `.v2`.
- **Producer:** never call Kafka from a handler. Write the event to the
  service's `outbox_messages` table in the same transaction as the state
  change; a relay publishes it.
- **Consumer:** insert `messageId` into `inbox_messages` (unique constraint)
  in the same transaction as the handling. If it already exists, skip.
- Commit the Kafka offset only after that transaction commits.
- Failures: retry topic `<topic>.<consumer>.retry` with backoff, then
  dead-letter topic `<topic>.<consumer>.dlt`. A failed message must not let
  later messages of the same conversation overtake it.
- Propagate W3C `traceparent` as a Kafka header.
- Topics are created explicitly (auto-create is disabled).
