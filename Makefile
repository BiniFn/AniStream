MIGRATION_NAME := $(filter-out migrate,$(MAKECMDGOALS))

.PHONY: migrate dev-all dev-api dev-proxy dev-worker sqlc genqlient build tidy docker-build-api docker-build-proxy docker-build-worker dev-up dev-down dev-logs help

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
dev-all: ## Run API, Proxy, and Worker together with tmux
	tmux new-session -s aniways -n dev \; \
		send-keys "make dev-api" C-m \; \
		split-window -h \; \
		send-keys "make dev-proxy" C-m \; \
		split-window -v \; \
		send-keys "make dev-worker" C-m \; \
		select-pane -t 0 \; \
		attach

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

# ----- Go Build & Tidy ----- #
build: ## Build API and Proxy binaries
	go build -o bin/api ./cmd/api
	go build -o bin/proxy ./cmd/proxy

tidy: ## Go mod tidy
	go mod tidy

# ----- Docker Build ----- #
docker-build-api: ## Build API Docker image
	docker build --target=api -t aniways-api -f infra/Dockerfile .

docker-build-proxy: ## Build Proxy Docker image
	docker build --target=proxy -t aniways-proxy -f infra/Dockerfile .

docker-build-worker: ## Build Worker Docker image
	docker build --target=worker -t aniways-worker -f infra/Dockerfile .

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
