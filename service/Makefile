# Make directives
MAKEFLAGS += --no-print-directory --warn-undefined-variables --no-builtin-rules

# Tooling and project metadata are managed from .project.yml.
-include .project.env
CZ_CHECK_REV_RANGE ?= HEAD~1..HEAD

# AWS
AWS_PROFILE ?= nonprod
AWS_REGION ?= eu-west-3
ECR_ACCOUNT_ID ?= 215627216890
SERVICE_NAME ?= golang-server

# Terraform configuration
ENV ?= preprod
TF_ENV ?= $(ENV)
TF_DIR := terraform
TF_ENV_DIR := $(TF_DIR)/envs/$(ENV)
TFVARS_FILE := $(TF_ENV_DIR)/$(ENV).tfvars
BACKEND_FILE := $(TF_ENV_DIR)/backend.hcl

# Service image
ECR_REPO ?= $(SERVICE_NAME)-$(TF_ENV)
IMAGE_TAG ?= $(shell git rev-parse --short=12 HEAD)
LOCAL_IMAGE ?= $(ECR_REPO):$(IMAGE_TAG)
IMAGE_URI := $(ECR_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com/$(ECR_REPO):$(IMAGE_TAG)

# Print help with descriptions
# Each recipe should be followed by `## <description>`
help: ## ⭐Show this help message
	@echo "Available commands (alphabetical):"
	@grep -Eh '^[[:alnum:]_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | \
		sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'

qc: ## ⭐Run all quality checks
	@echo "[qc] Running quality control checks..."
	@$(MAKE) fmt-check
	@$(MAKE) lint
	@$(MAKE) lint-openapi
	@$(MAKE) test
	@$(MAKE) cz-check
	@echo "[qc] All checks passed."

.PHONY: help qc

# DB targets
migrate-up: ## Run PostgreSQL migrations up
	migrate -path db/pg/migrations -database "$$DATABASE_URL" up

migrate-down: ## Rollback last PostgreSQL migration
	migrate -path db/pg/migrations -database "$$DATABASE_URL" down 1

migrate-mssql-up: ## Run MSSQL migrations up
	migrate -path db/mssql/migrations -database "$$MSSQL_DATABASE_URL" up

migrate-mssql-down: ## Rollback last MSSQL migration
	migrate -path db/mssql/migrations -database "$$MSSQL_DATABASE_URL" down 1

.PHONY: migrate-up migrate-down migrate-mssql-up migrate-mssql-down

# Golang targets
run: ## ⭐Start server (inmemory driver)
	PG_DRIVER=inmemory MSSQL_DRIVER=inmemory HTTP_ADDR=:8080 go run ./cmd/server

run-task-alive: ## Run alive task (inmemory driver); TASK_ID required
	PG_DRIVER=inmemory TASK_ID=$${TASK_ID:?TASK_ID is required} go run ./cmd/task_alive

fmt-check: ## Check Go code formatting
	@echo "[fmt-check] Running gofmt validation..."
	@files="$$(go list -f '{{range .GoFiles}}{{$$.Dir}}/{{.}} {{end}}{{range .TestGoFiles}}{{$$.Dir}}/{{.}} {{end}}' ./...)"; \
	if [ -z "$$files" ]; then \
		echo "[fmt-check] No Go files found."; \
		exit 0; \
	fi; \
	unformatted="$$(gofmt -l $$files)"; \
	if [ -n "$$unformatted" ]; then \
		echo "[fmt-check] Unformatted files:"; \
		printf '%s\n' "$$unformatted"; \
		exit 1; \
	fi; \
	echo "[fmt-check] Passed."

lint: ## Run golangci-lint
	@echo "[lint] Running golangci-lint..."
	@which golangci-lint
	@golangci-lint run
	@echo "[lint] Passed."

lint-openapi: ## Lint OpenAPI spec with vacuum
	@echo "[lint-openapi] Linting OpenAPI spec..."
	@vacuum lint -z openapi/openapi.yaml
	@echo "[lint-openapi] Passed."

test: ## Run all tests
	@echo "[test] Running go test..."
	@go test ./...
	@echo "[test] Passed."

test-integration: ## Run integration tests with coverage
	@echo "[test-integration] Running integration tests..."
	@mkdir -p coverage
	@go test -tags=integration -coverprofile=coverage/integration.out ./cmd/server/...
	@echo "[test-integration] Passed."

.PHONY: run run-task-alive fmt-check lint lint-openapi test test-integration

cz-check: ## Validate commit messages
	@if ! git rev-parse --verify HEAD >/dev/null 2>&1; then \
		echo "[cz-check] No commits found yet. Skipping commit message validation."; \
		exit 0; \
	fi
	@if git rev-parse --verify HEAD~1 >/dev/null 2>&1; then \
		echo "[cz-check] Validating commit messages in range $(CZ_CHECK_REV_RANGE)..."; \
		cz check --rev-range "$(CZ_CHECK_REV_RANGE)"; \
	else \
		echo "[cz-check] Only one commit found. Validating latest commit message..."; \
		cz check --message "$$(git log -1 --pretty=%s)"; \
	fi
	@echo "[cz-check] Passed."

cz-check-msg: ## Validate commit message file
	@if [ -z "$(COMMIT_MSG_FILE)" ]; then \
		echo "[cz-check-msg] COMMIT_MSG_FILE is required."; \
		echo "[cz-check-msg] Usage: make cz-check-msg COMMIT_MSG_FILE=.git/COMMIT_EDITMSG"; \
		exit 1; \
	fi
	@echo "[cz-check-msg] Validating commit message format..."
	@cz check --commit-msg-file "$(COMMIT_MSG_FILE)"
	@echo "[cz-check-msg] Passed."

.PHONY: cz-check cz-check-msg


docker-build: ## Build Docker image
	@echo "[docker-build] Building $(LOCAL_IMAGE)..."
	@docker build -t "$(LOCAL_IMAGE)" .
	@echo "[docker-build] Built $(LOCAL_IMAGE)."

ecr-login: ## Login to AWS ECR
	@echo "[ecr-login] Logging in to $(ECR_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com..."
	@aws ecr get-login-password --region "$(AWS_REGION)" --profile "$(AWS_PROFILE)" \
		| docker login --username AWS --password-stdin "$(ECR_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com"
	@echo "[ecr-login] Logged in."

ecr-tag: ## Tag image for ECR
	@echo "[ecr-tag] Tagging $(LOCAL_IMAGE) as $(IMAGE_URI)..."
	@docker tag "$(LOCAL_IMAGE)" "$(IMAGE_URI)"
	@echo "[ecr-tag] Tagged $(IMAGE_URI)."

ecr-push: ## Push image to ECR
	@echo "[ecr-push] Pushing $(IMAGE_URI)..."
	@docker push "$(IMAGE_URI)"
	@echo "[ecr-push] Pushed $(IMAGE_URI)."

image-push: ecr-login docker-build ecr-tag ecr-push ## Build and push image to ECR

print-image-uri: ## Print ECR image URI
	@printf '%s\n' "$(IMAGE_URI)"

.PHONY: docker-build ecr-login ecr-tag ecr-push image-push print-image-uri

# Terraform targets
tf-init: ## Initialize Terraform for an environment
	@echo "[tf-init] Initializing Terraform for environment: $(ENV)..."
	@test -f "$(BACKEND_FILE)" || (echo "[tf-init] Missing backend file: $(BACKEND_FILE)" && exit 1)
	@cd "$(TF_DIR)" && tofu init -upgrade -backend-config="envs/$(ENV)/backend.hcl"
	@echo "[tf-init] Initialized."

tf-validate: ## Validate Terraform configuration
	@echo "[tf-validate] Validating Terraform configuration..."
	@cd "$(TF_DIR)" && tofu validate
	@echo "[tf-validate] Validation passed."

tf-plan: ## Plan Terraform changes
	@echo "[tf-plan] Planning Terraform changes for environment: $(ENV)..."
	@test -f "$(TFVARS_FILE)" || (echo "[tf-plan] Missing vars file: $(TFVARS_FILE)." && exit 1)
	@cd "$(TF_DIR)" && tofu plan -var-file="envs/$(ENV)/$(ENV).tfvars" -var="app_image=$(IMAGE_URI)"
	@echo "[tf-plan] Plan complete."

tf-apply: ## Apply Terraform changes
	@echo "[tf-apply] Applying Terraform changes for environment: $(ENV)..."
	@test -f "$(TFVARS_FILE)" || (echo "[tf-apply] Missing vars file: $(TFVARS_FILE)." && exit 1)
	@cd "$(TF_DIR)" && tofu apply -var-file="envs/$(ENV)/$(ENV).tfvars" -var="app_image=$(IMAGE_URI)"
	@echo "[tf-apply] Apply complete."

tf-destroy: ## Destroy Terraform infrastructure
	@echo "========================================="
	@echo "WARNING: This will DESTROY all infrastructure in environment: $(ENV)"
	@echo "========================================="
	@printf "Type '%s' to confirm destruction: " "$(ENV)" && read -r confirm && [ "$$confirm" = "$(ENV)" ]
	@test -f "$(TFVARS_FILE)" || (echo "[tf-destroy] Missing vars file: $(TFVARS_FILE)." && exit 1)
	@echo "[tf-destroy] Destroying infrastructure for environment: $(ENV)..."
	@cd "$(TF_DIR)" && tofu destroy -var-file="envs/$(ENV)/$(ENV).tfvars"
	@echo "[tf-destroy] Destroy complete."

tf-output: ## Show Terraform outputs
	@echo "[tf-output] Fetching outputs for environment: $(ENV)..."
	@cd "$(TF_DIR)" && tofu output
	@echo "[tf-output] Output complete."

tf-fmt: ## Format Terraform files
	@echo "[tf-fmt] Formatting Terraform files..."
	@tofu fmt -recursive terraform/
	@echo "[tf-fmt] Formatting complete."

tf-fmt-check: ## Check Terraform formatting
	@echo "[tf-fmt-check] Checking Terraform formatting..."
	@tofu fmt -check -recursive terraform/
	@echo "[tf-fmt-check] Format check passed."

.PHONY: tf-init tf-validate tf-plan tf-apply tf-destroy tf-output tf-fmt tf-fmt-check
