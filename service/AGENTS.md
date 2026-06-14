# Repository Guidelines

## Project Structure & Module Organization
`cmd` contains the executable entrypoint and wiring (`main.go`, `run.go`). Business code is intentionally flatter: `health` for domain and service logic, `db` for persistence factory/backends, `config` for runtime config, and `httpserver` for HTTP handlers and routes. Database assets live in `db/pg/migrations` and `db/mssql/migrations`.

## Build, Test, and Development Commands
Use `make` targets as the default workflow:

- `make run` starts the server with the in-memory repository on `:8080`.
- `make test` runs `go test ./...` across all packages.
- `make lint` runs `golangci-lint`.
- `make fmt-check` fails if any tracked Go file needs `gofmt`.
- `make qc` runs format, lint, tests, and commit-message checks together.

For Postgres-backed runs, set `PG_DRIVER=postgres` and `DATABASE_URL=...`. For MSSQL-backed runs, set `MSSQL_DRIVER=mssql` and `MSSQL_DATABASE_URL=...`. Both can be active simultaneously.

## Coding Style & Naming Conventions
Follow standard Go formatting with `gofmt`; use tabs and idiomatic package layout. Keep package names short and lowercase (`health`, `config`). Exported identifiers use PascalCase; unexported helpers use camelCase. Test files follow Go defaults with `_test.go`, and transport packages use explicit aliases like `httptransport` when needed for clarity.

## Testing Guidelines
Tests use Go’s `testing` package. Prefer table-driven tests where behavior branches, and call `t.Parallel()` for independent cases, matching the existing suite. Keep unit tests beside the code they verify, for example `health/service_test.go`. Run `make test` before opening a PR; run `make qc` before merging.

## Commit & Pull Request Guidelines
Commit messages follow Conventional Commits via Commitizen, for example `feat(template): bootstrap DDD Go API server` or `chore(devx): add pre-commit quality gates`. Git hooks are managed by devenv; enter `devenv shell` to enable them. Pull requests should include a short description, note any config or migration changes, link related issues, and include sample requests or screenshots only when API behavior or operator workflow changes.

## Configuration
Default runtime settings are `PG_DRIVER=inmemory`, `MSSQL_DRIVER=inmemory`, `HTTP_ADDR=:8080`. `DATABASE_URL` is required only when `PG_DRIVER=postgres`; `MSSQL_DATABASE_URL` is required only when `MSSQL_DRIVER=mssql`.

## Terraform Strategy
Use the single-root OpenTofu strategy documented in `terraform/README.md`.

- Keep one root under `terraform/`.
- Use `terraform/envs/<env>/backend.hcl` plus `terraform/envs/<env>/<env>.tfvars` for environment isolation.
- Do not introduce separate Terraform root modules per environment; keep all `.tf` files in `terraform/`.

## Template Onboarding Gate
This repository starts as a template and may require project-specific renaming/setup.

- Before any implementation work, first walk the user through `README.md` section `## Template Onboarding Checklist`.
- Confirm with the user whether each checklist item is already completed.
- If the user confirms the checklist is fully completed, remove this `## Template Onboarding Gate` section from `AGENTS.md` and remove `## Template Onboarding Checklist` from `README.md`.
