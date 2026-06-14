# db

Persistence adapters for all domain repositories. This package translates between domain interfaces and database engines. No domain package imports anything from here — the dependency always flows inward (db → domain, never domain → db).

## Structure

```
db/
  factory.go          — NewPG / NewMSSQL: construct domain repos for a given driver
  inmem/              — In-memory adapter (used in tests and when no real DB is needed)
    repository.go     — Implements health.Repository in-memory
    cache.go          — Implements health.CachePort in-memory (for tests / local dev)
  pg/                 — PostgreSQL adapter (pgx/v5)
    queries/<domain>/ — One .sql file per repository method, embedded at compile time
    migrations/       — golang-migrate SQL migration files
    tx.go             — Transaction threading via context (see Transactions below)
    errors.go         — pg error code → domain sentinel mapper
    health.go         — Implements health.Repository against PostgreSQL
  mssql/              — Microsoft SQL Server adapter (go-mssqldb / database/sql)
    queries/<domain>/ — T-SQL files, same layout as pg/queries/
    migrations/       — golang-migrate SQL migration files
    tx.go             — Transaction threading via context (same contract as pg/tx.go)
    errors.go         — MSSQL error number → domain sentinel mapper
    health.go         — Implements health.Repository against MSSQL
  valkey/             — Valkey cache adapter (go-redis/v9, wire-compatible with Valkey)
    cache.go          — Implements health.CachePort backed by Valkey
    health.go         — cachedHealthRepo: decorator that wraps any health.Repository with a cache layer
```

## Design decisions

### Raw SQL, no ORM, no query builder

All SQL is hand-written. Queries live in `.sql` files under each engine's `queries/<domain>/` directory and are embedded into the binary at compile time via `//go:embed`. If a file is missing the build fails — there is no runtime file loading.

Each file contains exactly one query. The column order in `SELECT` / `RETURNING` / `OUTPUT` is the contract: the matching `rows.Scan` call in the Go repository method must list fields in the same order.

### Mapping rows to domain types

Repository methods are responsible for translating database rows into domain types. The general pattern is:

1. Run one or more queries against the database.
2. Scan results into local variables or a db-package-internal struct — never into a type imported from a db driver.
3. Construct and return the domain type explicitly, or an error.

```go
// db/pg/order.go — example with an intermediate row struct
type orderRow struct {
    ID         int64
    CustomerID int64
    TotalCents int64
    State      string
    CreatedAt  int64
}

func (r *orderRepository) GetOrder(ctx context.Context, id int64) (domain.Order, error) {
    var row orderRow
    err := conn(ctx, r.pool).QueryRow(ctx, getOrderSQL, id).
        Scan(&row.ID, &row.CustomerID, &row.TotalCents, &row.State, &row.CreatedAt)
    if err != nil {
        return domain.Order{}, mapErr(err, domain.ErrNotFound, domain.ErrConflict)
    }
    return domain.Order{
        ID:         row.ID,
        CustomerID: row.CustomerID,
        Total:      money.FromCents(row.TotalCents),
        State:      domain.OrderState(row.State),
        CreatedAt:  row.CreatedAt,
    }, nil
}
```

The `health` domain is a degenerate case where the query columns and domain fields are identical primitives, so the intermediate struct collapses to a direct scan. Do not treat that as the general pattern.

Nullable columns are represented as pointer fields (`*string`, `*int64`, `*time.Time`) in the scan target; pgx v5 and `database/sql` both set these to `nil` on a SQL `NULL`.

**Constraint:** domain structs must contain only Go stdlib types. `pgtype.*` or any other db-package type must never appear in a domain package.

### Factories

`NewPG(ctx, driver, getenv)` and `NewMSSQL(ctx, driver, getenv)` each return a typed struct of domain repositories:

```go
// driver="postgres"  → PostgreSQL-backed repos
// driver=anything    → in-memory repos
pgRepos, cleanup, err := db.NewPG(ctx, driver, os.Getenv)
defer cleanup()

// driver="mssql"     → MSSQL-backed repos
// driver=anything    → in-memory repos
msRepos, cleanup, err := db.NewMSSQL(ctx, driver, os.Getenv)
defer cleanup()
```

`cmd/` is the composition root — it calls the factory that matches the active driver and wires each domain service with the right repository. Neither factory knows about the other.

`ctx` is the application lifecycle context (from `Run`). It is used for the initial connection attempt so that a shutdown signal during startup cancels the connection rather than hanging.

### Error mapping

Each engine package has an `errors.go` with a `mapErr(err, errNotFound, errConflict error) error` function. It translates engine-specific error codes (pg `23505`, MSSQL `2627`, etc.) into the domain sentinel errors passed by the caller:

```go
// in a repository method:
return domain.Ping{}, mapErr(err, domain.ErrNotFound, domain.ErrConflict)
```

`mapErr` has no import of any domain package — the caller passes the sentinels, keeping the mapper reusable across domains.

Domain packages define their own sentinels:

```go
// health/errors.go
var ErrNotFound = errors.New("not found")
var ErrConflict  = errors.New("conflict")
```

Callers check errors with `errors.Is(err, health.ErrNotFound)`.

## Cache layer

`db/valkey` provides a `cachedHealthRepo` decorator that wraps any `health.Repository` and caches `GetHealth` results in Valkey. It is wired in `cmd/server/main.go` after the factory, so the factory functions stay single-responsibility:

