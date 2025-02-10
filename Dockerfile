# Стадия 1: Сборка приложения
FROM golang:1.20-alpine AS builder

WORKDIR /app

# Устанавливаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем код и собираем проект
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bot .

# Стадия 2: Финальный образ
FROM alpine:latest

WORKDIR /root/

# Копируем только исполняемый файл
COPY --from=builder /app/bot .

# Команда для запуска бота
CMD ["./bot"]