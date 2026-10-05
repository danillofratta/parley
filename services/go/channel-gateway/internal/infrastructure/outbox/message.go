package outbox

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/parley/parley/services/go/channel-gateway/internal/domain/events"
	"github.com/parley/parley/services/go/channel-gateway/internal/domain/valueobjects"
)

const topicMessageInbound = "messages.inbound"

type Message struct {
	ID         string
	TenantID   string
	Type       string
	Topic      string
	Key        string
	Payload    []byte
	OccurredAt time.Time
}

type envelope struct {
	MessageID       string    `json:"messageId"`
	Type            string    `json:"type"`
	TenantID        string    `json:"tenantId"`
	ConversationKey string    `json:"conversationKey"`
	OccurredAt      time.Time `json:"occurredAt"`
	Payload         any       `json:"payload"`
	// traceparent is added with OC-601.
}

type messageReceivedV1 struct {
	InboundMessageID  string    `json:"inboundMessageId"`
	Channel           string    `json:"channel"`
	ProviderMessageID string    `json:"providerMessageId"`
	Text              string    `json:"text"`
	ReceivedAt        time.Time `json:"receivedAt"`
}

// FromDomainEvent maps a domain event to the integration event that leaves the service.
func FromDomainEvent(event events.DomainEvent) (Message, error) {
	switch e := event.(type) {
	case events.MessageReceived:
		occurredAt := time.Unix(e.OccurredAtUnix(), 0).UTC()
		env := envelope{
			MessageID:       uuid.NewString(),
			Type:            e.EventType(),
			TenantID:        e.TenantID(),
			ConversationKey: conversationKey(e.Sender()),
			OccurredAt:      occurredAt,
			Payload: messageReceivedV1{
				InboundMessageID:  e.InboundMessageID(),
				Channel:           e.Sender().Channel().String(),
				ProviderMessageID: e.ProviderMessageID(),
				Text:              e.Text(),
				ReceivedAt:        time.Unix(e.ReceivedAtUnix(), 0).UTC(),
			},
		}

		payload, err := json.Marshal(env)
		if err != nil {
			return Message{}, fmt.Errorf("marshal %s: %w", env.Type, err)
		}

		return Message{
			ID:         env.MessageID,
			TenantID:   env.TenantID,
			Type:       env.Type,
			Topic:      topicMessageInbound,
			Key:        env.ConversationKey,
			Payload:    payload,
			OccurredAt: env.OccurredAt,
		}, nil
	default:
		return Message{}, fmt.Errorf("no integration event mapped for %T", event)
	}
}

// conversationKey is the Kafka message key "<channel>:<chat>"; it keeps a conversation's events in order.
func conversationKey(address valueobjects.ContactAddress) string {
	return address.Channel().String() + ":" + address.Chat()
}