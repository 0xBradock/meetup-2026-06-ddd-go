# Make directives
MAKEFLAGS += --no-print-directory --warn-undefined-variables --no-builtin-rules


# Print help with descriptions
# Each recipe should be followed by `## <description>`
help: ## ⭐Show this help message
	@echo "Available commands (alphabetical):"
	@grep -Eh '^[[:alnum:]_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | \
		sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'

.PHONY: help

# Golang targets
run: ## ⭐Start server (inmemory driver)
	PG_DRIVER=inmemory MSSQL_DRIVER=inmemory HTTP_ADDR=:8080 go run ./cmd/server
