# ---------- Build stage ----------
FROM golang:1.23.4-alpine AS builder

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o chat-server ./cmd/server

# ---------- Runtime stage ----------
FROM alpine:3.19

WORKDIR /app

# Certificates for HTTPS, DB, etc.
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/chat-server .

EXPOSE 8080

CMD ["./chat-server"]
