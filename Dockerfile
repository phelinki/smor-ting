# Root Dockerfile for Railway: build Go backend from smor_ting_backend

FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata
WORKDIR /app

# Copy only module files first for better caching
COPY smor_ting_backend/go.mod smor_ting_backend/go.sum ./
RUN go mod download

# Copy backend source
COPY smor_ting_backend/. .

# Build API binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-s -w" -o smor-ting-api ./cmd

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app

# Non-root user
RUN addgroup -g 1001 -S smor-ting && adduser -u 1001 -S smor-ting -G smor-ting
USER smor-ting

# Copy binary
COPY --from=builder /app/smor-ting-api ./smor-ting-api

EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./smor-ting-api"]


# Railway deployment trigger Wed Aug 13 00:22:46 CDT 2025
