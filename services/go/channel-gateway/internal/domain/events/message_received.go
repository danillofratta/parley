package events

import (
	"time"

	"github.com/google/uuid"

	"github.com/danillofratta/parley/services/go/channel-gateway/internal/domain/seedwork"
	"github.com/danillofratta/parley/services/go/channel-gateway/internal/domain/valueobjects"
)

// MessageReceived: a contact sent a message through a channel and Parley accepted it.
type MessageReceived struct {
	seedwork.EventBase
	Sender            valueobjects.ContactAddress
	ProviderMessageID string
	Text              string
}

func NewMessageReceived(
	inboundMessageID, tenantID uuid.UUID,
	sender valueobjects.ContactAddress,
	providerMessageID, text string,
	receivedAt time.Time,
) MessageReceived {
	return MessageReceived{
		EventBase:         seedwork.NewEventBase(inboundMessageID, tenantID, receivedAt),
		Sender:            sender,
		ProviderMessageID: providerMessageID,
		Text:              text,
	}
}
