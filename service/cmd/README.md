# cmd/

Each subdirectory under `cmd/` is an independent binary. They share domain packages and the repository abstraction but produce separate executables.

## Binaries

| Directory     | Type                 | Description                                                                |
|---------------|----------------------|----------------------------------------------------------------------------|
| `server/`     | Long-running service | HTTP API server — starts, listens, handles requests until signaled to stop |
| `task_alive/` | Short-lived ECS task | Example task binary — runs once and exits; copy this to add a new task     |

## Entry Point Pattern

Every binary follows the same signature:

```go
func run(ctx context.Context, args []string, getenv func(string) string, stdout, stderr io.Writer) error
```

`getenv` is injected instead of calling `os.Getenv` directly, which keeps configuration loading testable without real environment variables. `main.go` wires the process-level dependencies and calls `run`.

## Adding a New Task

1. Copy `cmd/task_alive/` to `cmd/task_<yourname>/`.
2. Rename the `"task"` log field to match your task name.
3. Read required EventBridge payload fields from `getenv` (e.g. `getenv("TASK_ID")`); hard-fail if they are missing.
4. Implement your logic in `run(...)` using the repository and `ctx`.
5. Add a Makefile target:
   ```makefile
   run-task-<yourname>: ## Run <yourname> task (inmemory driver); TASK_ID required
       PG_DRIVER=inmemory TASK_ID=$${TASK_ID:?TASK_ID is required} go run ./cmd/task<yourname>
   ```
6. Add the target to the `.PHONY` line alongside the other `run-task-*` entries.

See [README.md §Adding a new task](../README.md#adding-a-new-task) for deployment details (Dockerfile, Terraform, EventBridge Scheduler).

## Related

- [docs/architecture.md](../docs/architecture.md) — how `cmd/` fits into the hexagonal layout
