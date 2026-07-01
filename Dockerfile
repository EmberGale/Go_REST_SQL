FROM golang:alpine AS builder

WORKDIR /app

RUN apk update && apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

FROM alpine

WORKDIR /app

RUN adduser -D -g '' appuser && \
    apk add --no-cache ca-certificates

COPY --from=builder /app/server /app/server
COPY --from=builder /app/migrations /app/migrations

ENV APP_DATABASE__MIGRATIONS_PATH=./migrations

EXPOSE 8080

USER appuser

CMD ["/app/server"]