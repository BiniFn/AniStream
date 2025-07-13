# Makefile at repo root

# capture everything after "migrate" as the migration name
MIGRATION_NAME := $(filter-out migrate,$(MAKECMDGOALS))

.PHONY: help dev-api dev-cron dev-stream build test tidy docker migrate

# ----- Make migration files ----- #
migrate:           ## Generate migration files
	@if [ -z "$(MIGRATION_NAME)" ]; then \
	  echo "Usage: make migrate <migration_name>"; \
	  exit 1; \
	fi
	migrate create -seq -dir migrations -ext sql $(MIGRATION_NAME)

# prevent "make" from trying to build MIGRATION_NAME as its own target
%:
	@:

# ----- Local dev ----- #
dev-api:            ## Hot-reload API with Air
	air -c .air.toml

dev-proxy:
	go run ./cmd/proxy

# ----- SQLC ----- #
sqlc:               ## Generate SQLC code
	sqlc generate

# ----- Build & test ----- #
build:              ## Build all three binaries into ./bin
	go build -o bin/api      ./cmd/api
	go build -o bin/cron     ./cmd/cron
	go build -o bin/streamer ./cmd/streamer

test:               ## Run race-enabled test suite
	go test ./... -race -cover

tidy:               ## Go mod tidy
	go mod tidy

# ----- Docker ----- #
docker:             ## Build multi-stage Docker image
	docker build -t aniways:latest .

docker-compose:
	docker-compose --env-file .env.local up -d

# ----- GraphQL ----- #
genqlient:          ## Generate GraphQL client code
	@if [ ! -f schema.graphql ]; then \
		npx --yes graphqurl https://graphql.anilist.co --introspect > schema.graphql; \
		if [ $$? -ne 0 ]; then \
			echo "Failed to fetch schema.graphql"; \
			exit 1; \
		fi \
	fi
	genqlient genqlient.yaml

help:               ## Show this help
	@grep -E '^[a-zA-Z_\-]+:.*?## ' $(MAKEFILE_LIST) | \
	  awk 'BEGIN {FS=":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'
