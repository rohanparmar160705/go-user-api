# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod file
COPY go.mod ./
RUN go mod download
RUN go mod tidy

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/server/main.go

# Final stage
FROM alpine:latest

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/main .
# Note: In production you might want to pass env vars via docker-compose or k8s

# Install basic tools
RUN apk --no-cache add ca-certificates

EXPOSE 8080

CMD ["./main"]
