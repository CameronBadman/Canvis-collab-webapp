# Start with the official Go image to build the application
FROM golang:1.23-alpine as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go Modules manifests to the container
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod tidy

# Copy the source code into the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o canvas-api .

# Start a new stage to keep the image lean
FROM alpine:latest  

# Install necessary dependencies for running the Go app (Alpine has minimal libraries)
RUN apk --no-cache add ca-certificates

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/canvas-api .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./canvas-api"]
