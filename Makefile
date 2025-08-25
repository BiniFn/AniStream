MIGRATION_NAME := $(filter-out migrate,$(MAKECMDGOALS))

.PHONY: migrate dev-api dev-proxy dev-worker sqlc genqlient openapi build tidy docker-build-api docker-build-proxy docker-build-worker dev-docker-up dev-docker-down dev-docker-logs tmux help

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
dev-api: ## Run API with Air
	air -c .air.toml \
		-build.cmd "go build -o ./tmp/api ./cmd/api" \
		-build.bin "./tmp/api"

dev-proxy: ## Run proxy locally
	air -c .air.toml \
		-build.cmd "go build -o ./tmp/proxy ./cmd/proxy" \
		-build.full_bin "APP_ENV=development ./tmp/proxy"

dev-worker: ## Run worker locally
	air -c .air.toml \
		-build.cmd "go build -o ./tmp/worker ./cmd/worker" \
		-build.bin "./tmp/worker"

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

# ----- OpenAPI ----- #
openapi: ## Generate openapi docs
	swag init -g ./cmd/api/main.go -o ./docs/swagger
	swagger2openapi --patch ./docs/swagger/swagger.yaml -o ./docs/openapi.yaml
	rm -rf ./docs/swagger

# ----- Go Build & Tidy ----- #
build: ## Build API and Proxy binaries
	go build -o bin/api ./cmd/api
	go build -o bin/proxy ./cmd/proxy

tidy: ## Go mod tidy
	go mod tidy

# ----- Docker Build ----- #
docker-build-api: ## Build API Docker image
	docker build --target=api -t aniways-api -f docker/Dockerfile .

docker-build-proxy: ## Build Proxy Docker image
	docker build --target=proxy -t aniways-proxy -f docker/Dockerfile .

docker-build-worker: ## Build Worker Docker image
	docker build --target=worker -t aniways-worker -f docker/Dockerfile .

# ----- Docker Compose (Dev) ----- #
dev-docker-up: ## Start dev containers
	docker compose -p aniways -f docker/docker-compose.dev.yaml --env-file .env.local up -d

dev-docker-down: ## Stop dev containers
	docker compose -p aniways -f docker/docker-compose.dev.yaml --env-file .env.local down

dev-docker-logs: ## View logs for all containers
	docker compose -p aniways -f docker/docker-compose.dev.yaml --env-file .env.local logs -f

# ----- Tmux Commands ----- #
tmux:
	./scripts/aniways-tmux.sh

# ----- Help Menu ----- #
help: ## Show help
	@grep -E '^[a-zA-Z_\-]+:.*?## ' $(MAKEFILE_LIST) | \
	  awk 'BEGIN {FS=":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
