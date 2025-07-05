# Makefile at repo root
.PHONY: help dev-api dev-cron dev-stream build test tidy docker

# ----- Local dev ----- #
dev-api:            ## Hot-reload API with Air
	air -c .air.toml

dev-cron:           ## Run cron binary
	go run ./cmd/cron

dev-stream:         ## Run streaming proxy
	go run ./cmd/streamer

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

help:               ## Show this help
	@grep -E '^[a-zA-Z_\-]+:.*?## ' $(MAKEFILE_LIST) | \
	  awk 'BEGIN {FS=":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'
