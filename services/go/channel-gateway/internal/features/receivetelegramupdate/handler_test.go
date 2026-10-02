package receivetelegramupdate

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/parley/parley/services/go/channel-gateway/internal/domain/entities"
	"github.com/parley/parley/services/go/channel-gateway/internal/domain/repositories"
	"github.com/parley/parley/services/go/channel-gateway/internal/domain/rules"
)

type fakeClock struct{}

func (fakeClock) Now() time.Time { return time.Unix(1_700_000_000, 0) }

type fakeRepo struct {
	added []*entities.InboundMessage
	err   error
}

func (f *fakeRepo) Add(_ context.Context, m *entities.InboundMessage) error {
	if f.err != nil {
		return f.err
	}
	f.added = append(f.added, m)
	return nil
}

func (f *fakeRepo) FindByID(context.Context, string) (*entities.InboundMessage, error) {
	return nil, errors.New("not used")
}

func validCommand() ReceiveTelegramUpdateCommand {
	return ReceiveTelegramUpdateCommand{TenantID: "tenant-1", ProviderMessageID: "10", ExternalChatID: "42", Text: "hello"}
}

func TestValidMessageIsStored(t *testing.T) {
	repo := &fakeRepo{}
	res, err := NewReceiveTelegramUpdateHandler(repo, fakeClock{}).Handle(context.Background(), validCommand())
	if err != nil || res.InboundMessageID == "" || res.Duplicate || len(repo.added) != 1 {
		t.Fatalf("res=%+v err=%v stored=%d", res, err, len(repo.added))
	}
}

func TestDuplicateUpdateIsReportedNotFailed(t *testing.T) {
	repo := &fakeRepo{err: repositories.ErrAlreadyExists}
	res, err := NewReceiveTelegramUpdateHandler(repo, fakeClock{}).Handle(context.Background(), validCommand())
	if err != nil || !res.Duplicate {
		t.Fatalf("res=%+v err=%v", res, err)
	}
}

func TestMissingFieldsAreRejectedAsInvalidCommand(t *testing.T) {
	cmd := validCommand()
	cmd.ProviderMessageID = ""
	_, err := NewReceiveTelegramUpdateHandler(&fakeRepo{}, fakeClock{}).Handle(context.Background(), cmd)
	if !errors.Is(err, ErrInvalidCommand) || !cannotBeFixedByRetry(err) {
		t.Fatalf("err=%v", err)
	}
}

func TestDomainRuleViolationCannotBeFixedByRetry(t *testing.T) {
	// The validator and the domain rules check the same things today, so no
	// input reaches a rule through the handler; classify the error directly.
	err := fmt.Errorf("receive: %w", rules.NewBrokenRuleError("inbound_message.text_required", "an inbound message must have text"))
	if !cannotBeFixedByRetry(err) {
		t.Fatalf("err=%v", err)
	}
}

func TestStorageFailureIsRetryable(t *testing.T) {
	repo := &fakeRepo{err: errors.New("db down")}
	_, err := NewReceiveTelegramUpdateHandler(repo, fakeClock{}).Handle(context.Background(), validCommand())
	if err == nil || cannotBeFixedByRetry(err) {
		t.Fatalf("err=%v", err)
	}
}

func newTestEndpoint(t *testing.T, repo *fakeRepo) *Endpoint {
	t.Helper()
	ep, err := NewEndpoint(NewReceiveTelegramUpdateHandler(repo, fakeClock{}),
		EndpointConfig{WebhookSecret: "s3cret", TenantID: "tenant-1"},
		slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatal(err)
	}
	return ep
}

func post(ep *Endpoint, secret, body string) int {
	req := httptest.NewRequest(http.MethodPost, "/webhooks/telegram", strings.NewReader(body))
	if secret != "" {
		req.Header.Set(secretHeader, secret)
	}
	rec := httptest.NewRecorder()
	ep.ServeHTTP(rec, req)
	return rec.Code
}

func TestWrongSecretIsForbidden(t *testing.T) {
	if got := post(newTestEndpoint(t, &fakeRepo{}), "nope", `{}`); got != http.StatusForbidden {
		t.Fatalf("status=%d", got)
	}
}

func TestEndpointRequiresSecret(t *testing.T) {
	if _, err := NewEndpoint(nil, EndpointConfig{TenantID: "t"}, nil); err == nil {
		t.Fatal("expected error")
	}
}

func TestEndpointRequiresHandler(t *testing.T) {
	if _, err := NewEndpoint(nil, EndpointConfig{WebhookSecret: "s3cret", TenantID: "t"}, nil); err == nil {
		t.Fatal("expected error")
	}
}

func TestEndpointAcceptsNilLogger(t *testing.T) {
	ep, err := NewEndpoint(
		NewReceiveTelegramUpdateHandler(&fakeRepo{}, fakeClock{}),
		EndpointConfig{WebhookSecret: "s3cret", TenantID: "t"},
		nil,
	)
	if err != nil || ep == nil {
		t.Fatalf("ep=%v err=%v", ep, err)
	}
}

func TestTextUpdateIsStored(t *testing.T) {
	repo := &fakeRepo{}
	body := `{"update_id":10,"message":{"chat":{"id":42},"text":"hello"}}`
	if got := post(newTestEndpoint(t, repo), "s3cret", body); got != http.StatusOK || len(repo.added) != 1 {
		t.Fatalf("status=%d stored=%d", got, len(repo.added))
	}
}

func TestNonTextUpdateIsAcknowledgedAndIgnored(t *testing.T) {
	repo := &fakeRepo{}
	if got := post(newTestEndpoint(t, repo), "s3cret", `{"update_id":11}`); got != http.StatusOK || len(repo.added) != 0 {
		t.Fatalf("status=%d stored=%d", got, len(repo.added))
	}
}

func TestStorageFailureReturns500SoTelegramRetries(t *testing.T) {
	repo := &fakeRepo{err: errors.New("db down")}
	body := `{"update_id":10,"message":{"chat":{"id":42},"text":"hello"}}`
	if got := post(newTestEndpoint(t, repo), "s3cret", body); got != http.StatusInternalServerError {
		t.Fatalf("status=%d", got)
	}
}
