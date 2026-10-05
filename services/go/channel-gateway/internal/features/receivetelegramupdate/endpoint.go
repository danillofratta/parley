package receivetelegramupdate

import (
	"crypto/subtle"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/danillofratta/parley/services/go/channel-gateway/internal/domain/seedwork"
)

const (
	secretHeader = "X-Telegram-Bot-Api-Secret-Token"
	maxBodyBytes = 1 << 20 // 1 MiB
)

type EndpointConfig struct {
	WebhookSecret string
	TenantID      string
}

type Endpoint struct {
	handler *Handler // called directly: no mediator
	cfg     EndpointConfig
	log     *slog.Logger
}

// NewEndpoint fails closed: without a webhook secret every request would pass
// the constant-time comparison against an empty header.
func NewEndpoint(handler *Handler, cfg EndpointConfig, log *slog.Logger) (*Endpoint, error) {
	if handler == nil {
		return nil, errors.New("handler is required")
	}
	if cfg.WebhookSecret == "" {
		return nil, errors.New("telegram webhook secret is required")
	}
	if cfg.TenantID == "" {
		return nil, errors.New("tenant id is required")
	}
	if log == nil {
		log = slog.Default()
	}
	return &Endpoint{handler: handler, cfg: cfg, log: log}, nil
}

func (e *Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	got := r.Header.Get(secretHeader)
	if subtle.ConstantTimeCompare([]byte(got), []byte(e.cfg.WebhookSecret)) != 1 {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxBodyBytes))
	if err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			http.Error(w, "body too large", http.StatusRequestEntityTooLarge)
			return
		}
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}

	var req Request
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "failed to parse body", http.StatusBadRequest)
		return
	}

	// Release 1 accepts text messages only. Acknowledge everything else so
	// Telegram does not keep retrying an update we will never store.
	if req.Message == nil || req.Message.Text == "" {
		w.WriteHeader(http.StatusOK)
		return
	}

	result, err := e.handler.Handle(r.Context(), Command{
		TenantID:          e.cfg.TenantID,
		ProviderMessageID: strconv.FormatInt(req.UpdateID, 10),
		ExternalChatID:    strconv.FormatInt(req.Message.Chat.ID, 10),
		Text:              req.Message.Text,
	})

	switch {
	case err == nil:
		e.log.Info("telegram update received",
			"update_id", req.UpdateID,
			"inbound_message_id", result.InboundMessageID,
			"duplicate", result.Duplicate)
		w.WriteHeader(http.StatusOK)
	case cannotBeFixedByRetry(err):
		// Retrying the same payload cannot succeed; 200 stops Telegram's retries.
		e.log.Warn("telegram update rejected", "update_id", req.UpdateID, "error", err)
		w.WriteHeader(http.StatusOK)
	default:
		// Not stored: a non-2xx makes Telegram retry later, which is what we want.
		e.log.Error("telegram update not stored", "update_id", req.UpdateID, "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}

// cannotBeFixedByRetry: invalid input or a broken business rule would fail again on every redelivery.
func cannotBeFixedByRetry(err error) bool {
	var violation *seedwork.BusinessRuleViolation
	return errors.Is(err, ErrInvalidCommand) || errors.As(err, &violation)
}
