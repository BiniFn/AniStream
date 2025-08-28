MIGRATION_NAME := $(filter-out migrate,$(MAKECMDGOALS))

# ----- Migration ----- #
.PHONY: migrate
migrate: ## Generate migration files
	@if [ -z "$(MIGRATION_NAME)" ]; then \
	  echo "Usage: make migrate <migration_name>"; \
	  exit 1; \
	fi
	migrate create -seq -dir migrations -ext sql $(MIGRATION_NAME)

%:
	@:

# ----- Local Dev ----- #
.PHONY: dev-api dev-proxy dev-worker
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
.PHONY: sqlc genqlient
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
.PHONY: openapi
openapi: ## Generate openapi docs
	swag init -g ./cmd/api/main.go -o ./docs/swagger
	swagger2openapi --patch ./docs/swagger/swagger.yaml -o ./docs/openapi.yaml
	rm -rf ./docs/swagger

# ----- Go Build & Tidy ----- #
.PHONY: build tidy
build: ## Build API and Proxy binaries
	go build -o bin/api ./cmd/api
	go build -o bin/proxy ./cmd/proxy

tidy: ## Go mod tidy
	go mod tidy

# ----- Docker Build ----- #
.PHONY: docker-build-api docker-build-proxy docker-build-worker
docker-build-api: ## Build API Docker image
	docker build --target=api -t aniways-api -f docker/Dockerfile .

docker-build-proxy: ## Build Proxy Docker image
	docker build --target=proxy -t aniways-proxy -f docker/Dockerfile .

docker-build-worker: ## Build Worker Docker image
	docker build --target=worker -t aniways-worker -f docker/Dockerfile .

# ----- Docker Compose (Dev) ----- #
.PHONY: dev-docker-up dev-docker-down dev-docker-logs
dev-docker-up: ## Start dev containers
	docker compose -p aniways -f docker/docker-compose.dev.yaml --env-file .env.local up -d

dev-docker-down: ## Stop dev containers
	docker compose -p aniways -f docker/docker-compose.dev.yaml --env-file .env.local down

dev-docker-logs: ## View logs for all containers
	docker compose -p aniways -f docker/docker-compose.dev.yaml --env-file .env.local logs -f

# ----- Tmux Commands ----- #
.PHONY: tmux
tmux:
	./scripts/aniways-tmux.sh

# ----- Help Menu ----- #
.PHONY: help
help: ## Show help
	@grep -E '^[a-zA-Z_\-]+:.*?## ' $(MAKEFILE_LIST) | \
	  awk 'BEGIN {FS=":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
