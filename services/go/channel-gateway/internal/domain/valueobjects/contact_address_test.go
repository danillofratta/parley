package valueobjects

import (
	"errors"
	"testing"

	"github.com/danillofratta/parley/services/go/channel-gateway/internal/domain/enums"
	"github.com/danillofratta/parley/services/go/channel-gateway/internal/domain/seedwork"
)

func TestNewContactAddress(t *testing.T) {
	tests := []struct {
		name     string
		channel  enums.Channel
		chat     string
		wantCode string // empty means valid
	}{
		{"valid telegram chat", enums.ChannelTelegram, "111", ""},
		{"unknown channel", enums.ChannelUnknown, "111", "contact_address.unknown_channel"},
		{"blank chat", enums.ChannelTelegram, "  ", "contact_address.chat_required"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			address, err := NewContactAddress(tt.channel, tt.chat)

			if tt.wantCode == "" {
				if err != nil || address.IsZero() {
					t.Fatalf("address=%+v err=%v", address, err)
				}
				return
			}
			var violation *seedwork.BusinessRuleViolation
			if !errors.As(err, &violation) || violation.Rule.Code() != tt.wantCode {
				t.Fatalf("err = %v, want violation %s", err, tt.wantCode)
			}
		})
	}
}

func TestContactAddressesAreComparedByValue(t *testing.T) {
	a, _ := NewContactAddress(enums.ChannelTelegram, "111")
	b, _ := NewContactAddress(enums.ChannelTelegram, "111")

	if a != b {
		t.Fatal("equal addresses compared as different")
	}
}
