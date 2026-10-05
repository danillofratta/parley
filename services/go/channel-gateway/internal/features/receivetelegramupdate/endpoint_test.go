package receivetelegramupdate

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const textUpdate = `{"update_id":5001,"message":{"chat":{"id":111},"text":"oi"}}`

func newTestEndpoint(t *testing.T, messages *fakeMessages) *Endpoint {
	t.Helper()
	e, err := NewEndpoint(
		NewHandler(messages, fixedClock{t: testNow}),
		EndpointConfig{WebhookSecret: "secret", TenantID: testTenant},
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)
	if err != nil {
		t.Fatal(err)
	}
	return e
}

func post(e http.Handler, secret, body string) int {
	req := httptest.NewRequest(http.MethodPost, "/webhooks/telegram", strings.NewReader(body))
	req.Header.Set(secretHeader, secret)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec.Code
}

func TestNewEndpointFailsClosed(t *testing.T) {
	handler := NewHandler(newFakeMessages(), fixedClock{t: testNow})
	tests := []struct {
		name    string
		handler *Handler
		cfg     EndpointConfig
	}{
		{"without handler", nil, EndpointConfig{WebhookSecret: "secret", TenantID: testTenant}},
		{"without secret", handler, EndpointConfig{TenantID: testTenant}},
		{"without tenant", handler, EndpointConfig{WebhookSecret: "secret"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := NewEndpoint(tt.handler, tt.cfg, nil); err == nil {
				t.Fatal("expected an error")
			}
		})
	}
}

func TestEndpointRejectsInvalidSecret(t *testing.T) {
	messages := newFakeMessages()

	got := post(newTestEndpoint(t, messages), "wrong", textUpdate)

	if got != http.StatusForbidden || len(messages.added) != 0 {
		t.Fatalf("status=%d stored=%d", got, len(messages.added))
	}
}

func TestEndpointStoresTextUpdateOnce(t *testing.T) {
	messages := newFakeMessages()
	e := newTestEndpoint(t, messages)

	first := post(e, "secret", textUpdate)
	second := post(e, "secret", textUpdate)

	if first != http.StatusOK || second != http.StatusOK || len(messages.added) != 1 {
		t.Fatalf("first=%d second=%d stored=%d", first, second, len(messages.added))
	}
}

func TestEndpointAcknowledgesNonTextUpdateWithoutStoring(t *testing.T) {
	messages := newFakeMessages()

	got := post(newTestEndpoint(t, messages), "secret", `{"update_id":5002,"message":{"chat":{"id":111}}}`)

	if got != http.StatusOK || len(messages.added) != 0 {
		t.Fatalf("status=%d stored=%d", got, len(messages.added))
	}
}

func TestEndpointAcknowledgesBrokenRuleWithoutStoring(t *testing.T) {
	messages := newFakeMessages()

	got := post(newTestEndpoint(t, messages), "secret", `{"update_id":5003,"message":{"chat":{"id":111},"text":"   "}}`)

	if got != http.StatusOK || len(messages.added) != 0 {
		t.Fatalf("status=%d stored=%d", got, len(messages.added))
	}
}

func TestEndpointReturns500OnStorageFailureSoTelegramRetries(t *testing.T) {
	messages := newFakeMessages()
	messages.err = errors.New("database down")

	got := post(newTestEndpoint(t, messages), "secret", textUpdate)

	if got != http.StatusInternalServerError {
		t.Fatalf("status=%d", got)
	}
}
