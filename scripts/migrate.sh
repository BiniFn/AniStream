#!/bin/sh
set -eu

if ! command -v migrate >/dev/null 2>&1; then
    echo "migrate is not installed. Please install migrate and try again." >&2
    exit 1
fi

# Prompt for migration name
printf "Migration name: "
read -r MIGRATION_NAME

if [ -z "$MIGRATION_NAME" ]; then
    echo "Migration name cannot be empty" >&2
    exit 1
fi

migrate create -seq -dir migrations -ext sql "$MIGRATION_NAME"
echo "Migration files created for: $MIGRATION_NAME"