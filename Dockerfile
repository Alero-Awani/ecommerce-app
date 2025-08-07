# Use golang:alpine as base image for module caching
FROM golang:alpine AS modules

# Copy go files to a separate directory for module caching
COPY go.mod go.sum /modules/

# Set the working directory
WORKDIR /modules

# Install dependencies to cache them
RUN go mod download

# Use golang:alpine as base image for building the application
FROM golang:alpine AS builder

#Copy the cached modules from the previous stage
COPY --from=modules /go/pkg /go/pkg

# Copy source code to the builder stage
COPY . /app

# Set the working directory
WORKDIR /app

# Output binary is placed in /bin/app to separate it from the source code
RUN mkdir -p /bin && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /bin/app ./cmd/app

# Use alpine as the final image and copy the built binary
FROM alpine:latest
COPY --from=builder /bin/app /app

# Expose the port
EXPOSE 9000

# Start the application
ENTRYPOINT ["./app"]
