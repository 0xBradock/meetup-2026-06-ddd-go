# health/

The domain package. This is the centre of the hexagonal architecture: it defines what the application *does* without knowing *how* data is stored or *how* requests arrive.

**Nothing in this package imports from `db/` or `httpserver/`.** Adapters import from `health/` to satisfy its interfaces — never the reverse. This also means no web framework, no ORM, no driver dependency — just the Go standard library. The domain stays testable with `go test` alone, no infrastructure required.

See [docs/architecture.md](../docs/architecture.md) for the full explanation of Ports, Adapters, and the hexagonal rule.

## Files

| File            | DDD Role                | What it contains                                                                  |
|-----------------|-------------------------|-----------------------------------------------------------------------------------|
| `model.go`      | **Value Objects**       | `Health` and `Ping` — domain types with validation constructors                   |
| `service.go`    | **Application Service** | `Service` — orchestrates use cases (`GetHealth`, `Ping`)                          |
| `repository.go` | **Port**                | `Repository` interface — defines what persistence operations the domain needs     |
| `cache.go`      | **Port**                | `CachePort` interface — defines what cache operations the domain needs            |
| `errors.go`     | Sentinel errors         | `ErrNotFound`, `ErrConflict`, etc. — domain errors callers check with `errors.Is` |

### Value Objects (`model.go`)

`Health` and `Ping` carry domain meaning that plain primitives would lose. They have no identity — two `Health` values with the same `Status` are equal. `NewPing` is a guarded constructor that validates input and stamps the timestamp. Validation belongs here, not in the HTTP handler.

### Application Service (`service.go`)

`Service` coordinates use cases. It receives input, calls the `Repository` port, and returns output. It makes no decisions about SQL, HTTP, or caching — those are infrastructure concerns. A new use case is a new method on `Service`.

### Ports (`repository.go`, `cache.go`)

Ports are interfaces the domain defines to describe what it needs from the outside world. `Repository` says "I need to store and retrieve a Ping" without specifying how. `CachePort` says "I need to read and write cache entries" without specifying Valkey or Redis. The implementations live in `db/` — this package never imports them.

## Adding a New Use Case

1. Add a method to `Service` in `service.go` — keep it focused on one business operation.
2. If the use case needs persisted data, add the corresponding method to `Repository` in `repository.go`.
3. Implement the new method in each adapter under `db/` that needs to support it (`db/pg/`, `db/mssql/`, `db/inmem/`).
4. Write a unit test in `service_test.go` — pass an `inmem` repository or a hand-written stub. No database needed.

## Related

- [docs/architecture.md](../docs/architecture.md) — hexagonal diagram, Port/Adapter/Decorator definitions, DDD concepts
- [db/README.md](../db/README.md) — how the Repository port is implemented in each adapter
