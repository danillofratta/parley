package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/danillofratta/parley/services/go/channel-gateway/internal/domain/entities"
	"github.com/danillofratta/parley/services/go/channel-gateway/internal/domain/seedwork"
	"github.com/danillofratta/parley/services/go/channel-gateway/internal/infrastructure/outbox"
	"github.com/danillofratta/parley/services/go/channel-gateway/internal/infrastructure/persistence"
)

type InboundMessageRepository struct {
	pool    *pgxpool.Pool
	auditor persistence.Auditor
}

func NewInboundMessageRepository(pool *pgxpool.Pool, auditor persistence.Auditor) *InboundMessageRepository {
	return &InboundMessageRepository{pool: pool, auditor: auditor}
}

// Add stores the message and writes its events to the outbox in one transaction.
// The unique constraint on (tenant, channel, provider message id) acts as the Inbox:
// a redelivered message conflicts, writes nothing and produces no new event.
func (r *InboundMessageRepository) Add(ctx context.Context, message *entities.InboundMessage) error {
	r.auditor.Created(message)

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }() // no-op after a successful commit

	tag, err := tx.Exec(ctx, `
		INSERT INTO gateway.inbound_messages
			(id, tenant_id, channel, contact_chat, provider_message_id, received_at,
			 created_at, created_by, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (tenant_id, channel, provider_message_id) DO NOTHING`,
		message.ID(),
		message.TenantID(),
		message.Sender().Channel().String(),
		message.Sender().Chat(),
		message.ProviderMessageID(),
		message.ReceivedAt(),
		message.CreatedAt(),
		message.CreatedBy(),
		message.Version(),
	)
	if err != nil {
		return fmt.Errorf("insert inbound message: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return seedwork.ErrAlreadyExists
	}

	for _, event := range message.DomainEvents() {
		msg, err := outbox.FromDomainEvent(event)
		if err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, `
			INSERT INTO gateway.outbox_messages
				(id, tenant_id, type, topic, message_key, payload, occurred_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			msg.ID, msg.TenantID, msg.Type, msg.Topic, msg.Key, msg.Payload, msg.OccurredAt,
		); err != nil {
			return fmt.Errorf("insert outbox message: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit: %w", err)
	}
	message.ClearDomainEvents()
	return nil
}
