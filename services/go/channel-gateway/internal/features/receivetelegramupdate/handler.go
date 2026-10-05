package receivetelegramupdate

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/danillofratta/parley/services/go/channel-gateway/internal/domain/entities"
	"github.com/danillofratta/parley/services/go/channel-gateway/internal/domain/enums"
	"github.com/danillofratta/parley/services/go/channel-gateway/internal/domain/repositories"
	"github.com/danillofratta/parley/services/go/channel-gateway/internal/domain/seedwork"
	"github.com/danillofratta/parley/services/go/channel-gateway/internal/domain/valueobjects"
)

type Handler struct {
	messages repositories.InboundMessageRepository
	clock    Clock
}

func NewHandler(messages repositories.InboundMessageRepository, clock Clock) *Handler {
	return &Handler{messages: messages, clock: clock}
}

// Handle accepts a Telegram message once. Business rule violations are returned
// unchanged (*seedwork.BusinessRuleViolation) so the caller keeps the rule code.
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
		return Result{}, err
	}

	message, err := entities.ReceiveMessage(tenantID, sender, cmd.ProviderMessageID, cmd.Text, h.clock.Now().UTC())
	if err != nil {
		return Result{}, err
	}

	err = h.messages.Add(ctx, message)
	if errors.Is(err, seedwork.ErrAlreadyExists) {
		return Result{Duplicate: true}, nil
	}
	if err != nil {
		return Result{}, fmt.Errorf("add inbound message: %w", err)
	}

	return Result{InboundMessageID: message.ID().String()}, nil
}
