# Multi-stage build for OP Vault Local Edition
# Stage 1: Build frontend and backend
FROM node:19-alpine AS builder

# Install Go and build dependencies
RUN apk add --no-cache \
    go \
    git \
    build-base \
    poppler-utils

# Set Go environment
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH
ENV CGO_ENABLED=1

WORKDIR /app

# Copy package files
COPY package*.json ./
COPY go.mod go.sum ./

# Download dependencies
RUN npm ci --only=production && \
    npm install webpack webpack-cli --save-dev && \
    go mod download

# Copy source code
COPY . .

# Build frontend
RUN npm run build || webpack --mode production

# Build backend
RUN cd vault-web-server && \
    go build -o ../bin/vault-web-server .

# Stage 2: Runtime image
FROM node:19-alpine

# Install runtime dependencies
RUN apk add --no-cache \
    poppler-utils \
    ca-certificates

WORKDIR /app

# Copy built artifacts from builder
COPY --from=builder /app/bin/vault-web-server /app/bin/vault-web-server
COPY --from=builder /app/static /app/static
COPY --from=builder /app/web /app/web
COPY --from=builder /app/serverutil /app/serverutil

# Create directory for Qdrant data
RUN mkdir -p /app/data

# Expose port
EXPOSE 8100

# Set environment variables with defaults
ENV OLLAMA_ENDPOINT=http://ollama:11434
ENV OLLAMA_EMBED_MODEL=nomic-embed-text
ENV OLLAMA_CHAT_MODEL=llama3
ENV QDRANT_API_ENDPOINT=http://qdrant:6333

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8100/ || exit 1

# Run the application
CMD ["/app/bin/vault-web-server"]
