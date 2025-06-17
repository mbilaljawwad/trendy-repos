# Multi-stage Dockerfile for Go application with hot reloading

# Development stage with hot reloading
FROM golang:1.24-alpine AS development

# Install git and other dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Install Air for hot reloading
RUN go install github.com/air-verse/air@latest

# Set working directory
WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Create air config if it doesn't exist
RUN if [ ! -f .air.toml ]; then air init; fi

# Expose port
EXPOSE 8080

# Command for development with hot reloading
CMD ["air", "-c", ".air.toml"]

# Production build stage
FROM golang:1.23-alpine AS builder

# Install dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Production stage
FROM alpine:latest AS production

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy any config files if they exist
COPY --from=builder /app/*.env* ./

# Change ownership to non-root user
RUN chown -R appuser:appgroup /root

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Command to run the application
CMD ["./main"] 