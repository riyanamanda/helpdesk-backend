#!/bin/sh
set -e

RETRIES=30
until nc -z "$DB_HOST" "$DB_PORT" 2>/dev/null || [ "$RETRIES" -eq 0 ]; do
  echo "Waiting for database... ($RETRIES retries left)"
  RETRIES=$((RETRIES - 1))
  sleep 2
done

if [ "$RETRIES" -eq 0 ]; then
  echo "ERROR: Could not connect to database after 60 seconds"
  exit 1
fi

echo "Running database migrations..."
goose -dir ./migrations up

if [ "${RUN_SEED:-true}" = "true" ]; then
  echo "Running database seed..."
  ./seed
fi

echo "Starting API server..."
exec ./api
