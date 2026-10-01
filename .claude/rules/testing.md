# Testing

| Level | What | How |
| --- | --- | --- |
| Domain | Aggregates and value objects | Pure unit tests, no mocks, no I/O |
| Slice | Handler behaviour | Fakes for ports, or Testcontainers for persistence |
| Integration | Inbox/Outbox, SQL, Kafka | Testcontainers (PostgreSQL, Kafka) |
| Contract | Events match `contracts/events/` schemas | Schema validation in each consumer and producer |
| Architecture | Dependency rule | .NET: NetArchTest/ArchUnitNET; Python: import-linter; Go: import checks in a test |

- Every acceptance criterion of a story (OC-xxx) maps to at least one test.
- Test names describe behaviour: `DuplicateWebhookIsStoredOnce`.
- A fake proves the handler's behaviour, not the database guarantee.
  Guarantees (dedup, ordering, atomicity) need integration tests.
