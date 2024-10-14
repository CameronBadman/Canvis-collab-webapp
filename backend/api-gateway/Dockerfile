# Use the official Go image as the base image
FROM golang:1.22 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./
COPY go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o api-gateway ./cmd

# Start a new stage from scratch
FROM gcr.io/distroless/base

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/api-gateway .

# Command to run the executable
ENTRYPOINT ["/api-gateway"]

# Expose port 8080 to the outside world
EXPOSE 8080
