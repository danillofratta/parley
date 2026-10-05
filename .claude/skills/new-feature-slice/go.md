# Go slice template

Canonical example: `services/go/channel-gateway/internal/features/receivetelegramupdate`.

```
internal/features/<usecase>/
  request.go        // transport model (HTTP body as received)
  command.go        // application input
  result.go         // application output
  validator.go      // shape of the command (required fields)
  ports.go          // dependencies only this slice needs (e.g. Clock)
  handler.go        // orchestration: value objects -> aggregate factory -> repository
  endpoint.go       // HTTP: auth, decode, map Request -> Command, map outcome -> status
  handler_test.go
  endpoint_test.go
```

```go
// handler.go
func (h *Handler) Handle(ctx context.Context, cmd Command) (Result, error) {
	if err := validate(cmd); err != nil {
		return Result{}, err
	}
	tenantID, err := uuid.Parse(cmd.TenantID)
	if err != nil {
		return Result{}, fmt.Errorf("%w: tenant id is not a uuid", ErrInvalidCommand)
	}
	sender, err := valueobjects.NewContactAddress(enums.ChannelTelegram, cmd.ExternalChatID)
	if err != nil {
		return Result{}, err // *seedwork.BusinessRuleViolation
	}
	message, err := entities.ReceiveMessage(tenantID, sender, cmd.ProviderMessageID, cmd.Text, h.clock.Now().UTC())
	if err != nil {
		return Result{}, err
	}
	err = h.messages.Add(ctx, message) // aggregate + outbox, one transaction
	if errors.Is(err, seedwork.ErrAlreadyExists) {
		return Result{Duplicate: true}, nil
	}
	if err != nil {
		return Result{}, fmt.Errorf("add inbound message: %w", err)
	}
	return Result{InboundMessageID: message.ID().String()}, nil
}
```

```go
// endpoint.go: mapping outcomes to the protocol
func cannotBeFixedByRetry(err error) bool {
	var violation *seedwork.BusinessRuleViolation
	return errors.Is(err, ErrInvalidCommand) || errors.As(err, &violation)
}
```

Composition root (`internal/bootstrap`): build the `persistence.Auditor`, the
repository, the handler and the endpoint, then register the route.
