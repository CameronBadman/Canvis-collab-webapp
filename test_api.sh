#!/bin/bash

# API endpoint
API_URL="http://localhost:6969"

# Mock Firebase ID
MOCK_FIREBASE_ID="mock-firebase-id-123456"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Function to make API requests
make_request() {
    local method=$1
    local endpoint=$2
    local data=$3

    if [ -n "$data" ]; then
        response=$(curl -s -w "\n%{http_code}" -X $method \
            -H "Content-Type: application/json" \
            -H "X-Firebase-UID: $MOCK_FIREBASE_ID" \
            -d "$data" \
            $API_URL$endpoint)
    else
        response=$(curl -s -w "\n%{http_code}" -X $method \
            -H "X-Firebase-UID: $MOCK_FIREBASE_ID" \
            $API_URL$endpoint)
    fi

    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')

    echo "HTTP Status: $http_code"
    echo "Response body: $body"

    if [ $http_code -ge 200 ] && [ $http_code -lt 300 ]; then
        echo $body
    else
        echo "Request failed with status $http_code: $body"
    fi
}

# Function to check response
check_response() {
    local expected=$1
    local actual=$2
    local test_name=$3

    if [[ $actual == *"$expected"* ]]; then
        echo -e "${GREEN}✓ $test_name passed${NC}"
    else
        echo -e "${RED}✗ $test_name failed${NC}"
        echo "Expected: $expected"
        echo "Actual: $actual"
    fi
}

# Test CreateUser
echo "Testing CreateUser..."
create_user_response=$(make_request POST /user "{\"firebaseUID\":\"$MOCK_FIREBASE_ID\"}")
check_response "\"firebase_uid\":\"$MOCK_FIREBASE_ID\"" "$create_user_response" "CreateUser"

# Test GetUser
echo "Testing GetUser..."
get_user_response=$(make_request GET /user)
check_response "\"firebase_uid\":\"$MOCK_FIREBASE_ID\"" "$get_user_response" "GetUser"

# Test UpdateUser
echo "Testing UpdateUser..."
update_user_response=$(make_request PUT /user '{}')
check_response "\"firebase_uid\":\"$MOCK_FIREBASE_ID\"" "$update_user_response" "UpdateUser"

# Test GetUserCanvases
echo "Testing GetUserCanvases..."
get_canvases_response=$(make_request GET /user/canvases)
check_response "\"firebase_uid\":\"$MOCK_FIREBASE_ID\"" "$get_canvases_response" "GetUserCanvases"

# Test CreateCanvas
echo "Testing CreateCanvas..."
create_canvas_response=$(make_request POST /canvas '{"name":"Test Canvas","svg_data":"<svg></svg>"}')
canvas_id=$(echo $create_canvas_response | jq -r '.id')
check_response "\"name\":\"Test Canvas\"" "$create_canvas_response" "CreateCanvas"

# Test CreateCanvas
echo "Testing CreateCanvas..."
create_canvas_response=$(make_request POST /canvas '{"name":"Test Canvas","svg_data":"<svg></svg>"}')
echo "Raw CreateCanvas response: $create_canvas_response"
canvas_id=$(echo $create_canvas_response | jq -r '.id' || echo "")
check_response "\"name\":\"Test Canvas\"" "$create_canvas_response" "CreateCanvas"

# Test UpdateCanvas
echo "Testing UpdateCanvas..."
update_canvas_response=$(make_request PUT /canvas/$canvas_id '{"name":"Updated Canvas","svg_data":"<svg><circle></circle></svg>"}')
check_response "\"name\":\"Updated Canvas\"" "$update_canvas_response" "UpdateCanvas"

# Test DeleteCanvas
echo "Testing DeleteCanvas..."
delete_canvas_response=$(make_request DELETE /canvas/$canvas_id)
check_response "" "$delete_canvas_response" "DeleteCanvas"

# Test DeleteUser
echo "Testing DeleteUser..."
delete_user_response=$(make_request DELETE /user)
check_response "" "$delete_user_response" "DeleteUser"

echo "API tests completed."