FROM golang:1.22.5-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o ads-service cmd/main.go

# Создаем второй этап для минимального runtime
FROM alpine:3.20

WORKDIR /app

# Копируем папку config, migrations и исполняемый файл 
COPY --from=builder /app/ads-service .
COPY --from=builder /app/config ./config
COPY --from=builder /app/migrations ./migrations

CMD ["./ads-service"]