# Development stage
FROM golang:1.24-alpine AS development

# Add git for go mod download
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application
COPY . .

# Install air for hot reloading in development
RUN go install github.com/air-verse/air@latest

# Expose the application port
EXPOSE 8000

# Use air for hot reloading in development
CMD ["air", "-c", ".air.toml"]

# Production build stage
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Production stage
FROM alpine:latest AS production

RUN apk --no-cache add ca-certificates tzdata
RUN adduser -D appuser

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/server .

USER appuser

EXPOSE 8000

CMD ["./server"]