```go
// healthRepo is set by the active driver branch (pg, mssql, or inmemory)
if cfg.Cache.URL != "" {
    c, err := valkey.NewCache(ctx, cfg.Cache.URL, cfg.Cache.Prefix, cfg.Cache.TLSCACert)
    // ...
    healthRepo = valkey.NewCachedHealthRepo(healthRepo, c, cfg.Cache.HealthTTL, logger)
}
hs := health.NewService(healthRepo)
```

The cache is **optional** — omitting `CACHE_URL` skips the wrapping step entirely. On cache errors the decorator degrades gracefully by falling through to the inner repository.

`health.CachePort` (the interface the decorator consumes) is defined in the `health` package, following the hexagonal convention that ports are owned by the domain that uses them.

See `dbd/valkey/README.md` for the full configuration reference and key strategy.

## Transactions

### Single-domain (internal)

When a repository method needs a transaction for its own tables, it manages the transaction entirely internally — begins, executes, commits or rolls back, and returns domain types. The domain has no knowledge a transaction occurred.

### Cross-domain

When an operation in domain A must commit atomically with an operation in domain B:

1. Domain A defines a callback interface for what it needs from B:
   ```go
   // person/repository.go
   type EnterpriseLinker interface {
       LinkToEnterprise(ctx context.Context, personID int64) error
   }
   ```

2. Domain B (or its db adapter) implements that interface.

3. Domain A's `Repository` method accepts the linker:
   ```go
   CreatePerson(ctx context.Context, p Person, linker EnterpriseLinker) (Person, error)
   ```

4. The `db/pg` adapter for domain A starts a transaction, threads it into the context via `withTx`, then calls the linker. The `db/pg` adapter for domain B retrieves the transaction via `conn(ctx, pool)` and uses it automatically.

```go
// db/pg/person.go
func (r *personRepository) CreatePerson(ctx context.Context, p person.Person, linker person.EnterpriseLinker) (person.Person, error) {
    tx, err := r.pool.Begin(ctx)
    if err != nil { ... }
    defer tx.Rollback(ctx)

    txCtx := withTx(ctx, tx)          // thread the transaction
    // ... insert person using conn(txCtx, r.pool) ...
    if err := linker.LinkToEnterprise(txCtx, created.ID); err != nil { ... }

    return created, tx.Commit(ctx)
}
```

### Known risks of the tx-in-context pattern

**1. Silent data inconsistency (most dangerous)**
If `withTx` is omitted before calling the cross-domain callback, the callback runs on its own pool connection outside the transaction. The operation succeeds, the data is inconsistent, and there is no error or panic.

_Mitigation:_ always call `withTx` immediately before passing `ctx` to the callback. Never pass `ctx` to the callback before enriching it with the transaction.

**2. Invisible dependency**
A transaction in flight is not visible in any function signature. Transactional behaviour can only be understood by tracing the context.

_Mitigation:_ document the transaction scope in the method that begins it.

**3. Context misuse**
The Go standard library recommends context only for request-scoped values that cross API or process boundaries, not for mutable state like a database transaction. This pattern violates that guidance.

_Mitigation:_ never pass a tx-carrying context across goroutine boundaries; keep transaction scope as narrow as possible. The alternative — threading `pgx.Tx` as an explicit parameter — would require db adapters to import each other, which is worse.

## Adding a new domain

1. Create the domain interface in `<domain>/repository.go`.
2. Add sentinel errors to `<domain>/errors.go`.
3. Write SQL files in `db/pg/queries/<domain>/` (and `db/mssql/queries/<domain>/` if the domain uses MSSQL).
4. Implement the repository in `db/pg/<domain>.go` (and `db/mssql/<domain>.go`).
5. Add the new repository field to `PGRepositories` (and `MSSQLRepositories`) in `db/factory.go`.
6. Wire the new repository into the relevant domain service in `cmd/server/main.go`.

## Adding a new database engine

1. Create `db/<engine>/` mirroring the `db/pg/` layout.
2. Add a `tx.go` with a `conn` helper following the same contract as `db/pg/tx.go`.
3. Add an `errors.go` with a `mapErr` function mapping engine error codes to domain sentinels.
4. Add a `New<Engine>Repositories` factory function in `db/factory.go`.
5. Register the new driver string in `config.ValidatePersistenceEnv`.
6. Add migration targets to the `Makefile`.

## Testing

| Scope                     | Adapter         | How to run                                                                                             |
|---------------------------|-----------------|--------------------------------------------------------------------------------------------------------|
| Unit / service tests      | in-memory       | `go test ./...` — always runs                                                                          |
| pg adapter tests          | real PostgreSQL | `DATABASE_URL=postgres://... go test ./db/pg/...` — skipped if `DATABASE_URL` is unset                 |
| mssql adapter tests       | real MSSQL      | `MSSQL_DATABASE_URL=sqlserver://... go test ./db/mssql/...` — skipped if `MSSQL_DATABASE_URL` is unset |
| cache decorator tests     | `inmem.Cache`   | `go test ./db/valkey/...` — always runs; no live Valkey needed                                         |
| cache adapter integration | real Valkey     | `CACHE_URL=valkey://... go test ./db/valkey/...` — skipped if `CACHE_URL` is unset                     |

## Local development

| Tool                    | Purpose                                                                       |
|-------------------------|-------------------------------------------------------------------------------|
| `devenv up`             | Starts PostgreSQL and other open-source services via Nix (no Docker required) |
| `docker compose up -d`  | Starts MSSQL and any other services not available in nixpkgs                  |
| `make migrate-up`       | Runs PostgreSQL migrations against `$DATABASE_URL`                            |
| `make migrate-mssql-up` | Runs MSSQL migrations against `$MSSQL_DATABASE_URL`                           |
