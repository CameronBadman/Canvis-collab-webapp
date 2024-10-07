#!/bin/bash

# Exit immediately if a command exits with a non-zero status.
set -e

# Function to check if a command exists
command_exists() {
  command -v "$1" >/dev/null 2>&1
}

# Check for required commands
for cmd in docker kubectl; do
  if ! command_exists "$cmd"; then
    echo "Error: $cmd is not installed. Please install it and try again."
    exit 1
  fi
done

# Configuration
AUTH_SERVICE_DIR="./backend/auth-service"
KUBE_CONFIG_FILE="./auth-service-kubernetes.yaml"

echo "Starting deployment process..."

# Check if the necessary directories and files exist
if [ ! -d "$AUTH_SERVICE_DIR" ]; then
  echo "Error: Auth service directory not found. Are you in the project root?"
  exit 1
fi

if [ ! -f "$KUBE_CONFIG_FILE" ]; then
  echo "Error: Kubernetes configuration file not found. Are you in the project root?"
  exit 1
fi

# Build Docker image
echo "Building Docker image..."
docker build -t auth-service:latest "$AUTH_SERVICE_DIR"

# Verify image was built
if docker images | grep -q auth-service; then
  echo "Docker image built successfully."
else
  echo "Error: Failed to build Docker image."
  exit 1
fi

# Apply Kubernetes configuration
echo "Applying Kubernetes configuration..."
kubectl apply -f "$KUBE_CONFIG_FILE"

# Wait for pods to be ready
echo "Waiting for pods to be ready..."
kubectl wait --for=condition=ready pod -l app=auth-service --timeout=60s || echo "Warning: Timeout waiting for auth-service pod"
kubectl wait --for=condition=ready pod -l app=redis --timeout=60s || echo "Warning: Timeout waiting for redis pod"

# Check pod status
echo "Checking pod status..."
kubectl get pods

echo "Deployment process completed."

# Port forward for testing (uncomment if needed)
# echo "Setting up port forwarding..."
# kubectl port-forward service/auth-service 3000:3000 &
# echo "Auth service available at http://localhost:3000"