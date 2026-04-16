#!/bin/sh
set -e

# Run migrations if DATABASE_URL is set
if [ -n "$DATABASE_URL" ]; then
    echo "Running database migrations..."
    migrate -database "$DATABASE_URL" -path /app/internal/db/migrations up || echo "Migration warning: $?"
fi

# Start the application
exec ./server
