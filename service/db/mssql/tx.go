package mssql

import (
	"context"
	"database/sql"
)

// mssqlConn is the minimal interface satisfied by both *sql.DB and *sql.Tx.
// Keeping this interface narrow means every repository method works transparently
// inside or outside an active transaction.
type mssqlConn interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type txKey struct{}

// conn returns the active transaction from ctx if one was stored via withTx,
// otherwise falls back to db. Every repository method must use conn() for
// all query execution.
//
// For cross-domain transactions, store the active *sql.Tx in ctx using
// context.WithValue(ctx, txKey{}, tx) before passing ctx to the callback.
// See db/README.md for the full pattern and known risks.
func conn(ctx context.Context, db *sql.DB) mssqlConn {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return tx
	}
	return db
}
