// Package inmem provides in-memory implementations of all domain repository
// interfaces. It is used during local development and in unit tests where a
// real database is not needed or available.
package inmem

import (
	"context"
	"sync"
	"time"

	domain "go-svr/health"
)

// Repository is a thread-safe in-memory implementation of health.Repository.
type Repository struct {
	mu       sync.RWMutex
	lastPing domain.Ping
}

// NewRepository returns an empty in-memory Repository ready for use.
func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) GetHealth(context.Context) (domain.Health, error) {
	return domain.Health{
		Status:        "ok",
		CheckedAtUnix: time.Now().Unix(),
	}, nil
}

func (r *Repository) Ping(_ context.Context, message string) (domain.Ping, error) {
	ping, err := domain.NewPing(message)
	if err != nil {
		return domain.Ping{}, err
	}

	r.mu.Lock()
	r.lastPing = ping
	r.mu.Unlock()

	return ping, nil
}
