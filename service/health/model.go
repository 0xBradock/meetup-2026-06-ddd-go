// Package health is the core domain for service liveness and connectivity checks.
// It defines the domain model, the repository port, and the application service.
// No persistence or transport details belong here.
package health

import (
	"errors"
	"strings"
	"time"
)

// ErrEmptyPingMessage is returned by NewPing when the message is blank or
// contains only whitespace.
var ErrEmptyPingMessage = errors.New("ping message is empty")

// Health represents the current health status of the service at a point in time.
type Health struct {
	// Status is a human-readable health state, typically "ok".
	Status string
	// CheckedAtUnix is the UTC Unix timestamp (seconds) when the check was performed.
	CheckedAtUnix int64
}

// Ping is a record of a received ping message.
type Ping struct {
	// Message is the trimmed content of the ping.
	Message string
	// ReceivedAtUnix is the UTC Unix timestamp (seconds) when the ping was received.
	ReceivedAtUnix int64
}

// NewPing validates message and returns a Ping stamped with the current time.
// It returns ErrEmptyPingMessage if message is blank or whitespace-only.
func NewPing(message string) (Ping, error) {
	trimmed := strings.TrimSpace(message)
	if trimmed == "" {
		return Ping{}, ErrEmptyPingMessage
	}

	return Ping{
		Message:        trimmed,
		ReceivedAtUnix: time.Now().Unix(),
	}, nil
}
