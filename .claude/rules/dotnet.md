---
paths:
  - "services/dotnet/**"
---
# .NET services

- .NET 10 (LTS), SDK pinned in `global.json`. Nullable enabled,
  `TreatWarningsAsErrors`, file-scoped namespaces.
- Projects per service: `<Service>.Domain` (no framework references),
  `<Service>` (Minimal API host + `Features/` + `Infrastructure/`),
  and test projects. The Domain project boundary is enforced by the compiler.
- Slice files follow the table in `architecture.md`, one type per file:
  `<UseCase>Request`, `<UseCase>Response`, `<UseCase>Command`/`Query`,
  `<UseCase>Result`, `<UseCase>Validator`, `<UseCase>Handler`,
  `<UseCase>Endpoint`, `<UseCase>Registration`.
- No mediator library (no MediatR). Handlers implement
  `ICommandHandler<TCommand, TResult>` or `IQueryHandler<TQuery, TResult>`
  (defined in `Infrastructure/Cqrs`) and are injected into the endpoint.
- Dependency injection: built-in `Microsoft.Extensions.DependencyInjection`.
  Each slice exposes `Add<UseCase>(this IServiceCollection)` and
  `Map<UseCase>(this IEndpointRouteBuilder)` in its Registration file;
  `Features/FeatureRegistration.cs` lists every slice explicitly (no assembly
  scanning). Handlers, validators, repositories and the DbContext are Scoped.
- Validation with FluentValidation, injected into the handler.
- No AutoMapper: mapping is explicit in the endpoint (Request → Command,
  Result → Response).
- Commands, queries, results, requests, responses and value objects are
  `record`s; aggregates are `sealed` classes with private setters and
  factory methods.
- Write side: EF Core with explicit mapping in `Infrastructure/` (no data
  annotations in the Domain project). Read side: Dapper or EF `AsNoTracking`
  projections straight into the query Result.
- Tests: xUnit, Testcontainers, architecture tests with NetArchTest or ArchUnitNET.
