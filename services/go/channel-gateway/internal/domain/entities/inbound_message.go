package entities

import (
	"time"

	"github.com/parley/parley/services/go/channel-gateway/internal/domain/events"
	"github.com/parley/parley/services/go/channel-gateway/internal/domain/rules"
	"github.com/parley/parley/services/go/channel-gateway/internal/domain/valueobjects"
)

type InboundMessage struct {
	AggregateRoot
	sender            valueobjects.ContactAddress
	providerMessageID string
	receivedAtUnix    int64
}

func ReceiveMessage(
	tenantID string,
	sender valueobjects.ContactAddress,
	providerMessageID string,
	text string,
	receivedAt time.Time,
) (*InboundMessage, error) {
	if err := rules.CheckRules(
		rules.SenderMustBeInformed{Sender: sender},
		rules.ProviderMessageMustBeIdentified{ProviderMessageID: providerMessageID},
		rules.MessageMustHaveText{Text: text},
	); err != nil {
		return nil, err
	}

	root, err := NewAggregateRoot(tenantID)
	if err != nil {
		return nil, err
	}

	message := &InboundMessage{
		AggregateRoot:     *root,
		sender:            sender,
		providerMessageID: providerMessageID,
		receivedAtUnix:    receivedAt.Unix(),
	}

	message.Raise(events.NewMessageReceived(
		message.ID(),
		tenantID,
		sender,
		providerMessageID,
		text,
		receivedAt,
	))

	return message, nil
}

func (m *InboundMessage) Sender() valueobjects.ContactAddress {
	return m.sender
}

func (m *InboundMessage) ProviderMessageID() string {
	return m.providerMessageID
}

func (m *InboundMessage) ReceivedAtUnix() int64 {
	return m.receivedAtUnix
}