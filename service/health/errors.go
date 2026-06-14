package health

import "errors"

var (
	// ErrNotFound is returned when a requested resource does not exist.
	ErrNotFound = errors.New("not found")
	// ErrConflict is returned when an operation would violate a uniqueness constraint.
	ErrConflict = errors.New("conflict")
)
