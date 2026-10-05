package receivetelegramupdate

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/danillofratta/parley/services/go/channel-gateway/internal/domain/entities"
	"github.com/danillofratta/parley/services/go/channel-gateway/internal/domain/seedwork"
)

const testTenant = "00000000-0000-0000-0000-000000000001"

var testNow = time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC)

// fakeMessages is an in-memory repository that deduplicates like the real one.
// It proves the handler's behaviour, not the database guarantee (integration test, OC-602).
type fakeMessages struct {
	added map[string]*entities.InboundMessage
	err   error // forced failure, e.g. database down
}

func newFakeMessages() *fakeMessages {
	return &fakeMessages{added: map[string]*entities.InboundMessage{}}
}

func (f *fakeMessages) Add(_ context.Context, m *entities.InboundMessage) error {
	if f.err != nil {
		return f.err
	}
	key := m.TenantID().String() + "|" + m.Sender().Channel().String() + "|" + m.ProviderMessageID()
	if _, exists := f.added[key]; exists {
		return seedwork.ErrAlreadyExists
	}
	f.added[key] = m
	return nil
}

type fixedClock struct{ t time.Time }

func (c fixedClock) Now() time.Time { return c.t }

func validCommand() Command {
	return Command{
		TenantID:          testTenant,
		ProviderMessageID: "5001",
		ExternalChatID:    "111",
		Text:              "quero saber do meu pedido",
	}
}

func TestHandleAcceptsNewMessage(t *testing.T) {
	messages := newFakeMessages()
	h := NewHandler(messages, fixedClock{t: testNow})

	result, err := h.Handle(context.Background(), validCommand())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Duplicate || result.InboundMessageID == "" {
		t.Fatalf("unexpected result: %+v", result)
	}
	if len(messages.added) != 1 {
		t.Fatalf("added %d messages, want 1", len(messages.added))
	}
}

func TestHandleTreatsRedeliveryAsDuplicate(t *testing.T) {
	messages := newFakeMessages()
	h := NewHandler(messages, fixedClock{t: testNow})

	_, _ = h.Handle(context.Background(), validCommand())
	result, err := h.Handle(context.Background(), validCommand())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Duplicate {
		t.Fatal("redelivery not reported as duplicate")
	}
	if len(messages.added) != 1 {
		t.Fatalf("added %d messages, want 1", len(messages.added))
	}
}

func TestHandleRejectsInvalidInput(t *testing.T) {
	tests := []struct {
		name     string
		mutate   func(*Command)
		wantCode string // empty means ErrInvalidCommand
	}{
		{"missing tenant", func(c *Command) { c.TenantID = "" }, ""},
		{"malformed tenant", func(c *Command) { c.TenantID = "tenant-1" }, ""},
		{"missing provider message id", func(c *Command) { c.ProviderMessageID = "" }, ""},
		{"missing chat id", func(c *Command) { c.ExternalChatID = "" }, ""},
		{"blank text", func(c *Command) { c.Text = "   " }, "inbound_message.text_required"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			messages := newFakeMessages()
			h := NewHandler(messages, fixedClock{t: testNow})
			cmd := validCommand()
			tt.mutate(&cmd)

			_, err := h.Handle(context.Background(), cmd)

			var violation *seedwork.BusinessRuleViolation
			switch {
			case tt.wantCode == "" && !errors.Is(err, ErrInvalidCommand):
				t.Fatalf("err = %v, want ErrInvalidCommand", err)
			case tt.wantCode != "" && (!errors.As(err, &violation) || violation.Rule.Code() != tt.wantCode):
				t.Fatalf("err = %v, want violation %s", err, tt.wantCode)
			}
			if !cannotBeFixedByRetry(err) {
				t.Fatalf("err = %v should not be retried", err)
			}
			if len(messages.added) != 0 {
				t.Fatal("invalid input was stored")
			}
		})
	}
}

func TestHandleStorageFailureIsRetryable(t *testing.T) {
	messages := newFakeMessages()
	messages.err = errors.New("database down")
	h := NewHandler(messages, fixedClock{t: testNow})

	_, err := h.Handle(context.Background(), validCommand())

	if err == nil || cannotBeFixedByRetry(err) {
		t.Fatalf("err = %v, want a retryable error", err)
	}
}
