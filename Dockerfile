FROM golang:1.20
WORKDIR /app
COPY . .
RUN go build -o bot main.go bot.go poems.go logger.go
CMD ["./bot"]