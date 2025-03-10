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

# Function to delete all ConfigMaps in a namespace
delete_configmaps() {
  local namespace=$1
  echo "Deleting all ConfigMaps in namespace '$namespace'..."
  kubectl delete configmap --all -n "$namespace"
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

# 2. Validate ENV_FILE
if [[ ! -f "$ENV_FILE" ]]; then
  echo "Error: '$ENV_FILE' not found. Please ensure the file exists in the current directory."
  exit 1
fi

# 3. Extract sensitive data from the ENV_FILE
JWT_SECRET_KEY=$(grep -w "JWT_SECRET_KEY" "$ENV_FILE" | cut -d '=' -f2)
AUTH_REDIS_PASSWORD=$(grep -w "AUTH_REDIS_PASSWORD" "$ENV_FILE" | cut -d '=' -f2)
COGNITO_USER_POOL_ID=$(grep -w "COGNITO_USER_POOL_ID" "$ENV_FILE" | cut -d '=' -f2)
COGNITO_APP_CLIENT_ID=$(grep -w "COGNITO_APP_CLIENT_ID" "$ENV_FILE" | cut -d '=' -f2)
COGNITO_APP_CLIENT_SECRET=$(grep -w "COGNITO_APP_CLIENT_SECRET" "$ENV_FILE" | cut -d '=' -f2)
AWS_REGION=$(grep -w "AWS_REGION" "$ENV_FILE" | cut -d '=' -f2)

# New drawing Redis variables
DRAWING_REDIS_HOST=$(grep -w "DRAWING_REDIS_HOST" "$ENV_FILE" | cut -d '=' -f2)
DRAWING_REDIS_PORT=$(grep -w "DRAWING_REDIS_PORT" "$ENV_FILE" | cut -d '=' -f2)
DRAWING_REDIS_PASSWORD=$(grep -w "DRAWING_REDIS_PASSWORD" "$ENV_FILE" | cut -d '=' -f2)

# Check for missing required sensitive variables and print specific error messages
if [[ -z "$JWT_SECRET_KEY" ]]; then
  echo "Error: JWT_SECRET_KEY is missing in $ENV_FILE."
  exit 1
fi

if [[ -z "$AUTH_REDIS_PASSWORD" ]]; then
  echo "Error: AUTH_REDIS_PASSWORD is missing in $ENV_FILE."
  exit 1
fi

if [[ -z "$COGNITO_APP_CLIENT_SECRET" ]]; then
  echo "Error: COGNITO_APP_CLIENT_SECRET is missing in $ENV_FILE."
  exit 1
fi

if [[ -z "$DRAWING_REDIS_PASSWORD" ]]; then
  echo "Error: DRAWING_REDIS_PASSWORD is missing in $ENV_FILE."
  exit 1
fi

# 4. Delete existing ConfigMaps
delete_configmaps "$BACKEND_NAMESPACE"

# 5. Apply ConfigMaps (for non-sensitive data)
echo "Creating ConfigMaps for canvas, account, and general configuration..."
apply_configmap "canvas-config" "$BACKEND_NAMESPACE" "$ENV_FILE"
apply_configmap "account-config" "$BACKEND_NAMESPACE" "$ENV_FILE"
apply_configmap "redis-config" "$BACKEND_NAMESPACE" "$ENV_FILE"
apply_configmap "cognito-config" "$BACKEND_NAMESPACE" "$ENV_FILE"
apply_configmap "cassandra-config" "$BACKEND_NAMESPACE" "$ENV_FILE"

# 6. Apply Secrets (for sensitive data)
echo "Creating Secrets for sensitive values..."
apply_secret "jwt-secret" "$BACKEND_NAMESPACE" --from-literal=JWT_SECRET_KEY="$JWT_SECRET_KEY"
apply_secret "backend-redis-secret" "$BACKEND_NAMESPACE" --from-literal=REDIS_PASSWORD="$AUTH_REDIS_PASSWORD"
apply_secret "cognito-secret" "$BACKEND_NAMESPACE" \
  --from-literal=COGNITO_USER_POOL_ID="$COGNITO_USER_POOL_ID" \
  --from-literal=COGNITO_APP_CLIENT_ID="$COGNITO_APP_CLIENT_ID" \
  --from-literal=COGNITO_APP_CLIENT_SECRET="$COGNITO_APP_CLIENT_SECRET" \
  --from-literal=AWS_REGION="$AWS_REGION"

# New secret for drawing Redis
apply_secret "drawing-redis-secret" "redis-drawing" --from-literal=REDIS_PASSWORD="$DRAWING_REDIS_PASSWORD"


echo "Setup completed successfully!"
