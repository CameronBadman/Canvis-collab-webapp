# Build stage
FROM node:18-alpine as build
WORKDIR /app

# Copy package files and install dependencies
COPY package*.json ./
RUN npm install

# Copy the entire app
COPY . .

# Set environment variables for the backend API URL
ARG BACKEND
ENV BACKEND=${BACKEND:-http://nginx-api-gateway:8000}

# Debugging: Print out the URL to verify it's set correctly
RUN echo "BACKEND: $BACKEND"

# Build the app
RUN npm run build

# List the contents of the dist folder to ensure build was successful
RUN ls -la /app/dist

# Production stage
FROM nginx:stable-alpine

# Copy the built frontend files to Nginx's serving directory
COPY --from=build /app/dist /usr/share/nginx/html

# Copy the custom Nginx configuration (if needed for routing, etc.)
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose port 3002 (configured for the frontend)
EXPOSE 3002

# Start Nginx in the foreground
CMD ["nginx", "-g", "daemon off;"]
