# Channel Gateway

Go service at the edge of the **Channels** bounded context. It receives
messages from channel providers (Telegram in Release 1), accepts each message
exactly once and publishes `MessageReceived.v1` through its Outbox.

## Slices

| Slice | Entry point | Story |
| --- | --- | --- |
| `receivetelegramupdate` | `POST /webhooks/telegram` | OC-101, OC-102 |

## Events

| Direction | Event | Topic | Key |
| --- | --- | --- | --- |
| Out | `MessageReceived.v1` | `messages.inbound` | `conversationKey` (`<channel>:<chat>`) |

## Outbox relay

Runs inside the gateway process. Each cycle takes a PostgreSQL advisory lock
(one active publisher across replicas, preserving per-conversation order),
reads pending rows in `occurred_at` order, publishes them to Kafka with the
conversation key and marks `published_at`. It stops at the first failure and
retries on the next cycle (at-least-once).

## Layout

```
cmd/gateway/            entry point
internal/domain/        seedwork, enums, rules, valueobjects, events, entities, repositories
internal/features/      one folder per use case
internal/infrastructure/ config, clock, persistence (Auditor), outbox (mapping + relay), kafka, postgres
internal/bootstrap/     composition root
migrations/             schema "gateway"
```

## Run locally

1. From the repository root: `docker compose up -d`
2. Apply the migration (cmd):
   `docker exec -i parley-postgres psql -U gateway -d parley < migrations\0001_create_inbound_messages_and_outbox.sql`
3. Set `TELEGRAM_WEBHOOK_SECRET` and run `go run ./cmd/gateway`
4. Send `testdata/telegram_update.json` to `POST http://localhost:8081/webhooks/telegram`
   with header `X-Telegram-Bot-Api-Secret-Token`.

| Variable | Default |
| --- | --- |
| `HTTP_ADDR` | `:8081` |
| `GATEWAY_DATABASE_URL` | `postgres://gateway:gateway@localhost:5432/parley?sslmode=disable` |
| `TELEGRAM_WEBHOOK_SECRET` | required |
| `DEFAULT_TENANT_ID` | `00000000-0000-0000-0000-000000000001` |
| `KAFKA_BOOTSTRAP_SERVERS` | `localhost:9092` |
| `OUTBOX_POLL_INTERVAL` | `500ms` |
| `OUTBOX_BATCH_SIZE` | `100` |
| `OUTBOX_PUBLISH_TIMEOUT` | `10s` |

## Tests

`go vet ./...` and `go test ./...` (unit tests need no database).
