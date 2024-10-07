#!/bin/bash

# Exit immediately if a command exits with a non-zero status.
set -e

# Configuration
IMAGE_NAME="auth-service"
CONTAINER_NAME="auth-service-container"
PORT=3000

echo "Building Docker image..."
docker build -t $IMAGE_NAME .

echo "Running Docker container..."
docker run -d --name $CONTAINER_NAME -p $PORT:3000 $IMAGE_NAME

echo "Container started. You can access the service at http://localhost:$PORT"
echo "To view logs, run: docker logs $CONTAINER_NAME"
echo "To stop the container, run: docker stop $CONTAINER_NAME"