---
name: new-aggregate
description: Use when modeling a new aggregate, entity or value object in a Parley domain layer. Guides invariant discovery, boundaries and domain events.
---

# New aggregate

1. **Name it** with the ubiquitous language. If the term is missing, add it to
   `docs/domain/ubiquitous-language.md` first.
2. **List the invariants** it must protect ("a closed conversation cannot
   receive a reply approval"). No invariant, no aggregate: use a value object
   or a plain entity instead.
3. **Draw the boundary**: include only what must be consistent in the same
   transaction. Reference other aggregates by id.
4. **Design behaviour, not data**: public methods named after domain actions
   (`ApproveReply`, `RequestHandoff`), no public setters, a factory for creation.
5. **Value objects** for every concept with rules (`ConversationKey`,
   `Confidence`, `TenantId`): immutable, validated on creation, equality by value.
6. **Domain events** in past tense for every meaningful state change, recorded
   inside the aggregate and cleared after the handler saves them.
7. **Repository port** in `domain/`, one per aggregate root: `Get`, `Add`, and
   nothing query-like (queries belong to query slices).
8. **Tests first** for the invariants: pure unit tests, no mocks, no I/O.
9. Explain to the developer why each rule is inside the aggregate.
