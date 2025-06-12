FROM golang:1.24-alpine as development

# Install air for hot reloading in development (updated repository)
RUN go install github.com/air-verse/air@latest

# Install basic tools
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Expose port
EXPOSE 8080

# Default command - will be overridden by docker-compose
CMD ["go", "run", "./server/main.go"]
