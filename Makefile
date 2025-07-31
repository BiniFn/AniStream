MIGRATION_NAME := $(filter-out migrate,$(MAKECMDGOALS))

.PHONY: help migrate dev-api dev-proxy sqlc build tidy \
        docker-build-api docker-build-proxy \
        dev-up dev-down dev-logs \
        genqlient

# ----- Migration ----- #
migrate: ## Generate migration files
	@if [ -z "$(MIGRATION_NAME)" ]; then \
	  echo "Usage: make migrate <migration_name>"; \
	  exit 1; \
	fi
	migrate create -seq -dir migrations -ext sql $(MIGRATION_NAME)

%:
	@:

# ----- Local Dev ----- #
dev-api: ## Hot-reload API with Air
	air -c .air.toml

dev-proxy: ## Run proxy locally
	go run ./cmd/proxy

# ----- SQLC & GraphQL ----- #
sqlc: ## Generate SQLC code
	sqlc generate

genqlient: ## Generate GraphQL client code
	@if [ ! -f schema.graphql ]; then \
		npx --yes graphqurl https://graphql.anilist.co --introspect > schema.graphql; \
		if [ $$? -ne 0 ]; then \
			echo "Failed to fetch schema.graphql"; \
			exit 1; \
		fi \
	fi
	genqlient genqlient.yaml

# ----- Go Build & Tidy ----- #
build: ## Build API and Proxy binaries
	go build -o bin/api ./cmd/api
	go build -o bin/proxy ./cmd/proxy

tidy: ## Go mod tidy
	go mod tidy

# ----- Docker Build ----- #
docker-build-api: ## Build API Docker image
	docker build -t aniways-api -f infra/api/Dockerfile .

docker-build-proxy: ## Build Proxy Docker image
	docker build -t aniways-proxy -f infra/proxy/Dockerfile .

# ----- Docker Compose (Dev) ----- #
dev-up: ## Start dev containers
	docker compose -p aniways -f infra/docker-compose.dev.yaml --env-file .env.local up -d

dev-down: ## Stop dev containers
	docker compose -p aniways -f infra/docker-compose.dev.yaml --env-file .env.local down

dev-logs: ## View logs for all containers
	docker compose -p aniways -f infra/docker-compose.dev.yaml --env-file .env.local logs -f

# ----- Help Menu ----- #
help: ## Show help
	@grep -E '^[a-zA-Z_\-]+:.*?## ' $(MAKEFILE_LIST) | \
	  awk 'BEGIN {FS=":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
