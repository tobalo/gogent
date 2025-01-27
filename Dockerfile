# Build stage
FROM golang:1.23.4-bullseye AS builder

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apt-get update && apt-get install -y \
    git \
    gcc \
    g++ \
    sqlite3 \
    libsqlite3-dev \
    && rm -rf /var/lib/apt/lists/*

# Set environment variables
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=arm64

# Copy go mod files first for better layer caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary with explicit architecture flags
RUN CGO_CFLAGS="-O2" CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -o /app/gogent ./cmd/microlith

# Final stage
FROM debian:bullseye-slim

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    sqlite3 \
    curl \
    netcat \
    && rm -rf /var/lib/apt/lists/*

# Create non-root user
RUN useradd -r -s /bin/false appuser

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/gogent .

# Copy .env file if it exists (optional at build time)
COPY .env* ./

# Set ownership to non-root user
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose NATS port
EXPOSE 4222

# Environment variables for configuration
ENV API_KEY=""
ENV PROVIDER="OLLAMA"
ENV MODEL="deepseek-r1:1.5b"

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD nc -zv localhost 4222 || exit 1

# Run the binary
ENTRYPOINT ["./gogent"]