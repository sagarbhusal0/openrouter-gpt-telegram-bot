# ---- Build Stage ----
FROM golang:1.23.4 AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go binary
RUN go build -o openrouter-gpt-telegram-bot

# ---- Run Stage ----
FROM debian:bookworm-slim

WORKDIR /app

# Copy the binary from the builder
COPY --from=builder /app/openrouter-gpt-telegram-bot /app/openrouter-gpt-telegram-bot

# Copy config.yaml and other runtime files (if present)
COPY config.yaml ./config.yaml
COPY lang ./lang

# Expose port if your app uses one (optional)
# EXPOSE 8080

CMD ["/app/openrouter-gpt-telegram-bot"]