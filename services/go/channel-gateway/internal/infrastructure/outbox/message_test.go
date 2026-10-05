package outbox

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/danillofratta/parley/services/go/channel-gateway/internal/domain/enums"
	"github.com/danillofratta/parley/services/go/channel-gateway/internal/domain/events"
	"github.com/danillofratta/parley/services/go/channel-gateway/internal/domain/valueobjects"
)

func TestFromDomainEventMapsMessageReceived(t *testing.T) {
	sender, err := valueobjects.NewContactAddress(enums.ChannelTelegram, "111")
	if err != nil {
		t.Fatal(err)
	}
	tenant := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	messageID := uuid.MustParse("01920000-0000-7000-8000-000000000001")
	receivedAt := time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC)
	event := events.NewMessageReceived(messageID, tenant, sender, "5001", "oi", receivedAt)

	msg, err := FromDomainEvent(event)
	if err != nil {
		t.Fatal(err)
	}

	if msg.ID != event.EventID().String() {
		t.Fatalf("outbox id = %s, want the domain event id %s", msg.ID, event.EventID())
	}
	if msg.Topic != "messages.inbound" || msg.Key != "telegram:111" || msg.Type != "MessageReceived.v1" {
		t.Fatalf("unexpected routing: topic=%s key=%s type=%s", msg.Topic, msg.Key, msg.Type)
	}

	var env map[string]any
	if err := json.Unmarshal(msg.Payload, &env); err != nil {
		t.Fatal(err)
	}
	if env["messageId"] != msg.ID || env["conversationKey"] != "telegram:111" {
		t.Fatalf("unexpected envelope: %v", env)
	}
}
