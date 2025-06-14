# syntax=docker/dockerfile:1.4

# Build stage
FROM golang:1.24.4-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o templater ./cmd/templater

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN adduser -D -g '' templater

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/templater /app/templater

# Copy templates and config
COPY --from=builder /app/.templater.yaml /app/.templater.yaml

# Set ownership
RUN chown -R templater:templater /app

# Switch to non-root user
USER templater

# Set environment variables
ENV PATH="/app:${PATH}"

# Health check
HEALTHCHECK --interval=30s --timeout=5s \
    CMD templater --version || exit 1

# Default command
ENTRYPOINT ["/app/templater"] 