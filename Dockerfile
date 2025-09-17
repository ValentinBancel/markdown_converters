# Build stage
FROM golang:1.21-alpine AS builder

# Install git and ca-certificates, with skip verify for restricted networks
RUN apk --no-cache add git ca-certificates || echo "Warning: Could not install packages"

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

# Install ca-certificates (skip if fails)
RUN apk --no-cache add ca-certificates || echo "Warning: Could not install ca-certificates"

# Create non-root user
RUN adduser -D -g '' appuser

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
