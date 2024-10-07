#!/bin/bash

# API Gateway base URL (now pointing to the load balancer)
BASE_URL="http://localhost:80"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Function to print colored output
print_result() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}Success${NC}"
    else
        echo -e "${RED}Failed${NC}"
    fi
}

echo "Testing Load-Balanced API Gateway and Auth Service"

# Test health endpoint (if implemented in API Gateway)
echo -n "Testing health endpoint: "
HEALTH_RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" $BASE_URL/health)
if [ "$HEALTH_RESPONSE" == "200" ]; then
    echo -e "${GREEN}Success${NC}"
else
    echo -e "${RED}Failed (HTTP $HEALTH_RESPONSE)${NC}"
fi

# Function to attempt registration or login
attempt_auth() {
    local email="testuser$RANDOM@example.com"
    local password="testpassword123"

    echo "Attempting to register/login user: $email"
    
    REGISTER_RESPONSE=$(curl -s -X POST $BASE_URL/auth/register \
      -H "Content-Type: application/json" \
      -d "{\"email\": \"$email\", \"password\": \"$password\"}")

    if [[ $REGISTER_RESPONSE == *"User registered successfully"* ]]; then
        echo "User registered successfully"
        RESPONSE=$REGISTER_RESPONSE
    elif [[ $REGISTER_RESPONSE == *"Email already in use"* ]]; then
        echo "User already exists. Proceeding with login."
        LOGIN_RESPONSE=$(curl -s -X POST $BASE_URL/auth/login \
          -H "Content-Type: application/json" \
          -d "{\"email\": \"$email\", \"password\": \"$password\"}")
        RESPONSE=$LOGIN_RESPONSE
    else
        echo "Unexpected response: $REGISTER_RESPONSE"
        return 1
    fi

    # Extract UID and token from response
    USER_UID=$(echo $RESPONSE | grep -o '"uid":"[^"]*' | grep -o '[^"]*$')
    USER_TOKEN=$(echo $RESPONSE | grep -o '"token":"[^"]*' | grep -o '[^"]*$')

    if [ -z "$USER_TOKEN" ] || [ -z "$USER_UID" ]; then
        echo -e "${RED}Failed to extract token or UID from response${NC}"
        return 1
    fi

    echo "Extracted UID: $USER_UID"
    echo "Extracted token: $USER_TOKEN"
    return 0
}

# Attempt authentication
attempt_auth
if [ $? -ne 0 ]; then
    echo "Authentication failed. Exiting."
    exit 1
fi

# Logout
echo -n "Logging out: "
LOGOUT_RESPONSE=$(curl -s -X POST $BASE_URL/auth/logout \
  -H "Authorization: Bearer $USER_TOKEN")
echo $LOGOUT_RESPONSE
print_result $?

# Try to use the token after logout (should fail)
echo -n "Attempting to use token after logout: "
CHECKTOKEN_RESPONSE=$(curl -s -X GET $BASE_URL/auth/check-token/$USER_UID \
  -H "Authorization: Bearer $USER_TOKEN")
echo $CHECKTOKEN_RESPONSE
if [[ $CHECKTOKEN_RESPONSE == *"No token found for this user"* ]]; then
    echo -e "${GREEN}Token successfully invalidated${NC}"
else
    echo -e "${RED}Token still valid after logout${NC}"
fi

echo "Testing complete"