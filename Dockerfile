# Use Go 1.23 bookworm as base image
FROM golang:alpine AS builder

RUN apk update && \
apk add git ca-certificates

# Move to working directory /app
WORKDIR /app

# Copy the go.mod and go.sum files to the /app directory
COPY go.mod go.sum ./

# Install dependencies
RUN go mod download

# Copy the entire source code into the container
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o ecommerce

# Create a production stage to run the application binary
FROM alpine:latest AS production

# Install ca-certificates for HTTPS requests (if needed)
RUN apk --no-cache add ca-certificates

# Move to working directory /prod
WORKDIR /prod

# Copy binary from builder stage
COPY --from=builder /app ./

# Publish the port
EXPOSE 9000

# Start the application
ENTRYPOINT ["./ecommerce"]