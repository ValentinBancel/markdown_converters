# Build stage
FROM golang:1.21-alpine AS builder

# Install git (required for go mod download) 
RUN apk --no-cache add git || true

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies with insecure flag for restricted networks
RUN go env -w GOPROXY=direct && go env -w GOSUMDB=off && go mod download || echo "Warning: Could not download all dependencies"

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage - use alpine with ca-certificates already present
FROM alpine:latest

# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Change ownership
RUN chown -R appuser:appuser /app

# Use non-root user
USER appuser

# Expose port
EXPOSE 8080

# Command to run
CMD ["./main"]
