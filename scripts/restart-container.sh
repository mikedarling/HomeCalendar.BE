#!/bin/bash

# restart-container.sh - Stop, remove, rebuild, and restart the emcoded.calendar.be container

set -e  # Exit on any error

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# Get the project root (parent of scripts directory)
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

CONTAINER_NAME="emcoded_calendar_be_tls"
IMAGE_NAME="mikedarling/emcoded.calendar.be:latest"
CERT_PATH="/tmp/certs"

echo "🔄 Starting container restart process..."
echo "📁 Project root: $PROJECT_ROOT"

# Function to check if container exists
container_exists() {
    docker ps -a --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"
}

# Function to check if container is running
container_running() {
    docker ps --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"
}

# Stop container if running
if container_running; then
    echo "🛑 Stopping running container: $CONTAINER_NAME"
    docker stop "$CONTAINER_NAME"
else
    echo "ℹ️  Container $CONTAINER_NAME is not currently running"
fi

# Remove container if it exists
if container_exists; then
    echo "🗑️  Removing existing container: $CONTAINER_NAME"
    docker rm "$CONTAINER_NAME"
else
    echo "ℹ️  Container $CONTAINER_NAME does not exist"
fi

# Build the image
echo "🔨 Building Docker image: $IMAGE_NAME"
cd "$PROJECT_ROOT"
docker build -t "$IMAGE_NAME" .

# Ensure certificates exist
if [ ! -f "$CERT_PATH/fullchain.pem" ] || [ ! -f "$CERT_PATH/privkey.pem" ]; then
    echo "📜 Creating self-signed certificates at $CERT_PATH"
    mkdir -p "$CERT_PATH"
    openssl req -x509 -newkey rsa:4096 -keyout "$CERT_PATH/privkey.pem" -out "$CERT_PATH/fullchain.pem" -days 365 -nodes -subj "/C=US/ST=Local/L=Local/O=Dev/CN=localhost" >/dev/null 2>&1
    echo "✅ Certificates created"
else
    echo "✅ Certificates already exist at $CERT_PATH"
fi

# Start the container
echo "🚀 Starting new container: $CONTAINER_NAME"
docker run -d --name "$CONTAINER_NAME" \
    -p 8443:8443 \
    -v "$CERT_PATH/fullchain.pem:/etc/ssl/certs/fullchain.pem:ro" \
    -v "$CERT_PATH/privkey.pem:/etc/ssl/private/privkey.pem:ro" \
    -e TLS_CERT_FILE=/etc/ssl/certs/fullchain.pem \
    -e TLS_KEY_FILE=/etc/ssl/private/privkey.pem \
    "$IMAGE_NAME"

# Wait a moment for container to start
sleep 2

# Check if container is healthy
echo "🏥 Checking container health..."
if docker ps --filter "name=$CONTAINER_NAME" --filter "status=running" | grep -q "$CONTAINER_NAME"; then
    echo "✅ Container $CONTAINER_NAME is running successfully!"
    echo "🌐 Access the service at: https://localhost:8443"
    
    # Show container status
    echo ""
    echo "📊 Container status:"
    docker ps --filter "name=$CONTAINER_NAME" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
else
    echo "❌ Container failed to start. Checking logs..."
    docker logs "$CONTAINER_NAME"
    exit 1
fi

echo ""
echo "🎉 Container restart completed successfully!"