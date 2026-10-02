package repositories

import (
	"context"

	"github.com/parley/parley/services/go/channel-gateway/internal/domain/entities"
)

type InboundMessageRepository interface {
	Add(ctx context.Context, message *entities.InboundMessage) error
	FindByID(ctx context.Context, id string) (*entities.InboundMessage, error)
}