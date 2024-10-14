#!/bin/bash

# Configuration
AUTH_URL="http://localhost:8000/api/auth"
API_URL="http://localhost:8000/api/canvas"
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

# Function to make a curl request
make_request() {
    local method=$1
    local url=$2
    shift 2
    local response=$(curl -s -X "$method" -w "\n%{http_code}" "$url" "$@")
    local status_code=$(echo "$response" | tail -n1)
    local body=$(echo "$response" | sed '$d')
    echo "$status_code"
    echo "$body"
}

# Register user
print_color $YELLOW "Registering user..."
REGISTER_RESPONSE=$(make_request POST "$AUTH_URL/register" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}")

REGISTER_STATUS=$(echo "$REGISTER_RESPONSE" | head -n1)
REGISTER_BODY=$(echo "$REGISTER_RESPONSE" | tail -n1)

if [ "$REGISTER_STATUS" -eq 201 ]; then
    print_color $GREEN "User registration successful!"
else
    print_color $RED "User registration failed. Status: $REGISTER_STATUS"
    echo "Response: $REGISTER_BODY"
    exit 1
fi

# Login user
print_color $YELLOW "Logging in user..."
LOGIN_RESPONSE=$(make_request POST "$AUTH_URL/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}")

LOGIN_STATUS=$(echo "$LOGIN_RESPONSE" | head -n1)
LOGIN_BODY=$(echo "$LOGIN_RESPONSE" | tail -n1)

if [ "$LOGIN_STATUS" -eq 200 ]; then
    print_color $GREEN "User login successful!"
    TOKEN=$(echo "$LOGIN_BODY" | jq -r '.token')
    FIREBASE_UID=$(echo "$LOGIN_BODY" | jq -r '.uid')
    print_color $BLUE "Token: $TOKEN"
    print_color $BLUE "Firebase UID: $FIREBASE_UID"
else
    print_color $RED "User login failed. Status: $LOGIN_STATUS"
    echo "Response: $LOGIN_BODY"
    exit 1
fi

# Function to make an authenticated request
make_auth_request() {
    local method=$1
    local endpoint=$2
    shift 2
    make_request "$method" "$API_URL$endpoint" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $TOKEN" \
        -H "X-Firebase-UID: $FIREBASE_UID" \
        "$@"
}

# Test CreateUser
print_color $YELLOW "Testing CreateUser..."
CREATE_USER_RESPONSE=$(make_auth_request POST "/user" -d "{\"firebase_uid\":\"$FIREBASE_UID\"}")
CREATE_USER_STATUS=$(echo "$CREATE_USER_RESPONSE" | head -n1)
CREATE_USER_BODY=$(echo "$CREATE_USER_RESPONSE" | tail -n1)

if [ "$CREATE_USER_STATUS" -eq 201 ]; then
    print_color $GREEN "CreateUser successful!"
else
    print_color $RED "CreateUser failed. Status: $CREATE_USER_STATUS"
    echo "Response: $CREATE_USER_BODY"
fi

# Test GetUser
print_color $YELLOW "Testing GetUser..."
GET_USER_RESPONSE=$(make_auth_request GET "/user")
GET_USER_STATUS=$(echo "$GET_USER_RESPONSE" | head -n1)
GET_USER_BODY=$(echo "$GET_USER_RESPONSE" | tail -n1)

if [ "$GET_USER_STATUS" -eq 200 ]; then
    print_color $GREEN "GetUser successful!"
else
    print_color $RED "GetUser failed. Status: $GET_USER_STATUS"
    echo "Response: $GET_USER_BODY"
fi

# Test CreateCanvas
print_color $YELLOW "Testing CreateCanvas..."
CREATE_CANVAS_RESPONSE=$(make_auth_request POST "/canvas" -d "{\"name\":\"Test Canvas\"}")
CREATE_CANVAS_STATUS=$(echo "$CREATE_CANVAS_RESPONSE" | head -n1)
CREATE_CANVAS_BODY=$(echo "$CREATE_CANVAS_RESPONSE" | tail -n1)

if [ "$CREATE_CANVAS_STATUS" -eq 201 ]; then
    print_color $GREEN "CreateCanvas successful!"
    CANVAS_ID=$(echo "$CREATE_CANVAS_BODY" | jq -r '.id')
    print_color $BLUE "Canvas ID: $CANVAS_ID"
else
    print_color $RED "CreateCanvas failed. Status: $CREATE_CANVAS_STATUS"
    echo "Response: $CREATE_CANVAS_BODY"
    exit 1
