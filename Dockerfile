# Стадия сборки
FROM golang:1.20-alpine AS builder

WORKDIR /app

# Копируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем код и файл poems.json
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o bot .

# Финальный образ
FROM alpine:latest

WORKDIR /root/

# Копируем исполняемый файл и poems.json
COPY --from=builder /app/bot .
COPY --from=builder /app/poems.json .

# Команда запуска
CMD ["./bot"]