#!/bin/bash

# Configuration
LOAD_BALANCER_URL="http://localhost:8000"  # Load balancer URL
EMAIL="testuser$(date +%s)@example.com"
PASSWORD="TestPassword123!"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_color() {
    color=$1
    message=$2
    echo -e "${color}${message}${NC}"
}

# Function to make a curl request and handle potential connection errors
make_request() {
    local method=$1
    local url=$2
    shift 2
    echo "Making request to: $method $url" >&2
    local response=$(curl -s -X "$method" -w "\n%{http_code}" "$url" "$@")
    local status_code=$(echo "$response" | tail -n1)
    local body=$(echo "$response" | sed '$d')

    print_color $BLUE "Status Code: $status_code" >&2
    print_color $BLUE "Response Body: $body" >&2

    if [[ ! "$status_code" =~ ^[0-9]+$ ]]; then
        print_color $RED "Invalid status code. Is the load balancer running?" >&2
        return 1
    fi

    if [ "$status_code" -eq 000 ]; then
        print_color $RED "Connection error. Is the load balancer running?" >&2
        return 1
    fi

    echo "$status_code"
    echo "$body"
}

# Main script
print_color $YELLOW "Testing Account Creation via Load Balancer"
print_color $YELLOW "Load Balancer URL: $LOAD_BALANCER_URL"

# Health check for the load balancer
print_color $YELLOW "Performing health check..."
HEALTH_RESPONSE=$(make_request GET "$LOAD_BALANCER_URL/api/health")
if [ $? -ne 0 ]; then
    print_color $RED "Health check failed. Unable to connect to the load balancer."
    exit 1
fi

HEALTH_STATUS_CODE=$(echo "$HEALTH_RESPONSE" | head -n1)
HEALTH_BODY=$(echo "$HEALTH_RESPONSE" | tail -n1)

if [ "$HEALTH_STATUS_CODE" -eq 200 ]; then
    print_color $GREEN "Health check successful!"
else
    print_color $RED "Health check failed. Status code: $HEALTH_STATUS_CODE"
    print_color $RED "Response body: $HEALTH_BODY"
    exit 1
fi

# Sending registration request
print_color $YELLOW "Sending registration request..."
REGISTER_RESPONSE=$(make_request POST "$LOAD_BALANCER_URL/api/auth/register" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}")

REGISTER_STATUS_CODE=$(echo "$REGISTER_RESPONSE" | head -n1)
REGISTER_BODY=$(echo "$REGISTER_RESPONSE" | tail -n1)

if [ "$REGISTER_STATUS_CODE" -eq 200 ] || [ "$REGISTER_STATUS_CODE" -eq 201 ]; then
    print_color $GREEN "Account creation successful!"
else
    print_color $RED "Account creation failed. Status code: $REGISTER_STATUS_CODE"
    echo "Response: $REGISTER_BODY"
    exit 1
fi

# Sending login request
print_color $YELLOW "Sending login request..."
LOGIN_RESPONSE=$(make_request POST "$LOAD_BALANCER_URL/api/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}")

LOGIN_STATUS_CODE=$(echo "$LOGIN_RESPONSE" | head -n1)
LOGIN_BODY=$(echo "$LOGIN_RESPONSE" | tail -n1)

if [ "$LOGIN_STATUS_CODE" -eq 200 ] || [ "$LOGIN_STATUS_CODE" -eq 201 ]; then
    print_color $GREEN "Login successful!"
    print_color $YELLOW "Login response body:"
    echo "$LOGIN_BODY"
    
    # Extract token using jq
    TOKEN=$(echo "$LOGIN_BODY" | jq -r '.token')
    if [ -z "$TOKEN" ] || [ "$TOKEN" == "null" ]; then
        print_color $RED "Failed to extract token from response."
        exit 1
    else
        print_color $GREEN "Token successfully extracted: $TOKEN"
    fi
else
    print_color $RED "Login failed. Status code: $LOGIN_STATUS_CODE"
    echo "Response: $LOGIN_BODY"
    exit 1
fi

# Sending logout request
print_color $YELLOW "Sending logout request..."
LOGOUT_RESPONSE=$(make_request POST "$LOAD_BALANCER_URL/api/auth/logout" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json")

LOGOUT_STATUS_CODE=$(echo "$LOGOUT_RESPONSE" | head -n1)
LOGOUT_BODY=$(echo "$LOGOUT_RESPONSE" | tail -n1)

if [ "$LOGOUT_STATUS_CODE" -eq 200 ]; then
    print_color $GREEN "Logout successful!"
else
    print_color $RED "Logout failed. Status code: $LOGOUT_STATUS_CODE"
    echo "Response: $LOGOUT_BODY"
fi

print_color $YELLOW "Test completed."