fi

# Test GetCanvas
print_color $YELLOW "Testing GetCanvas..."
GET_CANVAS_RESPONSE=$(make_auth_request GET "/canvas/$CANVAS_ID")
GET_CANVAS_STATUS=$(echo "$GET_CANVAS_RESPONSE" | head -n1)
GET_CANVAS_BODY=$(echo "$GET_CANVAS_RESPONSE" | tail -n1)

if [ "$GET_CANVAS_STATUS" -eq 200 ]; then
    print_color $GREEN "GetCanvas successful!"
else
    print_color $RED "GetCanvas failed. Status: $GET_CANVAS_STATUS"
    echo "Response: $GET_CANVAS_BODY"
fi

# Test UpdateCanvas
print_color $YELLOW "Testing UpdateCanvas..."
UPDATE_CANVAS_RESPONSE=$(make_auth_request PUT "/canvas/$CANVAS_ID" -d "{\"name\":\"Updated Test Canvas\"}")
UPDATE_CANVAS_STATUS=$(echo "$UPDATE_CANVAS_RESPONSE" | head -n1)
UPDATE_CANVAS_BODY=$(echo "$UPDATE_CANVAS_RESPONSE" | tail -n1)

if [ "$UPDATE_CANVAS_STATUS" -eq 200 ]; then
    print_color $GREEN "UpdateCanvas successful!"
else
    print_color $RED "UpdateCanvas failed. Status: $UPDATE_CANVAS_STATUS"
    echo "Response: $UPDATE_CANVAS_BODY"
fi

# Test GetUserCanvases
print_color $YELLOW "Testing GetUserCanvases..."
GET_USER_CANVASES_RESPONSE=$(make_auth_request GET "/user/canvases")
GET_USER_CANVASES_STATUS=$(echo "$GET_USER_CANVASES_RESPONSE" | head -n1)
GET_USER_CANVASES_BODY=$(echo "$GET_USER_CANVASES_RESPONSE" | tail -n1)

if [ "$GET_USER_CANVASES_STATUS" -eq 200 ]; then
    print_color $GREEN "GetUserCanvases successful!"
else
    print_color $RED "GetUserCanvases failed. Status: $GET_USER_CANVASES_STATUS"
    echo "Response: $GET_USER_CANVASES_BODY"
fi

# Test DeleteCanvas
print_color $YELLOW "Testing DeleteCanvas..."
DELETE_CANVAS_RESPONSE=$(make_auth_request DELETE "/canvas/$CANVAS_ID")
DELETE_CANVAS_STATUS=$(echo "$DELETE_CANVAS_RESPONSE" | head -n1)
DELETE_CANVAS_BODY=$(echo "$DELETE_CANVAS_RESPONSE" | tail -n1)

if [ "$DELETE_CANVAS_STATUS" -eq 204 ]; then
    print_color $GREEN "DeleteCanvas successful!"
else
    print_color $RED "DeleteCanvas failed. Status: $DELETE_CANVAS_STATUS"
    echo "Response: $DELETE_CANVAS_BODY"
fi

# Test DeleteUser
print_color $YELLOW "Testing DeleteUser..."
DELETE_USER_RESPONSE=$(make_auth_request DELETE "/user")
DELETE_USER_STATUS=$(echo "$DELETE_USER_RESPONSE" | head -n1)
DELETE_USER_BODY=$(echo "$DELETE_USER_RESPONSE" | tail -n1)

if [ "$DELETE_USER_STATUS" -eq 204 ]; then
    print_color $GREEN "DeleteUser successful!"
else
    print_color $RED "DeleteUser failed. Status: $DELETE_USER_STATUS"
    echo "Response: $DELETE_USER_BODY"
fi

# Logout user
print_color $YELLOW "Logging out user..."
LOGOUT_RESPONSE=$(make_request POST "$AUTH_URL/logout" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json")

LOGOUT_STATUS=$(echo "$LOGOUT_RESPONSE" | head -n1)
LOGOUT_BODY=$(echo "$LOGOUT_RESPONSE" | tail -n1)

if [ "$LOGOUT_STATUS" -eq 200 ]; then
    print_color $GREEN "User logout successful!"
else
    print_color $RED "User logout failed. Status: $LOGOUT_STATUS"
    echo "Response: $LOGOUT_BODY"
fi

print_color $YELLOW "Canvas API test completed."