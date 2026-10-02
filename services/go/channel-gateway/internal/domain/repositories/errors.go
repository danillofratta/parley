package repositories

import "errors"

// ErrAlreadyExists is returned by Add when the aggregate (or its provider
// message id for the same tenant) is already stored. Handlers treat it as a
// duplicate delivery, not as a failure.
var ErrAlreadyExists = errors.New("inbound message already exists")
