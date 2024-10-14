FROM node:18-alpine

# Set working directory
WORKDIR /usr/src/app

# Copy package.json and package-lock.json
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the src directory
COPY src ./src

# Copy any other necessary files (e.g., config files)
# COPY config ./config  # Uncomment if you have a config directory

# List contents of the work directory for debugging
RUN echo "Contents of /usr/src/app:" && ls -la && echo "Contents of /usr/src/app/src:" && ls -la src

# Expose the port the app runs on
EXPOSE 3000

# Command to run the application
CMD ["node", "src/auth-service.js"]