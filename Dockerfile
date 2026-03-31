FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Устанавливаем goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Собираем приложение
RUN CGO_ENABLED=0 go build -o server ./cmd/server

# Копируем бинарник goose
RUN cp /go/bin/goose /goose

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /goose .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./server"]
