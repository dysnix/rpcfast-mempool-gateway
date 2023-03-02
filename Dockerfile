# First stage: Build the app using Golang base image
FROM golang:1.19 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the binary
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o rpcfast-mempool-gateway

# Second stage: Create a minimal image using a scratch base image
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/rpcfast-mempool-gateway /app

# Expose the port that the app listens on
EXPOSE 8080

# Start the app
CMD ["/app/rpcfast-mempool-gateway"]
