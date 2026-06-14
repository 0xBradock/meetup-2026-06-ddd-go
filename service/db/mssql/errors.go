package mssql

import (
	"database/sql"
	"errors"
	"fmt"

	mssqldb "github.com/microsoft/go-mssqldb"
)

const (
	mssqlErrUniqueViolation = 2627 // Violation of UNIQUE KEY constraint
	mssqlErrDuplicateKey    = 2601 // Cannot insert duplicate key row
	mssqlErrForeignKey      = 547  // INSERT statement conflicted with FOREIGN KEY constraint
)

// mapErr translates MSSQL errors into domain-level sentinel errors.
// Pass the domain-specific sentinels as errNotFound and errConflict so this
// function stays decoupled from any particular domain package.
func mapErr(err, errNotFound, errConflict error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w", errNotFound)
	}

	var mssqlErr mssqldb.Error
	if errors.As(err, &mssqlErr) {
		switch mssqlErr.Number {
		case mssqlErrUniqueViolation, mssqlErrDuplicateKey:
			return fmt.Errorf("%w", errConflict)
		case mssqlErrForeignKey:
			return fmt.Errorf("%w", errNotFound)
		}
	}

	return err
}
