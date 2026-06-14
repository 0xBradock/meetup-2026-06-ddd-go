package pg

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	pgErrUniqueViolation     = "23505"
	pgErrForeignKeyViolation = "23503"
)

// mapErr translates pgconn errors into domain-level sentinel errors.
// Pass the domain-specific sentinels as errNotFound and errConflict so this
// function stays decoupled from any particular domain package.
//
// Usage in a repository method:
//
//	return domain.Ping{}, mapErr(err, domain.ErrNotFound, domain.ErrConflict)
func mapErr(err, errNotFound, errConflict error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("%w", errNotFound)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgErrUniqueViolation:
			return fmt.Errorf("%w", errConflict)
		case pgErrForeignKeyViolation:
			return fmt.Errorf("%w", errNotFound)
		}
	}

	return err
}
