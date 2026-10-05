namespace Parley.Service.Domain.SeedWork;

/// <summary>Base of every entity: identity, audit trail and version. Equality by Id.</summary>
public abstract class Entity
{
    public Guid Id { get; protected init; } = Guid.CreateVersion7();
    public DateTimeOffset CreatedAt { get; private set; }
    public string CreatedBy { get; private set; } = string.Empty;
    public DateTimeOffset? UpdatedAt { get; private set; }
    public string? UpdatedBy { get; private set; }
    public int Version { get; private set; }

    /// <summary>Called by persistence (audit interceptor) on the first save.</summary>
    public void StampCreated(DateTimeOffset at, string by)
    {
        CreatedAt = at;
        CreatedBy = by;
        Version = 1;
    }

    /// <summary>Called by persistence (audit interceptor) on every update.</summary>
    public void StampUpdated(DateTimeOffset at, string by)
    {
        UpdatedAt = at;
        UpdatedBy = by;
        Version++;
    }

    protected static void CheckRule(IBusinessRule rule)
    {
        if (rule.IsBroken())
        {
            throw new BusinessRuleViolationException(rule);
        }
    }

    protected static void CheckRules(params IBusinessRule[] rules)
    {
        foreach (var rule in rules)
        {
            CheckRule(rule);
        }
    }

    public override bool Equals(object? obj) =>
        obj is Entity other && GetType() == other.GetType() && Id == other.Id;

    public override int GetHashCode() => Id.GetHashCode();
}
