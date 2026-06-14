package pg

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// pgxConn is the minimal interface satisfied by both *pgxpool.Pool and pgx.Tx.
// Keeping this interface narrow means every repository method works transparently
// inside or outside an active transaction.
type pgxConn interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

type txKey struct{}

// conn returns the active transaction from ctx if one was stored via withTx,
// otherwise falls back to pool. Every repository method must use conn() for
// all query execution.
//
// For cross-domain transactions, store the active pgx.Tx in ctx using
// context.WithValue(ctx, txKey{}, tx) before passing ctx to the callback.
// See db/README.md for the full pattern and known risks.
func conn(ctx context.Context, pool *pgxpool.Pool) pgxConn {
	if tx, ok := ctx.Value(txKey{}).(pgx.Tx); ok {
		return tx
	}
	return pool
}
