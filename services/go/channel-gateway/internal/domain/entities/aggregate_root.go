package entities

import (
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/parley/parley/services/go/channel-gateway/internal/domain/events"
)

type AggregateRoot struct {
	id           string
	tenantID     string
	domainEvents []events.DomainEvent
}

func NewAggregateRoot(tenantID string) (*AggregateRoot, error) {
	tenantID = strings.TrimSpace(tenantID)
	if tenantID == "" {
		return nil, fmt.Errorf("tenant id is required")
	}

	id, err := newID()
	if err != nil {
		return nil, fmt.Errorf("generate aggregate id: %w", err)
	}

	return &AggregateRoot{
		id:       id,
		tenantID: tenantID,
	}, nil
}

func (a *AggregateRoot) ID() string {
	return a.id
}

func (a *AggregateRoot) TenantID() string {
	return a.tenantID
}

func (a *AggregateRoot) Raise(event events.DomainEvent) {
	if event == nil {
		return
	}
	a.domainEvents = append(a.domainEvents, event)
}

func (a *AggregateRoot) DomainEvents() []events.DomainEvent {
	if len(a.domainEvents) == 0 {
		return nil
	}

	out := make([]events.DomainEvent, len(a.domainEvents))
	copy(out, a.domainEvents)
	return out
}

func (a *AggregateRoot) ClearDomainEvents() {
	a.domainEvents = nil
}

func newID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}