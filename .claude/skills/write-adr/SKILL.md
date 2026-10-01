---
name: write-adr
description: Use when an architecture decision is made or changed. Writes an ADR in docs/adr using the project template.
---

# Write an ADR

1. Next number: highest file in `docs/adr/` + 1, four digits, kebab-case title.
2. Copy `docs/adr/0000-template.md`.
3. Context: the forces and constraints, not the solution.
4. Options: at least two real alternatives with pros and cons.
5. Decision: one sentence, then the reasoning.
6. Consequences: good, bad, and what we will watch. Flag production gaps.
7. Status: Proposed → Accepted. A replaced decision is marked
   `Superseded by ADR-xxxx`, never deleted.
8. The developer writes the reasoning; Claude reviews it for missing options
   and hidden trade-offs.
