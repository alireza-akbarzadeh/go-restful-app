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
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o api-server ./cmd/server

# Build migration tool
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o migrate-tool ./cmd/migrate

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates

# Create a non-root user
RUN adduser -D -g '' appuser

# Create app directory
WORKDIR /app

# Create data directory for database and set permissions
RUN mkdir -p /app/data && chown -R appuser:appuser /app

# Copy binaries from builder
COPY --from=builder /app/api-server .
COPY --from=builder /app/migrate-tool .

# Copy migrations
COPY --from=builder /app/cmd/migrate/migrations ./cmd/migrate/migrations

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
ENV DATABASE_PATH=/app/data/data.db
ENV APP_ENV=production

# Run migrations and start the application
CMD ["sh", "-c", "./migrate-tool up && ./api-server"]
