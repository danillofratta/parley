# Go slice template (Channel Gateway: receive a Telegram webhook)

```
internal/features/receivewebhook/
  request.go        // Telegram update as received over HTTP
  command.go        // application input
  result.go         // application output
  validator.go      // boundary validation of the command
  handler.go        // use case orchestration
  endpoint.go       // HTTP: auth, decode Request, map, call handler
  handler_test.go
  endpoint_test.go
```
No `response.go`: the webhook answers with a status code only.

```go
// request.go
package receivewebhook

type Request struct {
	UpdateID int64           `json:"update_id"`
	Message  *RequestMessage `json:"message"`
}

type RequestMessage struct {
	Chat RequestChat `json:"chat"`
	Text string      `json:"text"`
}

type RequestChat struct {
	ID int64 `json:"id"`
}
```

```go
// command.go
package receivewebhook

type Command struct {
	TenantID          string
	ProviderMessageID string
	ExternalChatID    string
	Text              string
	RawPayload        []byte
}
```

```go
// result.go
package receivewebhook

type Result struct {
	Duplicate bool
}
```

```go
// validator.go
package receivewebhook

import "errors"

var ErrInvalidCommand = errors.New("invalid command")

func validate(cmd Command) error {
	if cmd.TenantID == "" || cmd.ProviderMessageID == "" || cmd.ExternalChatID == "" {
		return ErrInvalidCommand
	}
	return nil
}
```

```go
// handler.go
package receivewebhook

type Handler struct {
	records domain.InboundRecordRepository // port declared in domain
}

func NewHandler(records domain.InboundRecordRepository) *Handler {
	return &Handler{records: records}
}

func (h *Handler) Handle(ctx context.Context, cmd Command) (Result, error) {
	if err := validate(cmd); err != nil {
		return Result{}, err
	}
	record, err := domain.NewInboundRecord(cmd.TenantID, domain.ChannelTelegram,
		cmd.ProviderMessageID, cmd.ExternalChatID, cmd.Text, cmd.RawPayload)
	if err != nil {
		return Result{}, err
	}
	stored, err := h.records.SaveOnce(ctx, record) // record + outbox, one transaction
	if err != nil {
		return Result{}, fmt.Errorf("save inbound record: %w", err)
	}
	return Result{Duplicate: !stored}, nil
}
```

```go
// endpoint.go
package receivewebhook

type Endpoint struct {
	handler  *Handler // called directly: no mediator
	secret   string
	tenantID string
	log      *slog.Logger
}

func NewEndpoint(handler *Handler, secret, tenantID string, log *slog.Logger) *Endpoint {
	return &Endpoint{handler: handler, secret: secret, tenantID: tenantID, log: log}
}

func (e *Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 1. authenticate (constant-time), 2. limit and read body, 3. decode Request,
	// 4. map Request -> Command, 5. call e.handler.Handle, 6. map error/result -> status
}
```

```go
// internal/bootstrap/bootstrap.go (composition root)
func New(cfg config.Config, pool *pgxpool.Pool, log *slog.Logger) http.Handler {
	records := postgres.NewInboundRecordRepository(pool)

	receiveWebhook := receivewebhook.NewEndpoint(
		receivewebhook.NewHandler(records), cfg.WebhookSecret, cfg.TenantID, log)

	mux := http.NewServeMux()
	mux.Handle("POST /webhooks/telegram", receiveWebhook)
	return mux
}
```
