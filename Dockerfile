# Базовый образ с Go
FROM golang:1.20

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./

# Устанавливаем зависимости
RUN go mod download

# Копируем остальные файлы
COPY . .

# Собираем проект
RUN go build -o bot .

# Команда для запуска бота
CMD ["./bot"]