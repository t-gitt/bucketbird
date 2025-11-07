#!/bin/sh
set -e

TEMPLATE=/usr/share/nginx/html/config.js.template
TARGET=/usr/share/nginx/html/config.js

if [ -f "$TEMPLATE" ]; then
  # Allow BB_API_URL to default to empty string so fallback kicks in.
  envsubst '$$BB_API_URL $$BB_ALLOW_REGISTRATION $$BB_ENABLE_DEMO_LOGIN' < "$TEMPLATE" > "$TARGET"
fi

exec "$@"
