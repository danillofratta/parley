---
name: new-aggregate
description: Use when modeling a new aggregate, entity, value object, enum, business rule or domain event in any Parley service. Applies the domain model standard (Id, audit fields, Version, TenantId, CheckRule) identically in Go, .NET and Python.
---

# New aggregate

Read `.claude/rules/domain-model.md` first, including the domain folder
layout and its dependency order (enums <- rules <- value objects <- events <-
entities <- repositories). Put each concept in its folder. The seedwork for the stack must
already exist in the service (`domain/seedwork`); if not, copy it from
`.claude/skills/new-service/seedwork/<stack>/`.

1. **Name it** with the ubiquitous language. Missing term? Add it to
   `docs/domain/ubiquitous-language.md` first. Entities are named for what they
   are in the business, never after an identifier or a technical key.
2. **Classify each concept**: aggregate root (consistency boundary, has
   `TenantId`), entity inside the aggregate (identity, audit), value object
   (no identity, immutable), enum (closed set).
3. **List the invariants** and write one business-rule type per invariant
   (`ConversationMustBeOpen`) with `IsBroken`, `Code` (`conversation.must_be_open`)
   and `Message`.
4. **Factory method** named after the business action (`OpenConversation`):
   checks rules with `CheckRules`, creates the base (`Id` v7, `TenantId`), sets
   business attributes, raises the creation event.
5. **Behaviour methods** for every state change; each checks its rules and
   raises an event. No public setters. Audit fields are never set by the domain.
6. **Domain events** in past tense, carrying `EventId`, `OccurredAt`,
   `AggregateId`, `TenantId` (from the seedwork base) plus business fields.
7. **Repository port** in `domain/`: `GetById(tenantId, id)`, `Add`, `Update`;
   standard outcomes (not found, already exists, concurrency conflict).
8. **Table**: standard columns (`id`, `tenant_id`, `created_at`, `created_by`,
   `updated_at`, `updated_by`, `version`) plus business columns; business keys
   are `UNIQUE` including `tenant_id`.
9. **Tests first** for the rules and events: pure unit tests, no mocks, no I/O.
10. Explain to the developer why each concept is an entity, value object or enum.
