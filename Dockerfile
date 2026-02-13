# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install git for go mod download
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o jason-frontpage ./cmd/server

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Create workspace directory
RUN mkdir -p /workspace

# Copy binary
COPY --from=builder /app/jason-frontpage .

# Copy web assets
COPY --from=builder /app/web ./web

EXPOSE 8080

CMD ["./jason-frontpage"]