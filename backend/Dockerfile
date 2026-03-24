# Build stage
FROM --platform=linux/amd64 crpi-7tr092vgmwfahhs5.cn-shenzhen.personal.cr.aliyuncs.com/oct-home/golang:1.25-alpine AS builder

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

# Build migration tool
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /migrate ./scripts/migrate.go

# Production stage
FROM --platform=linux/amd64 crpi-7tr092vgmwfahhs5.cn-shenzhen.personal.cr.aliyuncs.com/oct-home/alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata wget mysql-client curl

# Copy binaries from builder
COPY --from=builder /baby-fans /baby-fans
COPY --from=builder /migrate /migrate

# Ensure binary is executable
RUN chmod +x /baby-fans /migrate

# Create necessary directories
RUN mkdir -p /app/storage/uploads /app/config /app/storage /app/certs

# Copy default config
COPY etc/config.yaml /app/config/config.yaml

# Expose port
EXPOSE 18081

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget -q --spider http://localhost:18081/health || exit 1

# Run migration then start the application
# MySQL will be tried first; if unavailable, falls back to SQLite
CMD ["/bin/sh", "-c", "/migrate && /baby-fans"]
