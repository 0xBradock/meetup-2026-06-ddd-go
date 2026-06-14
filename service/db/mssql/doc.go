// Package mssql implements domain repository interfaces against Microsoft SQL
// Server using database/sql and the go-mssqldb driver. SQL queries are embedded
// from the queries/ directory at compile time. No domain package imports this
// package — the dependency always flows inward.
package mssql
