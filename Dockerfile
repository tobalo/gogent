# Build stage
FROM golang:1.23.4-alpine3.19 AS builder

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files first for better layer caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/gogent ./cmd/microlith

# Final stage
FROM alpine:3.19

# Install CA certificates for HTTPS connections
RUN apk add --no-cache ca-certificates && \
    update-ca-certificates

# Create non-root user
RUN adduser -D -g '' appuser

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

# Environment variable for Gemini API key
ENV GEMINI_API_KEY=""

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD nc -zv localhost 4222 || exit 1

# Run the binary
ENTRYPOINT ["./gogent"]