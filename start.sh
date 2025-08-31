#!/bin/sh
set -e  # exit if any command fails

# Ensure required env vars are set
if [ -z "$GOOSE_DRIVER" ] || [ -z "$GOOSE_DBSTRING" ]; then
  echo "❌ Missing required environment variables: GOOSE_DRIVER and/or GOOSE_DBSTRING"
else
  # Run migrations
  echo "▶ Running database migrations..."
  goose -dir ./migrations up

fi

# Start the main service
echo "▶ Starting application..."
/app/main
