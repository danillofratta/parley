package persistence

import "time"

// Auditable is implemented by every entity through the seedwork.
type Auditable interface {
	StampCreated(at time.Time, by string)
	StampUpdated(at time.Time, by string)
}

// Auditor stamps audit fields with the current time and the acting user or service.
type Auditor struct {
	now   func() time.Time
	actor string
}

func NewAuditor(now func() time.Time, actor string) Auditor {
	return Auditor{now: now, actor: actor}
}

func (a Auditor) Created(e Auditable) { e.StampCreated(a.now().UTC(), a.actor) }
func (a Auditor) Updated(e Auditable) { e.StampUpdated(a.now().UTC(), a.actor) }
