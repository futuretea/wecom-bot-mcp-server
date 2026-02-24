# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /build

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary with version info
ARG VERSION=dev
ARG GIT_COMMIT=unknown
ARG BUILD_DATE=unknown
RUN CGO_ENABLED=0 go build \
    -ldflags "-s -w \
      -X 'github.com/futuretea/wecom-bot-mcp-server/pkg/core/version.Version=${VERSION}' \
      -X 'github.com/futuretea/wecom-bot-mcp-server/pkg/core/version.GitCommit=${GIT_COMMIT}' \
      -X 'github.com/futuretea/wecom-bot-mcp-server/pkg/core/version.BuildDate=${BUILD_DATE}'" \
    -o wecom-bot-mcp-server ./cmd/wecom-bot-mcp-server

# Final stage
FROM cgr.dev/chainguard/wolfi-base:latest AS runtime

# Create non-root user
RUN adduser -D -s /bin/sh wecom

USER wecom

ENTRYPOINT ["/usr/local/bin/wecom-bot-mcp-server"]

# Release image
FROM runtime AS release

COPY wecom-bot-mcp-server /usr/local/bin/wecom-bot-mcp-server

# Dev image
FROM runtime AS dev

# Copy the binary from builder
COPY --from=builder /build/wecom-bot-mcp-server /usr/local/bin/wecom-bot-mcp-server
