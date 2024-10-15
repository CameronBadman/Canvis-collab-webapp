#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to check and install dependencies
check_dependencies() {
    if ! command_exists docker; then
        echo "Docker is not installed. Please install Docker and try again."
        echo "Visit https://docs.docker.com/get-docker/ for installation instructions."
        exit 1
    fi

    if ! command_exists kubectl; then
        echo "kubectl is not installed. Please install kubectl and try again."
        echo "Visit https://kubernetes.io/docs/tasks/tools/install-kubectl/ for installation instructions."
        exit 1
    fi

    if ! command_exists helm; then
        echo "Helm is not installed. Would you like to install it? (y/n)"
        read -r response
        if [[ "$response" =~ ^([yY][eE][sS]|[yY])+$ ]]; then
            echo "Installing Helm..."
            curl https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | bash
        else
            echo "Helm is required to install Gloo Edge. Please install it and try again."
            exit 1
        fi
    fi
}

cleanup_kubernetes() {
    echo "Cleaning up existing Kubernetes deployments..."
    
    # Get the list of resources created by our application
    resources=(
        "deployment/auth-service"
        "deployment/frontend"
        "deployment/canvas-api"
        "deployment/redis"
        "statefulset/cassandra"
        "service/auth-service"
        "service/frontend"
        "service/canvas-api"
        "service/redis"
        "service/cassandra"
        "configmap/app-config"
        "secret/app-secrets"
    )

    # Delete specific resources
    for resource in "${resources[@]}"
    do
        echo "Deleting $resource..."
        kubectl delete $resource --ignore-not-found=true
    done

    # Wait for all pods related to our application to terminate
    echo "Waiting for all application pods to terminate..."
    kubectl wait --for=delete pod -l "app in (auth-service,frontend,canvas-api,redis,cassandra)" --timeout=60s || true

    echo "Cleanup complete."
}

# Function to build Docker images
build_images() {
    echo "Building auth-service image..."
    docker build -t auth-service:latest ./backend/auth-service

    echo "Building frontend image..."
    docker build -t frontend:latest ./frontend

    echo "Building canvas-api image..."
    docker build -t canvas-api:latest ./backend/canvas-api
}

install_gloo_edge() {
    echo "Installing Gloo Edge..."
    helm repo add gloo https://storage.googleapis.com/solo-public-helm
    helm repo update
    kubectl create namespace gloo-system --dry-run=client -o yaml | kubectl apply -f -

    # Install Gloo Edge with minimal resources
    helm upgrade --install gloo gloo/gloo \
        --namespace gloo-system \
        --set gateway.enabled=true \
        --set discovery.enabled=false \
        --set global.glooRbac.create=false \
        --set global.extensions.caching.enabled=false \
        --set global.extensions.extauth.enabled=false \
        --set global.extensions.ratelimit.enabled=false \
        --set gloo.deployment.resources.requests.memory=128Mi \
        --set gloo.deployment.resources.limits.memory=256Mi \
        --set gateway.deployment.resources.requests.memory=128Mi \
        --set gateway.deployment.resources.limits.memory=256Mi \
        --set gateway-proxy.deployment.resources.requests.memory=128Mi \
        --set gateway-proxy.deployment.resources.limits.memory=256Mi \
        --wait --timeout 5m

    echo "Waiting for Gloo Edge deployments to be ready..."
    kubectl wait --for=condition=available --timeout=300s deployment/gloo -n gloo-system
    kubectl wait --for=condition=available --timeout=300s deployment/gateway -n gloo-system
    kubectl wait --for=condition=available --timeout=300s deployment/gateway-proxy -n gloo-system

    echo "Verifying Gloo Edge installation..."
    kubectl get pods -n gloo-system
    if [ $? -ne 0 ]; then
        echo "Failed to get Gloo Edge pods. Installation may have failed."
        exit 1
    fi

    echo "Gloo Edge installation completed."
}

apply_kubernetes_configs() {
    echo "Applying Kubernetes configurations..."
    kubectl apply -f K8s/config-and-secrets.yaml
    kubectl apply -f K8s/redis.yaml
    kubectl apply -f K8s/cassandra.yaml
    kubectl apply -f K8s/auth-service.yaml
    kubectl apply -f K8s/canvas-api.yaml
    kubectl apply -f K8s/frontend.yaml
    
    echo "Applying Gloo Edge configuration..."
    kubectl apply -f gloo-config.yaml
}

# Function to wait for pods to be ready
wait_for_pods() {
    echo "Waiting for pods to be ready..."
    kubectl wait --for=condition=Ready pods --all --timeout=300s
}

check_deployment_status() {
    echo "Checking deployment status..."
    
    kubectl get pods
    echo ""
    kubectl get services
    echo ""
    kubectl get deployments
    echo ""
    kubectl get virtualservices.gateway.solo.io -n gloo-system
    echo ""
    kubectl get gateways.gateway.solo.io -n gloo-system
}

check_cluster_status() {
    echo "Checking cluster status..."
    kubectl cluster-info
    kubectl get nodes
    kubectl get pods --all-namespaces
}

# Main execution
echo "Starting deployment process..."

check_dependencies

check_cluster_status

# Perform cleanup
cleanup_kubernetes

build_images

install_gloo_edge

# Apply configs and handle potential errors
if ! apply_kubernetes_configs; then
    echo "Error applying Kubernetes configurations. Check the error messages above."
    check_deployment_status
    exit 1
fi

if ! wait_for_pods; then
    echo "Error: Some pods failed to become ready. Checking deployment status:"
    check_deployment_status
    exit 1
fi

echo "Deployment process completed. Checking final status:"
check_deployment_status

echo ""
echo "Gloo Edge API Gateway has been set up successfully."
echo "Your services should now be accessible through the Gloo Edge Gateway."
echo ""
echo "To access your services, use the following URL:"
echo "http://$(kubectl get service -n gloo-system gateway-proxy -o jsonpath='{.status.loadBalancer.ingress[0].ip}'):80"
echo ""
echo "Append the following paths to access specific services:"
echo "/api/auth - For Auth Service"
echo "/api/canvas - For Canvas API"
echo "/ - For Frontend"
echo ""
echo "If any services are missing, please check the deployment status above and the Kubernetes logs for more information."
echo "You can use 'kubectl logs <pod-name> -n gloo-system' to check the logs of Gloo Edge pods."
echo ""
echo "To clean up and delete all resources, run this script again."