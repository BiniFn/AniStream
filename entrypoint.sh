#!/bin/sh
set -e

# Load environment variables if ENV_FILE_PATH is provided
if [ -n "$ENV_FILE_PATH" ] && [ -f "$ENV_FILE_PATH" ]; then
  echo "Loading env from $ENV_FILE_PATH"
  set -a
  . "$ENV_FILE_PATH"
  set +a
fi

echo "Waiting for Postgres..."
until nc -z db 5432; do sleep 1; done

echo "Waiting for Redis..."
until nc -z redis 6379; do sleep 1; done

echo "Starting API..."
exec "$@"
