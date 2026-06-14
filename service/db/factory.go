// Package db is the composition root for all persistence adapters. It exposes
// factory functions that return domain repositories backed by a real database
// engine or an in-memory fallback, depending on the active driver. cmd/ is the
// only caller — it wires the repositories returned here into domain services.
package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"go-svr/db/inmem"
	dbmssql "go-svr/db/mssql"
	"go-svr/db/pg"
	domain "go-svr/health"
)

// PGRepositories holds domain repositories backed by PostgreSQL,
// or in-memory when the postgres driver is not active.
type PGRepositories struct {
	Health domain.Repository
}

// NewPG returns domain repositories backed by PostgreSQL when driver is
// "postgres", or in-memory for any other driver value. ctx is used for the
// initial connection attempt and should be the application lifecycle context.
// The caller is responsible for invoking the returned cleanup function on shutdown.
func NewPG(ctx context.Context, driver string, getenv func(string) string) (*PGRepositories, func(), error) {
	if driver != "postgres" {
		return &PGRepositories{
			Health: inmem.NewRepository(),
		}, func() { /* no-op: in-memory backend requires no cleanup */ }, nil
	}

	databaseURL := getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, nil, fmt.Errorf("DATABASE_URL required for postgres driver")
	}

	poolConfig, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, nil, fmt.Errorf("parse database url: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("connect postgres: %w", err)
	}

	return &PGRepositories{
		Health: pg.NewHealthRepository(pool),
	}, pool.Close, nil
}

// MSSQLRepositories holds domain repositories backed by Microsoft SQL Server,
// or in-memory when the mssql driver is not active.
type MSSQLRepositories struct {
	Health domain.Repository
}

// NewMSSQL returns domain repositories backed by Microsoft SQL Server when
// driver is "mssql", or in-memory for any other driver value. ctx is used for
// the initial connection ping and should be the application lifecycle context.
// The caller is responsible for invoking the returned cleanup function on shutdown.
func NewMSSQL(ctx context.Context, driver string, getenv func(string) string) (*MSSQLRepositories, func(), error) {
	if driver != "mssql" {
		return &MSSQLRepositories{
			Health: inmem.NewRepository(),
		}, func() { /* no-op: in-memory backend requires no cleanup */ }, nil
	}

	dsn := getenv("MSSQL_DATABASE_URL")
	if dsn == "" {
		return nil, nil, fmt.Errorf("MSSQL_DATABASE_URL required for mssql driver")
	}

	db, err := sql.Open("sqlserver", dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("open mssql connection: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, nil, fmt.Errorf("ping mssql: %w", err)
	}

	return &MSSQLRepositories{
		Health: dbmssql.NewHealthRepository(db),
	}, func() { _ = db.Close() }, nil
}
