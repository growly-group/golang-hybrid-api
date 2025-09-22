# Multi-stage build for smaller final image
FROM golang:1.24.2-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o main .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Default to running calculator service only
# Override with different TARGET_SERVICES in docker-compose or kubernetes
ENV TARGET_SERVICES=calculator-svc

# Expose default port (calculator service)
EXPOSE 8080

CMD ["./main"]