package receivetelegramupdate

import (
	"context"
	"errors"
	"fmt"

	"github.com/parley/parley/services/go/channel-gateway/internal/domain/entities"
	"github.com/parley/parley/services/go/channel-gateway/internal/domain/enums"
	"github.com/parley/parley/services/go/channel-gateway/internal/domain/repositories"
	"github.com/parley/parley/services/go/channel-gateway/internal/domain/valueobjects"
)

type ReceiveTelegramUpdateHandler struct {
	messages repositories.InboundMessageRepository
	clock    Clock
}

func NewReceiveTelegramUpdateHandler(messages repositories.InboundMessageRepository, clock Clock) *ReceiveTelegramUpdateHandler {
	return &ReceiveTelegramUpdateHandler{
		messages: messages,
		clock:    clock,
	}
}

func (h *ReceiveTelegramUpdateHandler) Handle(ctx context.Context, cmd ReceiveTelegramUpdateCommand) (ReceiveTelegramUpdateResult, error) {
	if err := validate(cmd); err != nil {
		return ReceiveTelegramUpdateResult{}, err
	}

	sender, err := NewContactAddress(enums.ChannelTelegram, cmd.ExternalChatID)
	if err != nil {
		return ReceiveTelegramUpdateResult{}, fmt.Errorf("%w: external chat id is not valid", ErrInvalidCommand)
	}

	message, err := entities.ReceiveMessage(cmd.TenantID, sender, cmd.ProviderMessageID, cmd.Text, h.clock.Now())
	if err != nil {
		return ReceiveTelegramUpdateResult{}, fmt.Errorf("%w: failed to receive message", ErrInvalidCommand)
	}

	err = h.messages.Add(ctx, message)
	if errors.Is(err, repositories.ErrAlreadyExists) {
		return ReceiveTelegramUpdateResult{Duplicate: true}, nil
	}
	if err != nil {
		return ReceiveTelegramUpdateResult{}, fmt.Errorf("add inbound message: %w", err)
	}

	return ReceiveTelegramUpdateResult{InboundMessageID: message.ID()}, nil
}
