package pg

import (
	"context"
	_ "embed"

	"github.com/jackc/pgx/v5/pgxpool"

	domain "go-svr/health"
)

//go:embed queries/health/get_health.sql
var getHealthSQL string

//go:embed queries/health/create_ping.sql
var createPingSQL string

type healthRepository struct {
	pool *pgxpool.Pool
}

// NewHealthRepository creates a PostgreSQL-backed health.Repository.
func NewHealthRepository(pool *pgxpool.Pool) domain.Repository {
	return &healthRepository{pool: pool}
}

func (r *healthRepository) GetHealth(ctx context.Context) (domain.Health, error) {
	row := conn(ctx, r.pool).QueryRow(ctx, getHealthSQL)

	var h domain.Health
	if err := row.Scan(&h.Status, &h.CheckedAtUnix); err != nil {
		return domain.Health{}, mapErr(err, domain.ErrNotFound, domain.ErrConflict)
	}
	return h, nil
}

func (r *healthRepository) Ping(ctx context.Context, message string) (domain.Ping, error) {
	row := conn(ctx, r.pool).QueryRow(ctx, createPingSQL, message)

	var p domain.Ping
	if err := row.Scan(&p.Message, &p.ReceivedAtUnix); err != nil {
		return domain.Ping{}, mapErr(err, domain.ErrNotFound, domain.ErrConflict)
	}
	return p, nil
}
