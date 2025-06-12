FROM golang:1.23.4 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Ensure build for Linux and disable CGO for static binary
ENV CGO_ENABLED=0
RUN GOOS=linux GOARCH=amd64 go build -o openrouter-gpt-telegram-bot

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/openrouter-gpt-telegram-bot /app/openrouter-gpt-telegram-bot
COPY config.yaml ./config.yaml
COPY lang ./lang
RUN chmod +x /app/openrouter-gpt-telegram-bot
CMD ["/app/openrouter-gpt-telegram-bot"]