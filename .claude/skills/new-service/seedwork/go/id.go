package seedwork

import "github.com/google/uuid"

// NewID returns a new time-ordered identity (UUID v7) for entities and events.
func NewID() uuid.UUID {
	return uuid.Must(uuid.NewV7())
}
