// Package pg implements domain repository interfaces against PostgreSQL using
// pgx/v5. SQL queries are embedded from the queries/ directory at compile time.
// No domain package imports this package — the dependency always flows inward.
package pg
