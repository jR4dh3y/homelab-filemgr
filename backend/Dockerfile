# Backend Dockerfile for Homelab File Manager
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /server ./cmd/server

# Final stage
FROM alpine:3.20

WORKDIR /app

# Install ca-certificates for HTTPS and tzdata for timezone support
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user for security (but we'll run as root for file access)
# The container needs root to access mounted volumes with various permissions
RUN addgroup -g 1000 filemanager && \
    adduser -u 1000 -G filemanager -s /bin/sh -D filemanager

# Copy binary from builder
COPY --from=builder /server /app/server

# Copy default config
COPY config.yaml /app/config.yaml

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/v1/health || exit 1

# Run as root to access mounted filesystems with proper permissions
# In production, consider using user namespaces or specific UID/GID mapping
ENTRYPOINT ["/app/server"]
