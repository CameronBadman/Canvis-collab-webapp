#!/bin/bash

# Set your JWT secret key
JWT_SECRET_KEY=""  # Replace with your actual secret key

# Generate JWT Token
generate_token() {
  # Define claims as a JSON string
  CLAIMS="{\"sub\":\"f7943ac8-5031-702b-fd3a-cbeb5faee7b1\",\"exp\":$(($(date +%s) + 3600)),\"iat\":$(date +%s)}"

  # Generate the JWT token using jwt-cli
  JWT_TOKEN=$(jwt encode --secret "$JWT_SECRET_KEY" "$CLAIMS")

  echo "Generated JWT Token:"
  echo "$JWT_TOKEN"
}

# Validate JWT Token
validate_token() {
  if [ -z "$1" ]; then
    echo "No JWT token provided for validation."
    exit 1
  fi

  # Validate the JWT token using jwt-cli
  jwt decode --secret "$JWT_SECRET_KEY" "$1" > /dev/null 2>&1

  if [ $? -eq 0 ]; then
    echo "Token is valid."
  else
    echo "Invalid or expired token."
  fi
}

# Main
echo "Welcome to the JWT generation and validation script."

# Generate a JWT token
generate_token

# Ask user if they want to validate the generated token
echo "Would you like to validate the generated token? (y/n)"
read VALIDATE_ANSWER

if [ "$VALIDATE_ANSWER" == "y" ]; then
  echo "Validating the token..."
  validate_token "$JWT_TOKEN"
fi
