package seedwork

import "github.com/google/uuid"

// AggregateRoot is an entity that owns a consistency boundary.
// It belongs to one tenant and records the domain events raised by its behaviour.
type AggregateRoot struct {
	Entity
	tenantID uuid.UUID
	events   []DomainEvent
}

// NewAggregateRoot creates the base of a new aggregate root for a tenant.
func NewAggregateRoot(tenantID uuid.UUID) (AggregateRoot, error) {
	if err := CheckRule(TenantMustBeInformed{TenantID: tenantID}); err != nil {
		return AggregateRoot{}, err
	}
	return AggregateRoot{Entity: NewEntity(), tenantID: tenantID}, nil
}

// RestoreAggregateRoot rebuilds the base of an aggregate root loaded from storage.
func RestoreAggregateRoot(entity Entity, tenantID uuid.UUID) AggregateRoot {
	return AggregateRoot{Entity: entity, tenantID: tenantID}
}

func (a *AggregateRoot) TenantID() uuid.UUID { return a.tenantID }

// Raise records a domain event. The repository writes it to the outbox when saving.
func (a *AggregateRoot) Raise(event DomainEvent) {
	a.events = append(a.events, event)
}

// DomainEvents returns the events raised since the aggregate was created or loaded.
func (a *AggregateRoot) DomainEvents() []DomainEvent {
	return append([]DomainEvent(nil), a.events...)
}

// ClearDomainEvents is called by the repository after a successful commit.
func (a *AggregateRoot) ClearDomainEvents() {
	a.events = nil
}

// TenantMustBeInformed: every aggregate belongs to exactly one tenant.
type TenantMustBeInformed struct {
	TenantID uuid.UUID
}

func (r TenantMustBeInformed) IsBroken() bool { return r.TenantID == uuid.Nil }
func (TenantMustBeInformed) Code() string     { return "aggregate.tenant_required" }
func (TenantMustBeInformed) Message() string  { return "every aggregate must belong to a tenant" }
