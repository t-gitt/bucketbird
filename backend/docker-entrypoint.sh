#!/bin/sh
set -e

# Build DSN from individual components if BB_DB_DSN is not set
if [ -z "$BB_DB_DSN" ]; then
  BB_DB_HOST="${BB_DB_HOST:-localhost}"
  BB_DB_PORT="${BB_DB_PORT:-5432}"
  BB_DB_USER="${BB_DB_USER:-postgres}"
  BB_DB_PASSWORD="${BB_DB_PASSWORD:-postgres}"
  BB_DB_NAME="${BB_DB_NAME:-postgres}"
  BB_DB_SSLMODE="${BB_DB_SSLMODE:-disable}"

  export BB_DB_DSN="postgres://${BB_DB_USER}:${BB_DB_PASSWORD}@${BB_DB_HOST}:${BB_DB_PORT}/${BB_DB_NAME}?sslmode=${BB_DB_SSLMODE}"
fi

# Default command is to start the API
CMD="${1:-api}"

case "$CMD" in
  migrate)
    shift
    echo "Running migrations..."
    /app/bucketbird migrate "$@"
    ;;
  api|serve)
    echo "Starting API server..."
    exec /app/bucketbird serve
    ;;
  *)
    # Allow running arbitrary commands
    exec "$@"
    ;;
esac
