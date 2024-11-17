#!/bin/bash

# Variables
NAMESPACE="db"  # Namespace where Cassandra is deployed
CONFIGMAP_NAME="cassandra-schema"  # Name of the ConfigMap
SCHEMA_FILE="schema.cql"  # Path to the schema file

# Check if the schema file exists
if [ ! -f "$SCHEMA_FILE" ]; then
  echo "Error: Schema file '$SCHEMA_FILE' not found. Please ensure it exists in the current directory."
  exit 1
fi

# Create the namespace if it doesn't exist
kubectl get namespace $NAMESPACE >/dev/null 2>&1
if [ $? -ne 0 ]; then
  echo "Namespace '$NAMESPACE' does not exist. Creating it..."
  kubectl create namespace $NAMESPACE
fi

# Create the ConfigMap
echo "Creating ConfigMap '$CONFIGMAP_NAME' in namespace '$NAMESPACE'..."
kubectl create configmap $CONFIGMAP_NAME \
  --from-file=$SCHEMA_FILE \
  --namespace=$NAMESPACE \
  --dry-run=client -o yaml | kubectl apply -f -

# Verify the ConfigMap creation
if [ $? -eq 0 ]; then
  echo "ConfigMap '$CONFIGMAP_NAME' created successfully in namespace '$NAMESPACE'."
else
  echo "Failed to create ConfigMap '$CONFIGMAP_NAME'."
  exit 1
fi

# List ConfigMaps in the namespace
echo "Listing ConfigMaps in namespace '$NAMESPACE':"
kubectl get configmaps -n $NAMESPACE
