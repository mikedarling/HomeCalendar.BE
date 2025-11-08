# Emcoded.Calendar.BE

## Overview

Backend for the Calendar Applications written in Go.

## Docker

This repo includes a `Dockerfile` and `docker-compose.yml` for running the service in a container.

## Quick Start

For development, use the provided script to build and run the container with TLS:

```bash
./scripts/restart-container.sh
```

This script will:
- Stop and remove any existing container
- Rebuild the Docker image
- Create self-signed certificates if needed
- Start the container with TLS on port 8443

Environment variables
- `APP_BASE_PATH` (optional) — path inside the container where static assets are served from. In the image this defaults to `/root`.

Examples
- Development (mount local `root` to see changes immediately):

```bash
docker run -d --name emcoded_calendar_be_dev -p 8080:8080 -v "$PWD/root:/root:ro" mikedarling/emcoded.calendar.be:latest
```

- Production (use baked-in assets from the image):

```bash
docker run -d --name emcoded_calendar_be -p 8080:8080 mikedarling/emcoded.calendar.be:latest
```

If you need to override the base path at runtime, pass `APP_BASE_PATH`:

```bash
docker run -d --name emcoded_calendar_be -p 8080:8080 -e APP_BASE_PATH=/root mikedarling/emcoded.calendar.be:latest
```

TLS / HTTPS
----------

To run the service over HTTPS, provide the TLS certificate and key files to the container and set the env vars `TLS_CERT_FILE` and `TLS_KEY_FILE`. By default the server will listen on port `8443` for TLS — override with `TLS_PORT`.

### Local Development with Self-Signed Certificates

For local development, create self-signed certificates:

```bash
# Create certificate directory
mkdir -p /tmp/certs

# Generate self-signed certificate
cd /tmp/certs
openssl req -x509 -newkey rsa:4096 -keyout privkey.pem -out fullchain.pem -days 365 -nodes -subj "/C=US/ST=Local/L=Local/O=Dev/CN=localhost"

# Run container with self-signed certificates
docker run -d --name emcoded_calendar_be_tls \
	-p 8443:8443 \
	-v /tmp/certs/fullchain.pem:/etc/ssl/certs/fullchain.pem:ro \
	-v /tmp/certs/privkey.pem:/etc/ssl/private/privkey.pem:ro \
	-e TLS_CERT_FILE=/etc/ssl/certs/fullchain.pem \
	-e TLS_KEY_FILE=/etc/ssl/private/privkey.pem \
	mikedarling/emcoded.calendar.be:latest
```

Access the service at: https://localhost:8443 (you'll need to accept the self-signed certificate warning)

### Local Development with mkcert (Recommended)

For a better local development experience with trusted certificates:

```bash
# Install mkcert (macOS)
brew install mkcert

# Create local CA
mkcert -install

# Generate certificates for localhost
mkdir -p ~/certs
cd ~/certs
mkcert localhost

# Run container with trusted certificates
docker run -d --name emcoded_calendar_be_tls \
	-p 8443:8443 \
	-v ~/certs/localhost.pem:/etc/ssl/certs/fullchain.pem:ro \
	-v ~/certs/localhost-key.pem:/etc/ssl/private/privkey.pem:ro \
	-e TLS_CERT_FILE=/etc/ssl/certs/fullchain.pem \
	-e TLS_KEY_FILE=/etc/ssl/private/privkey.pem \
	mikedarling/emcoded.calendar.be:latest
```

### Production Example

Example (production-style, mount certs read-only):

```bash
docker run -d --name emcoded_calendar_be_tls \
	-p 8443:8443 \
	-v /path/to/certs/fullchain.pem:/etc/ssl/certs/fullchain.pem:ro \
	-v /path/to/certs/privkey.pem:/etc/ssl/private/privkey.pem:ro \
	-e TLS_CERT_FILE=/etc/ssl/certs/fullchain.pem \
	-e TLS_KEY_FILE=/etc/ssl/private/privkey.pem \
	mikedarling/emcoded.calendar.be:latest
```