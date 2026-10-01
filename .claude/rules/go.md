---
paths:
  - "services/go/**"
---
# Go services

- Layout: `cmd/<binary>/main.go`, `internal/domain`, `internal/features/<usecase>`,
  `internal/infrastructure/<adapter>`, `internal/bootstrap`. Everything under `internal/`.
- Slice files follow the table in `architecture.md` (`request.go`, `command.go`,
  `result.go`, `validator.go`, `handler.go`, `endpoint.go`, ...).
- Dependency injection: manual constructor injection. Every component has a
  `NewX(deps...) *X` constructor; `internal/bootstrap` builds the whole graph
  and `main.go` only loads config, calls bootstrap and runs the server.
  No DI container, no global variables, no `init()` side effects.
- Repository ports live in `internal/domain`; slice-specific interfaces in the
  slice's `ports.go`. Keep interfaces small (one to three methods).
- Standard library first: `net/http` with method patterns, `log/slog` (JSON),
  `context`. Database: `pgx/v5`. Avoid web frameworks.
- `context.Context` is the first parameter of every I/O function.
- Wrap errors with `fmt.Errorf("...: %w", err)`; domain errors are sentinel
  values or typed errors in `domain`, checked with `errors.Is/As`.
- Graceful shutdown on SIGTERM/Ctrl+C; server timeouts always set.
- Tests: table-driven, `httptest` for endpoints, Testcontainers for integration.
- Format with `gofmt`; lint with `golangci-lint`.
