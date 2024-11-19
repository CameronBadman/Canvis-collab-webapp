#!/bin/bash

set -e  # Exit immediately if a command exits with a non-zero status
set -u  # Treat unset variables as an error

# Define namespaces and resources
DB_NAMESPACE="db"
BACKEND_NAMESPACE="backend"
ENV_FILE=".env"

# Function to create a namespace if it doesn't exist
create_namespace() {
  local namespace=$1
  if kubectl get namespace "$namespace" >/dev/null 2>&1; then
    echo "Namespace '$namespace' already exists."
  else
    echo "Creating namespace '$namespace'..."
    kubectl create namespace "$namespace"
  fi
}

# Function to create or replace a ConfigMap
apply_configmap() {
  local name=$1
  local namespace=$2
  local file=$3
  echo "Applying ConfigMap '$name' in namespace '$namespace'..."
  kubectl create configmap "$name" --from-env-file="$file" -n "$namespace" --dry-run=client -o yaml | kubectl apply -f -
}

# Function to create or replace a Secret
apply_secret() {
  local name=$1
  local namespace=$2
  shift 2
  echo "Applying Secret '$name' in namespace '$namespace'..."
  kubectl create secret generic "$name" "$@" -n "$namespace" --dry-run=client -o yaml | kubectl apply -f -
}

# Main script starts here

# 1. Ensure namespaces exist
create_namespace "$DB_NAMESPACE"
create_namespace "$BACKEND_NAMESPACE"

# 2. Ensure REDIS_HOST and REDIS_PORT are in the .env file for Redis
if [[ -f "$ENV_FILE" ]]; then
  # Add the REDIS_HOST entry to the ENV_FILE if not already present
  if ! grep -q "REDIS_HOST=" "$ENV_FILE"; then
    echo "Adding REDIS_HOST to $ENV_FILE"
    echo "REDIS_HOST=redis.backend.svc.cluster.local" >> "$ENV_FILE"
  fi

  # Add the REDIS_PORT entry to the ENV_FILE if not already present
  if ! grep -q "REDIS_PORT=" "$ENV_FILE"; then
    echo "Adding REDIS_PORT to $ENV_FILE"
    echo "REDIS_PORT=6379" >> "$ENV_FILE"  # Default Redis port, change if needed
  fi
else
  echo "Error: '$ENV_FILE' not found. Please ensure the file exists in the current directory."
  exit 1
fi

# 3. Apply ConfigMaps and Secrets
# ConfigMap for Cassandra
apply_configmap "cassandra-config" "$DB_NAMESPACE" "$ENV_FILE"
apply_configmap "cassandra-config" "$BACKEND_NAMESPACE" "$ENV_FILE"

# Secret for Cassandra
CASSANDRA_PASSWORD=$(grep -w "CASSANDRA_PASSWORD" "$ENV_FILE" | cut -d '=' -f2)
apply_secret "cassandra-secrets" "$DB_NAMESPACE" --from-literal=CASSANDRA_PASSWORD="$CASSANDRA_PASSWORD"
apply_secret "cassandra-secrets" "$BACKEND_NAMESPACE" --from-literal=CASSANDRA_PASSWORD="$CASSANDRA_PASSWORD"

# ConfigMap for Redis
apply_configmap "redis-config" "$BACKEND_NAMESPACE" "$ENV_FILE"

# Secret for Redis
REDIS_PASSWORD=$(grep -w "REDIS_PASSWORD" "$ENV_FILE" | cut -d '=' -f2)
apply_secret "backend-redis-secret" "$BACKEND_NAMESPACE" --from-literal=REDIS_PASSWORD="$REDIS_PASSWORD"

# ConfigMap for AWS Cognito
apply_configmap "cognito-config" "$BACKEND_NAMESPACE" "$ENV_FILE"

# Secret for AWS Cognito
COGNITO_USER_POOL_ID=$(grep -w "COGNITO_USER_POOL_ID" "$ENV_FILE" | cut -d '=' -f2)
COGNITO_APP_CLIENT_ID=$(grep -w "COGNITO_APP_CLIENT_ID" "$ENV_FILE" | cut -d '=' -f2)
COGNITO_APP_CLIENT_SECRET=$(grep -w "COGNITO_APP_CLIENT_SECRET" "$ENV_FILE" | cut -d '=' -f2)
AWS_REGION=$(grep -w "AWS_REGION" "$ENV_FILE" | cut -d '=' -f2)

apply_secret "cognito-secret" "$BACKEND_NAMESPACE" \
  --from-literal=COGNITO_USER_POOL_ID="$COGNITO_USER_POOL_ID" \
  --from-literal=COGNITO_APP_CLIENT_ID="$COGNITO_APP_CLIENT_ID" \
  --from-literal=COGNITO_APP_CLIENT_SECRET="$COGNITO_APP_CLIENT_SECRET" \
  --from-literal=AWS_REGION="$AWS_REGION"

echo "Setup completed successfully!"
