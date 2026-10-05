package kafka

import (
	"context"
	"fmt"

	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/danillofratta/parley/services/go/channel-gateway/internal/infrastructure/outbox"
)

// Publisher sends outbox messages to Kafka. franz-go uses an idempotent producer by
// default, so client retries neither duplicate nor reorder records of a partition.
type Publisher struct {
	client *kgo.Client
}

func NewPublisher(brokers []string, clientID string) (*Publisher, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.ClientID(clientID),
	)
	if err != nil {
		return nil, fmt.Errorf("create kafka client: %w", err)
	}
	return &Publisher{client: client}, nil
}

// Publish blocks until the broker acknowledges the record or ctx expires.
func (p *Publisher) Publish(ctx context.Context, msg outbox.PendingMessage) error {
	record := &kgo.Record{
		Topic: msg.Topic,
		Key:   []byte(msg.Key), // conversationKey: same key, same partition, same order
		Value: msg.Payload,
	}
	for k, v := range msg.Headers {
		record.Headers = append(record.Headers, kgo.RecordHeader{Key: k, Value: []byte(v)})
	}
	if err := p.client.ProduceSync(ctx, record).FirstErr(); err != nil {
		return fmt.Errorf("produce to %s: %w", msg.Topic, err)
	}
	return nil
}

func (p *Publisher) Close() {
	p.client.Close()
}
