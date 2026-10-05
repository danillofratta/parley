package seedwork

import "errors"

// Standard repository outcomes, with the same meaning in every service and stack.
var (
	ErrNotFound            = errors.New("aggregate not found")
	ErrAlreadyExists       = errors.New("aggregate already exists")
	ErrConcurrencyConflict = errors.New("aggregate was changed by another operation")
)
