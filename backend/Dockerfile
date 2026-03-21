# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /baby-fans ./cmd/server

# Production stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata wget

# Copy binary from builder
COPY --from=builder /baby-fans /baby-fans

# Ensure binary is executable
RUN chmod +x /baby-fans

# Create necessary directories
RUN mkdir -p /app/storage/uploads /app/config

# Copy default config if etc directory is mounted
# Note: config should be mounted at /app/config

# Expose port
EXPOSE 18081

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:18081/health || exit 1

# Run the application
CMD ["/bin/sh", "-c", "/baby-fans"]
