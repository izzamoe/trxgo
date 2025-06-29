# Multi-stage build untuk efisiensi
FROM golang:1.24-alpine AS builder

# Install dependencies yang diperlukan
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
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/server/main.go

# Final stage - minimal image
FROM alpine:latest

# Install ca-certificates untuk HTTPS calls
RUN apk --no-cache add ca-certificates mysql-client

WORKDIR /root/

# Copy binary dari builder stage
COPY --from=builder /app/main .

# Copy .env file jika ada
COPY --from=builder /app/.env* ./

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the binary
CMD ["./main"]
