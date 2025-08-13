# Multi-stage build for production
FROM golang:1.23-alpine AS builder

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

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build \
    -a -installsuffix cgo \
    -ldflags="-s -w" \
    -o smor-ting-api \
    cmd/main.go

# Production stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S smor-ting && \
    adduser -u 1001 -S smor-ting -G smor-ting

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/smor-ting-api .

# Create necessary directories
RUN mkdir -p /app/logs /app/backups /app/static && \
    chown -R smor-ting:smor-ting /app

# Switch to non-root user
USER smor-ting

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
ENTRYPOINT ["./smor-ting-api"]# Last updated: Wed Aug 13 00:00:25 CDT 2025
