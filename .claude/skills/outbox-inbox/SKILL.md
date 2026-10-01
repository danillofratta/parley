---
name: outbox-inbox
description: Use when a service must publish or consume Kafka events. Implements the transactional Outbox, relay and Inbox deduplication in Go, .NET or Python.
---

# Outbox and Inbox

## Tables (per service schema)

```sql
CREATE TABLE <schema>.outbox_messages (
    id               UUID PRIMARY KEY,          -- becomes envelope.messageId
    tenant_id        UUID NOT NULL,
    type             TEXT NOT NULL,             -- e.g. MessageReceived.v1
    topic            TEXT NOT NULL,
    message_key      TEXT NOT NULL,             -- conversationKey
    payload          JSONB NOT NULL,            -- full envelope
    headers          JSONB NOT NULL DEFAULT '{}',
    occurred_at      TIMESTAMPTZ NOT NULL,
    published_at     TIMESTAMPTZ,
    attempts         INT NOT NULL DEFAULT 0
);
CREATE INDEX ix_outbox_unpublished ON <schema>.outbox_messages (occurred_at)
    WHERE published_at IS NULL;

CREATE TABLE <schema>.inbox_messages (
    message_id       UUID NOT NULL,
    consumer         TEXT NOT NULL,
    processed_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (message_id, consumer)
);
```

## Producer (command handler)

1. Change the aggregate.
2. Insert one outbox row per integration event.
3. Commit both in **one** transaction. Never call Kafka here.

## Relay

1. `SELECT ... WHERE published_at IS NULL ORDER BY occurred_at LIMIT n FOR UPDATE SKIP LOCKED`.
2. Publish each row with key = `message_key`, headers incl. `traceparent`.
   Use an idempotent producer (`enable.idempotence=true`, `acks=all`).
3. Mark `published_at` only after the broker acknowledges.
4. At-least-once: duplicates downstream are expected; the Inbox absorbs them.
5. Ordering: publish rows of the same `message_key` in `occurred_at` order.

## Consumer

1. Decode the envelope and validate against the JSON Schema.
2. In one transaction: `INSERT INTO inbox_messages ... ON CONFLICT DO NOTHING`;
   if no row was inserted, skip; otherwise run the command handler.
3. Commit the transaction, then commit the Kafka offset.
4. On failure: retry topic with backoff, then dead-letter topic.

## Stack notes

- Go and Python: implement by hand (learning goal).
- .NET: a library may provide Outbox/Inbox (Wolverine, MassTransit); check its
  license and record the choice in an ADR. Explain what it does underneath.
- Flag simplifications: polling relay vs CDC (Debezium) is ADR 0003.
