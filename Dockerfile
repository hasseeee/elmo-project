# Build stage
FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o main \
    ./cmd/server/main.go

# Final stage
FROM --platform=$TARGETPLATFORM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Create non-root user
RUN addgroup -g 1001 -S elmo && \
    adduser -u 1001 -S elmo -G elmo

# Copy the binary from builder
COPY --from=builder /app/main .

# Change ownership to non-root user
RUN chown -R elmo:elmo /app

# Switch to non-root user
USER elmo

# Expose port 8080
EXPOSE 8080


# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Command to run the application
CMD ["./main"]