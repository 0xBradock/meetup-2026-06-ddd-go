<p align="center">

<h1>golang server</h1>

![hero image](./docs/imgs/hero_gopher.png)

> Template Go API server with DDD layering, HTTP transport, and pluggable persistence (`inmemory` or `postgres`).

</p>

> [!NOTE]
> Read the [onboarding](./ONBOARDING.md) guide

## Quick Reference

| Task                 | Command                                                 |
| -------------------- | ------------------------------------------------------- |
| Run server (local)   | `make run`                                              |
| Unit tests           | `make test`                                             |
| Integration tests    | `make test-integration`                                 |
| Full quality check   | `make qc`                                               |
| Lint                 | `make lint`                                             |
| Run PG migrations    | `make migrate-up` (requires `DATABASE_URL`)             |
| Run MSSQL migrations | `make migrate-mssql-up` (requires `MSSQL_DATABASE_URL`) |

## Documentation

- **[Architecture Overview](docs/architecture.md)** - Hexagonal diagram and layer rules
- **[Architecture Decision Records](docs/adr/)** - Significant technical decisions
