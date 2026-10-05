package outbox

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// relayLockKey identifies the gateway's outbox relay in pg_try_advisory_xact_lock.
// Only the instance holding it publishes in a cycle, which keeps per-key order across replicas.
const relayLockKey int64 = 4_201_001

// PendingMessage is an outbox row waiting to be published.
type PendingMessage struct {
	ID      string
	Topic   string
	Key     string
	Payload []byte
	Headers map[string]string
}

// Publisher sends one message to the broker and returns only after the broker acknowledged it.
type Publisher interface {
	Publish(ctx context.Context, msg PendingMessage) error
}

type RelayConfig struct {
	PollInterval   time.Duration
	BatchSize      int
	PublishTimeout time.Duration
}

// Relay publishes outbox rows in occurred_at order and marks them as published.
type Relay struct {
	pool      *pgxpool.Pool
	publisher Publisher
	cfg       RelayConfig
	log       *slog.Logger
}

func NewRelay(pool *pgxpool.Pool, publisher Publisher, cfg RelayConfig, log *slog.Logger) *Relay {
	return &Relay{pool: pool, publisher: publisher, cfg: cfg, log: log}
}

// Run polls the outbox until ctx is cancelled. It returns nil on a clean shutdown.
func (r *Relay) Run(ctx context.Context) error {
	r.log.Info("outbox relay started", "poll_interval", r.cfg.PollInterval, "batch_size", r.cfg.BatchSize)
	ticker := time.NewTicker(r.cfg.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			r.log.Info("outbox relay stopped")
			return nil
		case <-ticker.C:
			published, err := r.publishBatch(ctx)
			if err != nil && ctx.Err() == nil {
				r.log.Error("outbox relay cycle failed", "error", err)
			}
			if published > 0 {
				r.log.Info("outbox messages published", "count", published)
			}
		}
	}
}

// publishBatch runs one cycle in a single transaction:
// take the relay lock, read pending rows in order, publish until the first failure,
// then mark what was published (and count the failed attempt).
func (r *Relay) publishBatch(ctx context.Context) (int, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }() // no-op after a successful commit

	var leader bool
	if err := tx.QueryRow(ctx, `SELECT pg_try_advisory_xact_lock($1)`, relayLockKey).Scan(&leader); err != nil {
		return 0, fmt.Errorf("acquire relay lock: %w", err)
	}
	if !leader {
		return 0, nil // another instance is publishing in this cycle
	}

	rows, err := tx.Query(ctx, `
		SELECT id, topic, message_key, payload, headers
		FROM gateway.outbox_messages
		WHERE published_at IS NULL
		ORDER BY occurred_at, id
		LIMIT $1`, r.cfg.BatchSize)
	if err != nil {
		return 0, fmt.Errorf("select pending messages: %w", err)
	}
	pending, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (PendingMessage, error) {
		var m PendingMessage
		err := row.Scan(&m.ID, &m.Topic, &m.Key, &m.Payload, &m.Headers)
		return m, err
	})
	if err != nil {
		return 0, fmt.Errorf("read pending messages: %w", err)
	}

	published := make([]string, 0, len(pending))
	for _, msg := range pending {
		publishCtx, cancel := context.WithTimeout(ctx, r.cfg.PublishTimeout)
		err := r.publisher.Publish(publishCtx, msg)
		cancel()
		if err != nil {
			// Stop here: publishing later rows would break the order of this conversation.
			r.log.Warn("outbox publish failed", "outbox_id", msg.ID, "topic", msg.Topic, "error", err)
			if _, err := tx.Exec(ctx,
				`UPDATE gateway.outbox_messages SET attempts = attempts + 1 WHERE id = $1`, msg.ID); err != nil {
				return 0, fmt.Errorf("count failed attempt: %w", err)
			}
			break
		}
		published = append(published, msg.ID)
	}

	if len(published) > 0 {
		if _, err := tx.Exec(ctx,
			`UPDATE gateway.outbox_messages SET published_at = now() WHERE id = ANY($1::uuid[])`,
			published); err != nil {
			return 0, fmt.Errorf("mark published: %w", err)
		}
	}

	// If this commit fails after Kafka acknowledged, the rows are published again
	// next cycle: at-least-once delivery, absorbed by consumers' Inbox.
	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("commit: %w", err)
	}
	return len(published), nil
}
