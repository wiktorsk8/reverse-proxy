#!/bin/bash

set -e

CONTAINER_NAME="reverse-proxy"
IMAGE_NAME="reverse-proxy-img"
PORT="8000"

echo "Building image: $IMAGE_NAME..."
docker build . -t "$IMAGE_NAME"

echo "Removing existing containers named $CONTAINER_NAME"
docker rm -f "$CONTAINER_NAME" 2>/dev/null || true

docker run \
  -d \
  --name "$CONTAINER_NAME" \
  -p "$PORT:$PORT"\
  "$IMAGE_NAME"

echo "Exposed port: $PORT"
echo "Container $CONTAINER_NAME is up and running! ✅"