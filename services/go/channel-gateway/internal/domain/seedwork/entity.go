package seedwork

import (
	"time"

	"github.com/google/uuid"
)

// Entity is the base of every entity: identity, audit trail and version.
// Embed it (through AggregateRoot for roots). Entities are equal when their ids are equal.
type Entity struct {
	id        uuid.UUID
	createdAt time.Time
	createdBy string
	updatedAt *time.Time
	updatedBy *string
	version   int
}

// NewEntity creates the base of a new entity with a fresh id.
func NewEntity() Entity {
	return Entity{id: NewID()}
}

// RestoreEntity rebuilds the base of an entity loaded from storage.
func RestoreEntity(id uuid.UUID, createdAt time.Time, createdBy string,
	updatedAt *time.Time, updatedBy *string, version int) Entity {
	return Entity{id: id, createdAt: createdAt, createdBy: createdBy,
		updatedAt: updatedAt, updatedBy: updatedBy, version: version}
}

func (e *Entity) ID() uuid.UUID         { return e.id }
func (e *Entity) CreatedAt() time.Time  { return e.createdAt }
func (e *Entity) CreatedBy() string     { return e.createdBy }
func (e *Entity) UpdatedAt() *time.Time { return e.updatedAt }
func (e *Entity) UpdatedBy() *string    { return e.updatedBy }
func (e *Entity) Version() int          { return e.version }

// SameIdentityAs compares entities by identity, never by attributes.
func (e *Entity) SameIdentityAs(other *Entity) bool {
	return other != nil && e.id == other.id
}

// StampCreated is called by persistence (the Auditor) on the first save.
func (e *Entity) StampCreated(at time.Time, by string) {
	e.createdAt = at
	e.createdBy = by
	e.version = 1
}

// StampUpdated is called by persistence (the Auditor) on every update.
func (e *Entity) StampUpdated(at time.Time, by string) {
	e.updatedAt = &at
	e.updatedBy = &by
	e.version++
}
