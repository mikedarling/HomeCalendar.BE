#!/usr/bin/env sh
# Simple healthcheck script that probes /healthz over HTTP or HTTPS
# Uses TLS_CERT_FILE/TLS_KEY_FILE presence to decide TLS usage and TLS_PORT env var

set -e

# Default ports
HTTP_PORT=${HTTP_PORT:-8080}
TLS_PORT=${TLS_PORT:-8443}

if [ -n "$TLS_CERT_FILE" ] && [ -n "$TLS_KEY_FILE" ]; then
  URL="https://127.0.0.1:${TLS_PORT}/healthz"
  # Allow insecure for self-signed certs — in prod use valid certs
  curl --fail --silent --insecure "$URL" >/dev/null 2>&1
else
  URL="http://127.0.0.1:${HTTP_PORT}/healthz"
  curl --fail --silent "$URL" >/dev/null 2>&1
fi

exit 0
