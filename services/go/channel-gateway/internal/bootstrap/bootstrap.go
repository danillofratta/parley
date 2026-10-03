package bootstrap

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/danillofratta/parley/services/go/channel-gateway/internal/features/receivetelegramupdate"
	"github.com/danillofratta/parley/services/go/channel-gateway/internal/infrastructure/clock"
	"github.com/danillofratta/parley/services/go/channel-gateway/internal/infrastructure/config"
	"github.com/danillofratta/parley/services/go/channel-gateway/internal/infrastructure/persistence"
	"github.com/danillofratta/parley/services/go/channel-gateway/internal/infrastructure/postgres"
)

const actor = "system:channel-gateway"

// New builds the dependency graph. It is the only place that knows concrete types.
func New(cfg config.Config, pool *pgxpool.Pool, log *slog.Logger) (http.Handler, error) {
	// Fail closed: a misconfigured tenant must stop the service, not drop messages.
	if _, err := uuid.Parse(cfg.TenantID); err != nil {
		return nil, fmt.Errorf("DEFAULT_TENANT_ID: %w", err)
	}

	auditor := persistence.NewAuditor(time.Now, actor)
	messages := postgres.NewInboundMessageRepository(pool, auditor)

	receiveTelegramUpdate := receivetelegramupdate.NewEndpoint(
		receivetelegramupdate.NewHandler(messages, clock.System{}),
		receivetelegramupdate.EndpointConfig{
			WebhookSecret: cfg.WebhookSecret,
			TenantID:      cfg.TenantID,
		},
		log,
	)

	mux := http.NewServeMux()
	mux.Handle("POST /webhooks/telegram", receiveTelegramUpdate)
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		if err := pool.Ping(r.Context()); err != nil {
			http.Error(w, "database unavailable", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	return mux, nil
}