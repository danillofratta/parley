package entities

import (
	"time"

	"github.com/google/uuid"

	"github.com/danillofratta/parley/services/go/channel-gateway/internal/domain/events"
	"github.com/danillofratta/parley/services/go/channel-gateway/internal/domain/rules"
	"github.com/danillofratta/parley/services/go/channel-gateway/internal/domain/seedwork"
	"github.com/danillofratta/parley/services/go/channel-gateway/internal/domain/valueobjects"
)

// InboundMessage is a message a contact sent to a tenant through a channel.
// Aggregate root. Parley accepts each provider message only once per tenant and channel.
type InboundMessage struct {
	seedwork.AggregateRoot
	sender            valueobjects.ContactAddress
	providerMessageID string
	receivedAt        time.Time
}

// ReceiveMessage accepts a message from a contact and raises MessageReceived.
func ReceiveMessage(
	tenantID uuid.UUID,
	sender valueobjects.ContactAddress,
	providerMessageID string,
	text string,
	receivedAt time.Time,
) (*InboundMessage, error) {
	if err := seedwork.CheckRules(
		rules.SenderMustBeInformed{Sender: sender},
		rules.ProviderMessageMustBeIdentified{ProviderMessageID: providerMessageID},
		rules.MessageMustHaveText{Text: text},
	); err != nil {
		return nil, err
	}

	root, err := seedwork.NewAggregateRoot(tenantID)
	if err != nil {
		return nil, err
	}

	message := &InboundMessage{
		AggregateRoot:     root,
		sender:            sender,
		providerMessageID: providerMessageID,
		receivedAt:        receivedAt,
	}
	message.Raise(events.NewMessageReceived(
		message.ID(), tenantID, sender, providerMessageID, text, receivedAt))
	return message, nil
}

func (m *InboundMessage) Sender() valueobjects.ContactAddress { return m.sender }
func (m *InboundMessage) ProviderMessageID() string           { return m.providerMessageID }
func (m *InboundMessage) ReceivedAt() time.Time               { return m.receivedAt }
