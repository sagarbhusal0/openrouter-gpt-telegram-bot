FROM golang:1.23.4 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o openrouter-gpt-telegram-bot

FROM debian:bookworm-slim
WORKDIR /app
# Install CA certificates
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
# Create logs directory and set permissions
RUN mkdir -p /app/logs && chmod 777 /app/logs
COPY --from=builder /app/openrouter-gpt-telegram-bot /app/
COPY config.yaml ./config.yaml
COPY lang ./lang
CMD ["/app/openrouter-gpt-telegram-bot"]