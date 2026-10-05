namespace Parley.Service.Domain.SeedWork;

/// <summary>Persistence port for one aggregate root. No query methods.</summary>
public interface IRepository<T> where T : AggregateRoot
{
    /// <returns>The aggregate, or null when not found.</returns>
    Task<T?> GetById(Guid tenantId, Guid id, CancellationToken ct);

    /// <exception cref="AggregateAlreadyExistsException">A unique business key already exists.</exception>
    Task Add(T aggregate, CancellationToken ct);

    /// <exception cref="ConcurrencyConflictException">The stored Version differs.</exception>
    Task Update(T aggregate, CancellationToken ct);
}

public sealed class AggregateAlreadyExistsException(string aggregate, string key)
    : Exception($"{aggregate} '{key}' already exists.");

public sealed class ConcurrencyConflictException(string aggregate, Guid id)
    : Exception($"{aggregate} {id} was changed by another operation.");
