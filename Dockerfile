# Build stage
FROM golang:1.21-alpine AS builder

# Install git and ca-certificates
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files (go.sum may not exist for projects without external deps)
COPY go.mod ./
COPY go.sum* ./

# Download dependencies (if go.mod exists)
RUN if [ -f go.mod ]; then go mod download || true; fi

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1000 -S appgroup && \
    adduser -u 1000 -S appuser -G appgroup

# Set working directory
WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/main .

# Copy templates
COPY --from=builder /app/templates ./templates

# Change ownership
RUN chown -R appuser:appgroup /root/

# Use non-root user
USER appuser

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]
