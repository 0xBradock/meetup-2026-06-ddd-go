package mssql

import (
	"context"
	"database/sql"
	_ "embed"

	_ "github.com/microsoft/go-mssqldb" // register sqlserver driver

	domain "go-svr/health"
)

//go:embed queries/health/get_health.sql
var getHealthSQL string

//go:embed queries/health/create_ping.sql
var createPingSQL string

type healthRepository struct {
	db *sql.DB
}

// NewHealthRepository creates a Microsoft SQL Server-backed health.Repository.
func NewHealthRepository(db *sql.DB) domain.Repository {
	return &healthRepository{db: db}
}

func (r *healthRepository) GetHealth(ctx context.Context) (domain.Health, error) {
	row := conn(ctx, r.db).QueryRowContext(ctx, getHealthSQL)

	var h domain.Health
	if err := row.Scan(&h.Status, &h.CheckedAtUnix); err != nil {
		return domain.Health{}, mapErr(err, domain.ErrNotFound, domain.ErrConflict)
	}
	return h, nil
}

func (r *healthRepository) Ping(ctx context.Context, message string) (domain.Ping, error) {
	row := conn(ctx, r.db).QueryRowContext(ctx, createPingSQL, message)

	var p domain.Ping
	if err := row.Scan(&p.Message, &p.ReceivedAtUnix); err != nil {
		return domain.Ping{}, mapErr(err, domain.ErrNotFound, domain.ErrConflict)
	}
	return p, nil
}
