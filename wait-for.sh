#!/bin/sh
set -e

echo "Waiting for Postgres..."
until nc -z db 5432; do
  sleep 1
done

echo "Waiting for Redis..."
until nc -z redis 6379; do
  sleep 1
done

echo "Dependencies are ready. Starting app..."
exec "$@"
