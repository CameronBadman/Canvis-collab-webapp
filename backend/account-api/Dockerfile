# Use a small Go image to build the application
FROM golang:1.18 as builder

# Set the working directory in the container
WORKDIR /app

# Copy go.mod and go.sum files and install dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the source code
COPY . .

# Build the Go application
RUN go build -o account-api .

# Use a minimal image to run the built application
FROM gcr.io/distroless/base-debian10

# Copy the built binary from the builder stage
COPY --from=builder /app/account-api /account-api

# Set the entrypoint to run the binary
ENTRYPOINT ["/account-api"]
