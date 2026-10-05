package entities

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/danillofratta/parley/services/go/channel-gateway/internal/domain/enums"
	"github.com/danillofratta/parley/services/go/channel-gateway/internal/domain/events"
	"github.com/danillofratta/parley/services/go/channel-gateway/internal/domain/seedwork"
	"github.com/danillofratta/parley/services/go/channel-gateway/internal/domain/valueobjects"
)

func TestReceiveMessage(t *testing.T) {
	tenant := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	sender, err := valueobjects.NewContactAddress(enums.ChannelTelegram, "111")
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC)

	t.Run("creates the aggregate and raises MessageReceived", func(t *testing.T) {
		message, err := ReceiveMessage(tenant, sender, "5001", "oi", now)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if message.ID() == uuid.Nil || message.ID().Version() != 7 {
			t.Fatalf("id = %v, want a UUID v7", message.ID())
		}
		if message.TenantID() != tenant || message.Version() != 0 {
			t.Fatalf("tenant=%v version=%d", message.TenantID(), message.Version())
		}

		raised := message.DomainEvents()
		if len(raised) != 1 {
			t.Fatalf("raised %d events, want 1", len(raised))
		}
		event, ok := raised[0].(events.MessageReceived)
		if !ok {
			t.Fatalf("event is %T, want events.MessageReceived", raised[0])
		}
		if event.AggregateID() != message.ID() || event.TenantID() != tenant ||
			!event.OccurredAt().Equal(now) || event.Text != "oi" {
			t.Fatalf("unexpected event: %+v", event)
		}
	})

	tests := []struct {
		name       string
		tenant     uuid.UUID
		sender     valueobjects.ContactAddress
		providerID string
		text       string
		wantCode   string
	}{
		{"without tenant", uuid.Nil, sender, "5001", "oi", "aggregate.tenant_required"},
		{"without sender", tenant, valueobjects.ContactAddress{}, "5001", "oi", "inbound_message.sender_required"},
		{"without provider id", tenant, sender, "", "oi", "inbound_message.provider_message_id_required"},
		{"without text", tenant, sender, "5001", "  ", "inbound_message.text_required"},
	}
	for _, tt := range tests {
		t.Run("rejects a message "+tt.name, func(t *testing.T) {
			_, err := ReceiveMessage(tt.tenant, tt.sender, tt.providerID, tt.text, now)

			var violation *seedwork.BusinessRuleViolation
			if !errors.As(err, &violation) || violation.Rule.Code() != tt.wantCode {
				t.Fatalf("err = %v, want violation %s", err, tt.wantCode)
			}
		})
	}
}
