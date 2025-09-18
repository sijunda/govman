# Dockerfile for GOVMAN
# Multi-stage build for minimal image size

# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -X github.com/sijunda/govman/internal/version.Version=docker" \
    -a -installsuffix cgo \
    -o govman \
    ./cmd/govman

# Final stage - minimal image
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add \
    ca-certificates \
    git \
    curl \
    tar \
    gzip \
    && adduser -D -s /bin/sh govman

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy CA certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy binary
COPY --from=builder /app/govman /usr/local/bin/govman

# Create govman directories
RUN mkdir -p /home/govman/.govman/{bin,cache,versions} && \
    chown -R govman:govman /home/govman/.govman

# Switch to non-root user
USER govman
WORKDIR /home/govman

# Set default environment
ENV HOME=/home/govman
ENV PATH="/home/govman/.govman/bin:${PATH}"

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD govman version || exit 1

# Default command
ENTRYPOINT ["govman"]
CMD ["--help"]

# Labels for metadata
LABEL org.opencontainers.image.title="GOVMAN - Go Version Manager"
LABEL org.opencontainers.image.description="Cross-platform Go version manager"
LABEL org.opencontainers.image.url="https://github.com/sijunda/govman"
LABEL org.opencontainers.image.source="https://github.com/sijunda/govman"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.vendor="sijunda"