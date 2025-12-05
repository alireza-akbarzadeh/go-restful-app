# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
# CGO_ENABLED=0 is preferred for pure Go binaries (standard for Postgres/pgx)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api-server ./cmd/server

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates

# Create a non-root user
RUN adduser -D -g '' appuser

# Create app directory
WORKDIR /app

# Copy binaries from builder
COPY --from=builder /app/api-server .

# Copy swagger docs (if needed)
COPY --from=builder /app/docs ./docs

# Set ownership of all copied files to appuser
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Set environment variables
ENV PORT=8080
ENV APP_ENV=production

# Start the application
CMD ["./api-server"]
