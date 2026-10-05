namespace Parley.Service.Domain.SeedWork;

/// <summary>A business fact that already happened, named in past tense.</summary>
public interface IDomainEvent
{
    Guid EventId { get; }
    DateTimeOffset OccurredAt { get; }
    Guid AggregateId { get; }
    Guid TenantId { get; }
}

public abstract record DomainEvent(Guid AggregateId, Guid TenantId, DateTimeOffset OccurredAt) : IDomainEvent
{
    public Guid EventId { get; init; } = Guid.CreateVersion7();
}
