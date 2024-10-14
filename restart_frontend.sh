#!/bin/bash

# Name of the frontend container
FRONTEND_CONTAINER="canvis-collab-webapp-frontend-1"

# Restart the frontend container
echo "Restarting $FRONTEND_CONTAINER..."
docker restart $FRONTEND_CONTAINER

# Check if the restart was successful
if [ $? -eq 0 ]; then
    echo "$FRONTEND_CONTAINER restarted successfully."
else
    echo "Failed to restart $FRONTEND_CONTAINER."
fi
