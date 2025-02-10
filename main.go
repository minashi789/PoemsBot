package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	if err := InitLogger("bot.log"); err != nil {
		log.Fatalf("Ошибка инициализации логгера: %v", err)
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Ошибка загрузки .env файла: %v", err)
	}

	//const botToken = "7597451477:AAFV6FNLNiSmQHzwmnySTRvP4-fNXidI-IQ"

	LoadPoemsFromFile("poems.json")

	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatalf("Токен бота не найден. Проверьте переменную окружения BOT_TOKEN.")
	}
	log.Printf("Токен бота: %s", botToken)

	StartBot(botToken)
}
