# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies for SQLite
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary with CGO enabled for SQLite
RUN CGO_ENABLED=1 GOOS=linux go build -a -o jason-frontpage ./cmd/server

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates sqlite-libs

WORKDIR /root/

# Create workspace directory for persistent data
RUN mkdir -p /workspace/data

# Copy binary
COPY --from=builder /app/jason-frontpage .

# Copy web assets
COPY --from=builder /app/web ./web

EXPOSE 8080

CMD ["./jason-frontpage"]
