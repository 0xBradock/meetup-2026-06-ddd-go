<p align="center">

<h1>golang server</h1>

![hero image](./docs/imgs/hero_gopher.png)

> Template Go API server with DDD layering, HTTP transport, and pluggable persistence (`inmemory` or `postgres`).
</p>

> [!NOTE]
> Read the [onboarding](./ONBOARDING.md) guide

## Quick Reference

| Task                 | Command                                                 |
|----------------------|---------------------------------------------------------|
| Run server (local)   | `make run`                                              |
| Unit tests           | `make test`                                             |
| Integration tests    | `make test-integration`                                 |
| Full quality check   | `make qc`                                               |
| Lint                 | `make lint`                                             |
| Run PG migrations    | `make migrate-up` (requires `DATABASE_URL`)             |
| Run MSSQL migrations | `make migrate-mssql-up` (requires `MSSQL_DATABASE_URL`) |

## Documentation

- **[Onboading](./ONBOARDING.md)** - Guide to be used to onboard new teammates to new or existing projects
- **[Architecture Overview](docs/architecture.md)** - Hexagonal diagram and layer rules
- **[Architecture Decision Records](docs/adr/)** - Significant technical decisions
- **[Deployment Guides](docs/deployment/)** - Step-by-step deployment procedures
- **[Terraform Documentation](terraform/README.md)** - Infrastructure deployment
- **[Contributing Guide](CONTRIBUTING.md)** - Branch strategy, commit messages, MR checklist
- **[TODO List](docs/TODO.md)** - Future improvements and open tasks

## Setup Project

