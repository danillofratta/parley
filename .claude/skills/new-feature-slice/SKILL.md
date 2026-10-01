---
name: new-feature-slice
description: Use when adding a use case (command, query or message consumer) to any Parley service. Creates a vertical slice with one file per role (request, response, command/query, result, validator, handler, endpoint/consumer), no mediator, dependencies injected through the composition root. Go, .NET or Python.
---

# New feature slice

## 1. Classify the use case

- **Command**: changes state ("register", "approve", "send").
- **Query**: reads state ("list", "get").
- **Entry point**: HTTP endpoint (Request/Response) or Kafka consumer (Message).

Name it after the intent, using the ubiquitous language
(`docs/domain/ubiquitous-language.md`). Link the story ID (OC-xxx).

## 2. Check the domain first (commands only)

- Which aggregate owns the change? If none exists, run `new-aggregate`.
- Which invariant does the domain method protect?
- Which domain event does it record, and which integration event leaves the
  service? New events go through `new-integration-event`.

## 3. Create one file per role

Use the table "Slice anatomy" in `.claude/rules/architecture.md`.
Create only the roles the use case needs (a webhook has no Response body;
a consumer has a Message instead of Request/Response).

Stack templates: [go.md](go.md), [dotnet.md](dotnet.md), [python.md](python.md).

## 4. Wire it in the composition root

- Go: construct handler and endpoint in `internal/bootstrap`, register the route.
- .NET: `Add<UseCase>` / `Map<UseCase>` in the Registration file, listed in
  `Features/FeatureRegistration.cs`.
- Python: add providers to `bootstrap/container.py`.

## 5. Checklist before finishing

- [ ] The entry point calls its handler directly (no mediator, no dispatcher).
- [ ] The slice does not import another slice or the DI container.
- [ ] Request/Response never reach the handler; Command/Result never reach HTTP.
- [ ] Command handler changes exactly one aggregate and writes events to the Outbox in the same transaction.
- [ ] Query handler does not load aggregates and does not write.
- [ ] Consumer deduplicates through the Inbox and holds no business logic.
- [ ] Every acceptance criterion has a test.
- [ ] No personal data in logs; `tenantId` on every read and write.
- [ ] Explain to the developer why each file exists and where it is wired.
