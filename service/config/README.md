# config/

Runtime configuration loading and validation. All environment variable access is centralised here so that the rest of the codebase never calls `os.Getenv` directly — it receives a `getenv func(string) string` from the binary entrypoint, which makes configuration testable without real environment variables.

## Environment Variables

### Server (`LoadServerConfig`)

| Variable               | Default    | Description                                          |
|------------------------|------------|------------------------------------------------------|
| `PG_DRIVER`            | `inmemory` | PostgreSQL driver: `inmemory` or `postgres`          |
| `MSSQL_DRIVER`         | `inmemory` | MSSQL driver: `inmemory` or `mssql`                  |
| `ENVIRONMENT`          | `recette`  | Deployment environment: `prod`, `preprod`, `recette` |
| `DATABASE_URL`         | —          | Required when `PG_DRIVER=postgres`                   |
| `MSSQL_DATABASE_URL`   | —          | Required when `MSSQL_DRIVER=mssql`                   |
| `HTTP_ADDR`            | `:8080`    | TCP address the HTTP server listens on               |
| `HTTP_READ_TIMEOUT`    | `5s`       | Maximum duration for reading the full request        |
| `HTTP_WRITE_TIMEOUT`   | `5s`       | Maximum duration for writing the response            |
| `HTTP_IDLE_TIMEOUT`    | `5s`       | Maximum idle time for keep-alive connections         |
| `SERVER_SHUTDOWN_TIME` | `5`        | Graceful shutdown window in seconds                  |

### Tasks (`LoadTaskConfig`)

| Variable       | Default    | Description                                            |
|----------------|------------|--------------------------------------------------------|
| `PG_DRIVER`    | `inmemory` | Same as server                                         |
| `ENVIRONMENT`  | `recette`  | Same as server                                         |
| `DATABASE_URL` | —          | Same as server                                         |
| `TASK_ID`      | —          | Required; injected by EventBridge Scheduler at runtime |

## Adding a New Config Field

1. Add the field to the relevant struct in `config.go`.
2. Read it inside the corresponding `Load*Config` function using `getenv`.
3. Apply validation (fail-fast with a descriptive error for required fields).
4. Update the table above.

## Related

- [docs/architecture.md](../docs/architecture.md) — where `config/` sits in the overall layout
