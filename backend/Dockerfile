# Multi-stage Dockerfile for Golang Backend

# Stage 1: Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Stage 2: Development stage (with hot reload)
FROM golang:1.23-alpine AS development

# Install air for hot reload
RUN go install github.com/cosmtrek/air@v1.49.0

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

EXPOSE 8080

# Will be overridden by docker-compose command
CMD ["air", "-c", ".air.toml"]

# Stage 3: Production stage
FROM alpine:latest AS production

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Set timezone to Africa/Addis_Ababa
ENV TZ=Africa/Addis_Ababa

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy config files if any
COPY --from=builder /app/config ./config

EXPOSE 8080

CMD ["./main"]
