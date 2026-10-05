namespace Parley.Service.Domain.SeedWork;

/// <summary>
/// Base of value objects: immutable and compared by value (record equality).
/// Validate in the factory method and keep constructors private.
/// </summary>
public abstract record ValueObject;
