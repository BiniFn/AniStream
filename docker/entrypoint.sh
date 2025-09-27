#!/bin/sh
set -e

# Load environment variables if ENV_FILE_PATH is provided
if [ -n "$ENV_FILE_PATH" ] && [ -f "$ENV_FILE_PATH" ]; then
  echo "Loading env from $ENV_FILE_PATH"
  set -a
  . "$ENV_FILE_PATH"
  set +a
fi

DB_HOST=$(echo "$DATABASE_URL" | sed -E 's|^[^@]+@([^:/]+).*|\1|')
DB_PORT=$(echo "$DATABASE_URL" | sed -E 's|.*:([0-9]+)/.*|\1|')

REDIS_HOST=$(echo "$REDIS_ADDR" | cut -d: -f1)
REDIS_PORT=$(echo "$REDIS_ADDR" | cut -d: -f2)

# Fallbacks
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
REDIS_HOST=${REDIS_HOST:-localhost}
REDIS_PORT=${REDIS_PORT:-6379}

echo "Waiting for Postgres at $DB_HOST:$DB_PORT..."
until nc -z "$DB_HOST" "$DB_PORT"; do sleep 1; done

echo "Waiting for Redis at $REDIS_HOST:$REDIS_PORT..."
until nc -z "$REDIS_HOST" "$REDIS_PORT"; do sleep 1; done

echo "Starting application..."
exec "$@"
