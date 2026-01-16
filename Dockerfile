# Build Stage
FROM golang:1.24.4-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Download generic dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
# We are building the new 'app' entrypoint
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/main ./cmd/app/main.go

# Final Stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies (certificates for HTTPS if needed)
RUN apk --no-cache add ca-certificates

# Copy binary from builder
COPY --from=builder /app/main .

# Copy static assets and templates
COPY --from=builder /app/web ./web
COPY --from=builder /app/assets ./assets
COPY --from=builder /app/.env.example .env

# Expose port
EXPOSE 8080

# Run
CMD ["/app/main"]
