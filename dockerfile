# Development stage
FROM golang:1.24-alpine AS development

# Add git for go mod download
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
# Use cache for go mod download
RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  go mod download

# Copy the rest of the application
COPY . .

# Install air for hot reloading in development
RUN go install github.com/air-verse/air@v1.62.0

# Expose the application port
EXPOSE 8000

# Use air for hot reloading in development
CMD ["air", "-c", ".air.toml"]

# Production build stage
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the application
# Reduce the size of the binary by stripping debug symbols "-ldflags=\"-s -w\""
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -trimpath -mod=readonly -o server ./cmd/server

# Production stage
FROM alpine:3.19 AS production

# Combine runtime dependency installation and user creation
RUN apk --no-cache add ca-certificates tzdata wget && adduser -D appuser

WORKDIR /app

# Copy the binary from the builder stage and set ownership
COPY --from=builder --chown=appuser:appuser /app/server .
USER appuser

EXPOSE 8000

# Use exec form for HEALTHCHECK
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s \
  CMD ["wget", "-qO-", "http://localhost:8000/health"]

CMD ["./server"]
