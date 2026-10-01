# .NET slice template (Conversations: support agent sends a manual reply)

```
Features/Conversations/SendManualReply/
  SendManualReplyRequest.cs
  SendManualReplyResponse.cs
  SendManualReplyCommand.cs
  SendManualReplyResult.cs
  SendManualReplyValidator.cs
  SendManualReplyHandler.cs
  SendManualReplyEndpoint.cs
  SendManualReplyRegistration.cs
```

```csharp
// SendManualReplyRequest.cs  (HTTP body)
public sealed record SendManualReplyRequest(string Text);

// SendManualReplyResponse.cs (HTTP body)
public sealed record SendManualReplyResponse(Guid MessageId);

// SendManualReplyCommand.cs
public sealed record SendManualReplyCommand(Guid TenantId, Guid ConversationId, string Text);

// SendManualReplyResult.cs
public sealed record SendManualReplyResult(Guid MessageId);
```

```csharp
// SendManualReplyValidator.cs
public sealed class SendManualReplyValidator : AbstractValidator<SendManualReplyCommand>
{
    public SendManualReplyValidator()
    {
        RuleFor(c => c.TenantId).NotEmpty();
        RuleFor(c => c.ConversationId).NotEmpty();
        RuleFor(c => c.Text).NotEmpty().MaximumLength(4096);
    }
}
```

```csharp
// SendManualReplyHandler.cs
public sealed class SendManualReplyHandler(
    IValidator<SendManualReplyCommand> validator,
    IConversationRepository conversations,
    IOutbox outbox,
    IUnitOfWork unitOfWork) : ICommandHandler<SendManualReplyCommand, SendManualReplyResult>
{
    public async Task<SendManualReplyResult> Handle(SendManualReplyCommand command, CancellationToken ct)
    {
        await validator.ValidateAndThrowAsync(command, ct);

        var conversation = await conversations.Get(
            new TenantId(command.TenantId), new ConversationId(command.ConversationId), ct)
            ?? throw new ConversationNotFound(command.ConversationId);

        var message = conversation.SendManualReply(command.Text); // invariant: not closed

        outbox.AddRange(conversation.DomainEvents.ToIntegrationEvents());
        await unitOfWork.Commit(ct); // aggregate + outbox in one transaction

        return new SendManualReplyResult(message.Id.Value);
    }
}
```

```csharp
// SendManualReplyEndpoint.cs
public static class SendManualReplyEndpoint
{
    public static async Task<IResult> Handle(
        Guid conversationId,
        SendManualReplyRequest request,
        ITenantContext tenant,
        ICommandHandler<SendManualReplyCommand, SendManualReplyResult> handler, // injected, called directly
        CancellationToken ct)
    {
        var command = new SendManualReplyCommand(tenant.Id, conversationId, request.Text);
        var result = await handler.Handle(command, ct);
        return Results.Ok(new SendManualReplyResponse(result.MessageId));
    }
}
```

```csharp
// SendManualReplyRegistration.cs
public static class SendManualReplyRegistration
{
    public static IServiceCollection AddSendManualReply(this IServiceCollection services) =>
        services
            .AddScoped<IValidator<SendManualReplyCommand>, SendManualReplyValidator>()
            .AddScoped<ICommandHandler<SendManualReplyCommand, SendManualReplyResult>, SendManualReplyHandler>();

    public static IEndpointRouteBuilder MapSendManualReply(this IEndpointRouteBuilder app)
    {
        app.MapPost("/conversations/{conversationId:guid}/replies", SendManualReplyEndpoint.Handle);
        return app;
    }
}
```

```csharp
// Features/FeatureRegistration.cs (composition root, explicit list)
public static class FeatureRegistration
{
    public static IServiceCollection AddFeatures(this IServiceCollection services) =>
        services.AddSendManualReply();
    // .AddRegisterInboundMessage() ...

    public static IEndpointRouteBuilder MapFeatures(this IEndpointRouteBuilder app) =>
        app.MapSendManualReply();
}
```

A query slice uses `<UseCase>Query` and `IQueryHandler<,>`; its handler reads
with Dapper or `AsNoTracking()` straight into the Result and never touches
repositories or aggregates.
