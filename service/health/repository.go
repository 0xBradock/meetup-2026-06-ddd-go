package health

import "context"

// Repository is the persistence port for the health domain.
// Implementations live in db/pg, db/mssql, and db/inmem — never here.
type Repository interface {
	// GetHealth returns the current health status of the service.
	GetHealth(ctx context.Context) (Health, error)
	// Ping stores a ping message and returns it with a received timestamp.
	Ping(ctx context.Context, message string) (Ping, error)
}
