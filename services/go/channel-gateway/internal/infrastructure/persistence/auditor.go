package persistence

import "time"

type Auditable interface {
	CreatedAt() time.Time
	UpdatedAt() time.Time
}

type Auditor struct {
	createdAt time.Time
	updatedAt time.Time
}

func (a Auditor) CreatedAt() time.Time {
	return a.createdAt
}

func (a Auditor) UpdatedAt() time.Time {
	return a.updatedAt
}

func NewAuditor(createdAt, updatedAt time.Time) Auditor {
	return Auditor{
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}