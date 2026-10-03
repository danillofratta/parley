package outbox

const topicMessageInbound = "message.inbound"

type Message struct {
	Id string
	TenantId string
	Type string
	Topic string
	Key string
	Payload []byte
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
func FromDomainEvent(event seedwork.DomainEvent) (Message, error) {
	switch e := event.(type) {
	case events.MessageReceived:
		env := envelope{
			MessageID:       e.EventID().String(),
			Type:            "MessageReceived.v1",
			TenantID:        e.TenantID().String(),
			ConversationKey: conversationKey(e.Sender),
			OccurredAt:      e.OccurredAt(),
			Payload: messageReceivedV1{
				InboundMessageID:  e.AggregateID().String(),
				Channel:           e.Sender.Channel().String(),
				ProviderMessageID: e.ProviderMessageID,
				Text:              e.Text,
				ReceivedAt:        e.OccurredAt(),
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
			Topic:      topicMessagesInbound,
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