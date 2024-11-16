# Use a small Go image for building
FROM golang:1.23 as builder

# Set the working directory in the container
WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod tidy && go mod download

# Copy the application source code
COPY . .

# Build the Go application
RUN go build -o account-api .

# Final stage: Use a runtime image with the correct glibc version
FROM debian:bookworm-slim

# Install necessary runtime dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Set up working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/account-api /app/account-api

# Ensure the binary is executable
RUN chmod +x /app/account-api

# Expose the application's port
EXPOSE 8080

# Set the entrypoint to run the binary
ENTRYPOINT ["/app/account-api"]
