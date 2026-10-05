package seedwork

import (
	"time"

	"github.com/google/uuid"
)

// DomainEvent is a business fact that already happened, named in past tense.
type DomainEvent interface {
	EventID() uuid.UUID
	OccurredAt() time.Time
	AggregateID() uuid.UUID
	TenantID() uuid.UUID
}

// EventBase holds the fields every domain event has. Embed it in concrete events.
type EventBase struct {
	eventID     uuid.UUID
	occurredAt  time.Time
	aggregateID uuid.UUID
	tenantID    uuid.UUID
}

func NewEventBase(aggregateID, tenantID uuid.UUID, occurredAt time.Time) EventBase {
	return EventBase{eventID: NewID(), occurredAt: occurredAt, aggregateID: aggregateID, tenantID: tenantID}
}

func (e EventBase) EventID() uuid.UUID     { return e.eventID }
func (e EventBase) OccurredAt() time.Time  { return e.occurredAt }
func (e EventBase) AggregateID() uuid.UUID { return e.aggregateID }
func (e EventBase) TenantID() uuid.UUID    { return e.tenantID }
