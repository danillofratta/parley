namespace Parley.Service.Domain.SeedWork;

/// <summary>Entity that owns a consistency boundary, belongs to a tenant and raises domain events.</summary>
public abstract class AggregateRoot : Entity
{
    private readonly List<IDomainEvent> _domainEvents = [];

    public Guid TenantId { get; private init; }

    public IReadOnlyCollection<IDomainEvent> DomainEvents => _domainEvents.AsReadOnly();

    protected AggregateRoot(Guid tenantId)
    {
        CheckRule(new TenantMustBeInformed(tenantId));
        TenantId = tenantId;
    }

    /// <summary>For ORM materialization only.</summary>
    protected AggregateRoot() { }

    protected void Raise(IDomainEvent domainEvent) => _domainEvents.Add(domainEvent);

    /// <summary>Called by persistence after the outbox rows are committed.</summary>
    public void ClearDomainEvents() => _domainEvents.Clear();
}

public sealed class TenantMustBeInformed(Guid tenantId) : IBusinessRule
{
    public bool IsBroken() => tenantId == Guid.Empty;
    public string Code => "aggregate.tenant_required";
    public string Message => "Every aggregate must belong to a tenant.";
}
