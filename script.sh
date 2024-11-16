#!/bin/bash

set -e  # Exit immediately if a command fails

echo "Creating namespaces if they don't exist..."
kubectl create namespace db || echo "Namespace 'db' already exists"
kubectl create namespace backend || echo "Namespace 'backend' already exists"

echo "Creating ConfigMap and Secret in 'db' namespace..."
# Cassandra ConfigMap
kubectl create configmap cassandra-config --from-env-file=.env -n db || echo "ConfigMap 'cassandra-config' already exists in 'db'"

# Cassandra Secret
kubectl create secret generic cassandra-secrets --from-literal=CASSANDRA_PASSWORD=cassandra1 -n db || echo "Secret 'cassandra-secrets' already exists in 'db'"

echo "Copying Cassandra ConfigMap and Secret to 'backend' namespace..."
# Cassandra ConfigMap
kubectl get configmap cassandra-config -n db -o yaml | sed "s/namespace: db/namespace: backend/" | kubectl apply -f -

# Cassandra Secret
kubectl get secret cassandra-secrets -n db -o yaml | sed "s/namespace: db/namespace: backend/" | kubectl apply -f -

echo "Creating Redis ConfigMap and Secret in 'backend' namespace..."
# Redis ConfigMap
kubectl create configmap redis-config --from-literal=REDIS_HOST=redis.backend.svc.cluster.local --from-literal=REDIS_PORT=6379 -n backend || echo "ConfigMap 'redis-config' already exists in 'backend'"

# Redis Secret
kubectl create secret generic backend-redis-secret --from-literal=REDIS_PASSWORD=password -n backend || echo "Secret 'backend-redis-secret' already exists in 'backend'"

echo "Namespace setup completed successfully."
