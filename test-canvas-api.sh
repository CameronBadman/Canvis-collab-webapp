#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Base URL of your API
# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Base URL of your API
API_URL="http://localhost:8080"

# Function to print colored output
print_result() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}PASS${NC}: $2"
    else
        echo -e "${RED}FAIL${NC}: $2"
    fi
}

# Test 1: Create a new canvas
echo "Test 1: Creating a new canvas"
echo "Sending POST request to ${API_URL}/api/canvas/create"
RESPONSE=$(curl -v -X POST "${API_URL}/api/canvas/create" 2>&1)
echo "Full curl output:"
echo "$RESPONSE"

CANVAS_ID=$(echo $RESPONSE | grep -o '"canvasId":"[^"]*' | sed 's/"canvasId":"//')

if [[ $CANVAS_ID =~ ^[A-Za-z0-9]{9}$ ]]; then
    print_result 0 "Canvas created successfully with ID: $CANVAS_ID"
else
    print_result 1 "Failed to create canvas. Response: $RESPONSE"
    exit 1
fi

# Test 2: Retrieve the created canvas
echo "Test 2: Retrieving the created canvas"
RESPONSE=$(curl -s -X GET "${API_URL}/api/canvas/${CANVAS_ID}")
RETRIEVED_ID=$(echo $RESPONSE | jq -r '.canvasId')

if [ "$RETRIEVED_ID" == "$CANVAS_ID" ]; then
    print_result 0 "Canvas retrieved successfully"
else
    print_result 1 "Failed to retrieve canvas. Response: $RESPONSE"
    exit 1
fi

# Test 3: Attempt to retrieve a non-existent canvas
echo "Test 3: Attempting to retrieve a non-existent canvas"
RESPONSE=$(curl -s -X GET "${API_URL}/api/canvas/nonexistent")
ERROR_MESSAGE=$(echo $RESPONSE | jq -r '.error')

if [ "$ERROR_MESSAGE" == "Canvas not found" ]; then
    print_result 0 "Correctly handled non-existent canvas"
else
    print_result 1 "Unexpected response for non-existent canvas. Response: $RESPONSE"
    exit 1
fi

echo "All tests completed successfully!"