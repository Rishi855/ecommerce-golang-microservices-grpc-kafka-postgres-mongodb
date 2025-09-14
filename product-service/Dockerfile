# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Install necessary tools
RUN apk --no-cache add netcat-openbsd

# Create necessary directories
RUN mkdir -p internal/config

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy the config directory
COPY --from=builder /app/internal/config/.env ./internal/config/

# Command to run the executable
CMD ["./main"]
