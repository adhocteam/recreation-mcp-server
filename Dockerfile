# Build stage
FROM golang:1.25.12-alpine3.22 AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Copy dependency files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary with security flags
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o mcp-server ./cmd/server

# Runtime stage
FROM alpine:3.22.2

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata && \
    rm -rf /var/cache/apk/*

# Create non-root user
RUN addgroup -g 1000 mcp && \
    adduser -D -u 1000 -G mcp mcp

# Set up working directory
WORKDIR /home/mcp

# Copy binary from builder with correct ownership
COPY --from=builder --chown=mcp:mcp /app/mcp-server .

# Switch to non-root user
USER mcp

# Run the MCP server
CMD ["./mcp-server"]
