# ================================
# Stage 1: Build all binaries
# ================================
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copy go mod files first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build all binaries
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /bin/api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /bin/cron ./cmd/cron
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /bin/worker ./cmd/worker

# ================================
# Stage 2: API Service
# ================================
FROM alpine:3.19 AS api

RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app

# Copy binary
COPY --from=builder /bin/api /app/api

# Copy .env if needed for local testing (use env vars in production)
# COPY .env .env

EXPOSE 8080

CMD ["/app/api"]

# ================================
# Stage 3: Cron Service
# ================================
FROM alpine:3.19 AS cron

RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app

COPY --from=builder /bin/cron /app/cron

CMD ["/app/cron"]

# ================================
# Stage 4: Worker Service
# ================================
FROM alpine:3.19 AS worker

RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app

COPY --from=builder /bin/worker /app/worker

CMD ["/app/worker"]
