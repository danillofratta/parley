---
paths:
  - "services/python/**"
---
# Python services

- Python 3.13, managed with `uv`. `src/` layout: `src/<package>/domain`,
  `src/<package>/features/<use_case>`, `src/<package>/infrastructure`,
  `src/<package>/bootstrap`.
- Slice files follow the table in `architecture.md` (`message.py`,
  `command.py`, `result.py`, `validator.py`, `handler.py`, `consumer.py`, ...).
- Dependency injection: `dependency-injector`. The `Container` lives in
  `bootstrap/container.py` and is the only place that knows concrete classes.
  Slices use plain constructor injection and never import the container or
  use `@inject`/`Provide` markers (the slice stays framework-free and testable).
- Lint and format with `ruff`; type-check with `mypy --strict` (or pyright).
- Pydantic v2 for request, response, message, command and result models;
  frozen dataclasses for domain value objects.
- Async I/O throughout (`asyncio`, async Kafka and PostgreSQL clients).
- The LangGraph graph of a use case lives inside its slice (`graph.py`).
  Nodes call domain logic and ports; they never talk to Kafka or the database.
- LLM access only through `infrastructure/llm.py`, which builds the chat model
  from `LLM_PROVIDER` and `LLM_MODEL` (default `groq`). Never import a
  provider package inside a slice.
- Tests: `pytest` with `pytest-asyncio`; architecture checked with import-linter.