1. Setup enviroment and tooling by running the commands in [devenv - first time setup](./docs/devenv.md#first-time-setup).
2. Update project variables in [.project.yml](./.project.yml)
3. Run `./scripts/sync-config.sh` which create all configurations for the project with the pinned versions:

- `.mise.toml` for mise/devenv tooling;
- `.versions.env` for shell and OpenTofu variables;
- `.project.env` for shell and Make defaults;
- `.gitlab/versions.yml` for GitLab CI parse-time variables.
- `.gitlab/project.yml` for GitLab CI project variables.

4. Run `make run` to run the http server locally.
5. Run `curl -v http://localhost:8080/health` 

## Template Onboarding Checklist

> [!NOTE]
> This repository is a template. Before your first real deployment, update service naming and defaults to match your project.
> After all items are completed and confirmed, this section can be removed.

- [ ] Complete the [devenv first-time setup](docs/devenv.md#first-time-setup).
- [ ] Set `project.service_name` in `.project.yml`, then run `./scripts/sync-config.sh`.
- [ ] Update `terraform/envs/<env>/backend.hcl` (state key) and fill `TODO` placeholders in `terraform/envs/<env>/<env>.tfvars`.
- [ ] Run `go test ./...` and `tofu fmt -check -recursive terraform` after renaming.
- [ ] Add `SONAR_TOKEN` in your Gitlab repository variables.

## Naming Contract

- `SERVICE_NAME` is the single naming input for this template.
- The committed source for that value is `.project.yml`.
- Runtime logs use `SERVICE_NAME` via `config.LoadServerConfig`.
- Build and deploy defaults derive image repository and ECS naming from `SERVICE_NAME` and `TF_ENV`.
- Terraform examples use `SERVICE_NAME` placeholders in state key and image references.

## Architecture

```text
.
├── cmd/
│   ├── server/           # HTTP server binary (main entrypoint).
│   └── task_alive/   # Example task binary — copy to add a new task.
├── config/       # Runtime configuration parsing and validation.
├── db/           # Persistence factory, repositories, and migrations.
├── health/       # Health domain model, validation rules, use-cases, and repository port.
├── httpserver/   # HTTP routes and handlers.
├── docs/         # ADRs, deployment notes, devenv guide, and project documentation.
├── scripts/      # Project automation and generated-config sync scripts.
├── terraform/    # OpenTofu infrastructure configuration.
├── .project.yml  # Project identity and tool version source of truth.
├── devenv.nix    # Development shell, scripts, and Git hook definitions.
├── devenv.yaml   # devenv input definitions.
└── Makefile      # Thin workflow wrappers for local development and deployment.
```

The application follows a small hexagonal layout: `health` owns domain behavior, `db` provides persistence adapters, `httpserver` handles transport, and `cmd` composes the runtime dependencies. See [docs/architecture.md](docs/architecture.md) for the full diagram and layer rules.

Each subdirectory under `cmd/` is an independent binary. The HTTP server (`cmd/server`) is the long-running service.
Task binaries (`cmd/task*`) are short-lived ECS tasks or Fargate spot jobs that share the same domain packages and repository abstraction.

Each folder has its own README: [cmd/](cmd/README.md) · [config/](config/README.md) · [health/](health/README.md) · [httpserver/](httpserver/README.md) · [db/](db/README.md) · [scripts/](scripts/README.md) · [docs/](docs/README.md)

## Persistence Factory

`db/factory.go` exposes two factory functions, each returning a typed struct of domain repositories:

- `NewPG(ctx, driver, getenv)` — PostgreSQL-backed repos when `driver="postgres"`, in-memory otherwise.
- `NewMSSQL(ctx, driver, getenv)` — MSSQL-backed repos when `driver="mssql"`, in-memory otherwise.

| Driver     | Factory    | Required env var     |
|------------|------------|----------------------|
| `inmemory` | either     | —                    |
| `postgres` | `NewPG`    | `DATABASE_URL`       |
| `mssql`    | `NewMSSQL` | `MSSQL_DATABASE_URL` |

All SQL is hand-written — no ORM, no query builder. PostgreSQL uses raw `pgx/v5`; MSSQL uses `database/sql` with `go-mssqldb`.

### Cache layer

An optional Valkey (Redis-compatible) cache can wrap any repository.
When `CACHE_URL` is set, the server wraps the active repository in a cache decorator before passing it to the domain service.
Omitting `CACHE_URL` skips the cache entirely. On cache errors the decorator falls through to the underlying repository transparently.

See [db/README.md](db/README.md) for the full design, transaction strategy, and testing guide.

## Development

```bash
make run                              # Run HTTP server with in-memory repository on :8080
TASK_ID=test-id make run-task-alive   # Run the alive task (exits when done)

make test       # Run tests
make qc         # Run all quality checks (local and CI/CD-ready)

# Run commit quality checks (version range)
CZ_CHECK_REV_RANGE=origin/main..HEAD make qc
```

Each binary follows the same `run(ctx, args, getenv, stdout, stderr)` signature — `getenv` is injected so configuration loading is testable without real environment variables.

Environment variables (defaults in parentheses):

- `PG_DRIVER` (`inmemory`; use `postgres` to enable PostgreSQL)
- `MSSQL_DRIVER` (`inmemory`; use `mssql` to enable Microsoft SQL Server)
- `ENVIRONMENT` (`recette`; allowed: `prod`, `preprod`, `recette`)
- `DATABASE_URL` (required when `PG_DRIVER=postgres`)
- `MSSQL_DATABASE_URL` (required when `MSSQL_DRIVER=mssql`)
- `HTTP_ADDR` (`:8080`) — server only
- `HTTP_READ_TIMEOUT` (`5s`) — server only
- `HTTP_WRITE_TIMEOUT` (`5s`) — server only
- `HTTP_IDLE_TIMEOUT` (`5s`) — server only
- `SERVER_SHUTDOWN_TIME` (`5` seconds) — server only
- `TASK_ID` (required) — task binaries only; injected by EventBridge Scheduler at runtime

Sample requests:

```bash
curl -s http://localhost:8080/health
```

Run with Postgres:

```bash
PG_DRIVER=postgres DATABASE_URL=postgres://user:pass@localhost:5432/app?sslmode=disable go run ./cmd/server
```

## Adding a new task

1. Copy `cmd/task_alive/` to `cmd/task<yourname>/`.
2. Rename the `"task"` log field to match your task name.
3. Read any EventBridge payload fields from `getenv` (e.g. `getenv("TASK_ID")`); hard-fail if required fields are missing.
4. Implement your logic in `run(...)` using `pgRepos` and `ctx`.
5. Add a Makefile target:
   ```makefile
   run-task-<yourname>: ## Run <yourname> task (inmemory driver); TASK_ID required
       PG_DRIVER=inmemory TASK_ID=$${TASK_ID:?TASK_ID is required} go run ./cmd/task<yourname>
   ```
6. Add the target to the `.PHONY` line alongside the other `run-task-*` entries.

**Deploying with EventBridge Scheduler:**
- Add a `go build` line and `COPY` for the new binary in the Dockerfile (follow the `task_alive` pattern).
- Add an `aws_ecs_task_definition` in `terraform/ecs.tf` (follow the `task_alive` pattern — same image, override `entryPoint`).
- Add an `aws_scheduler_schedule` in `terraform/scheduler.tf` (follow the existing pattern). EventBridge Scheduler injects `TASK_ID` via `<aws.scheduler.execution-id>` in the container override input.
- See [ADR-011](docs/adr/011-eventbridge-scheduler-task-trigger.md) for the full rationale.

## AWS Fargate (OpenTofu)

OpenTofu infrastructure is under `terraform/`.

- Strategy: single root with per-environment backend and var files (see `terraform/README.md`)
- Backend configs: `terraform/envs/<env>/backend.hcl`
- Local var files: `terraform/envs/<env>/<env>.tfvars` (tracked; fill in `TODO` placeholders before first apply)
- Resource definitions: root `terraform/*.tf` files (split by concern)
- HTTP is split into public and private exposures (HTTPS follow-up noted)
- ECS tasks, RDS PostgreSQL, and RDS Proxy are private
- App receives `DATABASE_URL` via Secrets Manager and connects through RDS Proxy
- S3 backend per environment (no locking yet; follow-up noted)
- ECR repository per environment (`identity-prod`, `identity-preprod`)

### Infrastructure Deployment

Use the provided Makefile targets for convenient deployment:

```bash
# Initialize OpenTofu for preprod (default)
make tf-init

# Initialize for a specific environment
make tf-init ENV=prod

# Plan infrastructure changes
make tf-plan ENV=preprod

# Apply infrastructure changes
make tf-apply ENV=preprod

# View outputs
make tf-output ENV=preprod

# Destroy infrastructure (requires confirmation)
make tf-destroy ENV=preprod
```

See [terraform/README.md](terraform/README.md) for detailed deployment workflows and troubleshooting.

### Prerequisites

Fill in the `TODO` placeholders in `terraform/envs/<env>/<env>.tfvars`:

- `public_subnets` — public subnet IDs (not CIDR blocks) for the public ALB
- `allowed_private_source_sg_ids` — SGs of internal callers
- `internal_alb_certificate_arn` — ACM certificate for the gRPC listener
- `app_image` — overridden at deploy time via `TF_VAR_app_image` from CI/CD

### Building and Pushing Docker Images

```bash
# Build and push to ECR
make image-push AWS_PROFILE=nonprod IMAGE_TAG=v1.0.0

# Print the image URI for Terraform
make print-image-uri IMAGE_TAG=v1.0.0
```
