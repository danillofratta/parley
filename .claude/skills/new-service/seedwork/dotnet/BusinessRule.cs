namespace Parley.Service.Domain.SeedWork;

/// <summary>An invariant named in the ubiquitous language (e.g. ConversationMustBeOpen).</summary>
public interface IBusinessRule
{
    bool IsBroken();

    /// <summary>Stable code for clients, e.g. "conversation.must_be_open".</summary>
    string Code { get; }

    string Message { get; }
}

/// <summary>The single exception type the domain throws for a broken rule.</summary>
public sealed class BusinessRuleViolationException(IBusinessRule rule)
    : Exception($"{rule.Code}: {rule.Message}")
{
    public IBusinessRule Rule { get; } = rule;
}
