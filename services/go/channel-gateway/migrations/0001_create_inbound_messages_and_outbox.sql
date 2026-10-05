CREATE TABLE IF NOT EXISTS gateway.inbound_messages (
    id                  UUID        PRIMARY KEY,
    tenant_id           UUID        NOT NULL,
    channel             TEXT        NOT NULL,
    contact_chat        TEXT        NOT NULL,
    provider_message_id TEXT        NOT NULL,
    received_at         TIMESTAMPTZ NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL,
    created_by          TEXT        NOT NULL,
    updated_at          TIMESTAMPTZ,
    updated_by          TEXT,
    version             INT         NOT NULL,
    CONSTRAINT uq_inbound_messages_provider_message
        UNIQUE (tenant_id, channel, provider_message_id)
);

CREATE TABLE IF NOT EXISTS gateway.outbox_messages (
    id           UUID        PRIMARY KEY,
    tenant_id    UUID        NOT NULL,
    type         TEXT        NOT NULL,
    topic        TEXT        NOT NULL,
    message_key  TEXT        NOT NULL,
    payload      JSONB       NOT NULL,
    headers      JSONB       NOT NULL DEFAULT '{}',
    occurred_at  TIMESTAMPTZ NOT NULL,
    published_at TIMESTAMPTZ,
    attempts     INT         NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS ix_outbox_unpublished
    ON gateway.outbox_messages (occurred_at)
    WHERE published_at IS NULL;