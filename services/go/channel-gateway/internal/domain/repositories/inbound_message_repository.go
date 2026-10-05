package repositories

import (
	"context"

	"github.com/danillofratta/parley/services/go/channel-gateway/internal/domain/entities"
)

// InboundMessageRepository is the persistence port of InboundMessage.
// The standard defines GetById, Add and Update; a repository exposes only what use cases need.
type InboundMessageRepository interface {
	// Add stores a new message and writes its events to the outbox in one transaction.
	// It returns seedwork.ErrAlreadyExists when the provider message was already received.
	Add(ctx context.Context, message *entities.InboundMessage) error
}
