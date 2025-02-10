package main

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var adminID int64

func main() {

	if err := InitLogger("bot.log"); err != nil {
		log.Fatalf("Ошибка инициализации логгера: %v", err)
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Ошибка загрузки .env файла: %v", err)
	}

	LoadPoemsFromFile("poems.json")

	adminIDStr := os.Getenv("ADMIN_ID")
	var err error
	adminID, err = strconv.ParseInt(adminIDStr, 10, 64)
	if err != nil {
		log.Fatalf("Неверный формат ADMIN_ID: %v", err)
	}

	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatalf("Токен бота не найден. Проверьте переменную окружения BOT_TOKEN.")
	}
	log.Printf("Токен бота: %s", botToken)

	StartBot(botToken)
}
