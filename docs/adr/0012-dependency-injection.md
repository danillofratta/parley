# ADR-0012: Dependency injection per stack

- Status: Accepted
- Date: 2026-09-30

## Context

Every slice depends on repositories, the Outbox, validators and clients.
Handlers must be testable with fakes, and concrete implementations must be
chosen in one place per service. The three stacks have different idioms.

## Options

### .NET
- Built-in `Microsoft.Extensions.DependencyInjection`: standard, fast, enough
  for constructor injection and scoped lifetimes.
- Third-party containers (Autofac, etc.): richer features we do not need.

### Python
- Manual composition root: no dependency, but a lot of repetitive wiring code
  and manual lifetime management.
- `dependency-injector`: widely used declarative container with singletons,
  factories and resources (async pools).
- Framework-specific injection (FastAPI `Depends`): only covers HTTP; our main
  entry points are Kafka consumers.

### Go
- Manual constructor injection in a composition root: the idiomatic default;
  compile-time checked, no reflection, readable.
- Uber Fx: runtime container with lifecycle hooks; errors appear at startup
  instead of compile time; more "magic" than Go code usually accepts.
- Google Wire: generates the wiring code at compile time; worth it for large
  graphs; check its maintenance status before adopting.

## Decision

- .NET: built-in DI; each slice registers itself in its Registration file;
  `FeatureRegistration.cs` lists slices explicitly (no assembly scanning).
- Python: `dependency-injector`, confined to `bootstrap/container.py`; slices
  use plain constructor injection and never import the container.
- Go: manual constructor injection in `internal/bootstrap`; no container.

## Consequences

- Positive: handlers are testable with fakes in every stack; one place per
  service shows the full dependency graph.
- Negative: explicit registration means one more line per new slice; Go
  wiring grows linearly with the number of components.
- Revisit Go if a service's wiring becomes hard to read (Wire would be the
  first candidate).
