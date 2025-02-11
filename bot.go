package main

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

//const adminID int64 = ID // Замените на ваш Telegram ID

// StartBot запускает Telegram-бота
func StartBot(botToken string) {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	if err := LoadPoemsFromFile("poems.json"); err != nil {
		Logger.Printf("Ошибка загрузки стихов: %v", err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // игнорируем любые непрямо адресованные сообщения
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		// Обработка команд
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				msg.Text = "Привет! Отправь мне эмоцию 😊, 😞, 😂, ❤️, 🥔, 🩸 и я пришлю тебе стих! Чтобы ознакомиться со всем функционалом бота используйте команду /help"
				Logger.Printf("Пользователь %s (%d) выполнил команду /start", update.Message.From.UserName, update.Message.From.ID)
			case "random":
				msg.Text = GetRandomPoem()
				Logger.Printf("Пользователь %s (%d) выполнил команду /random", update.Message.From.UserName, update.Message.From.ID)
			case "addpoem":
				if update.Message.From.ID != adminID {
					msg.Text = "У вас нет прав для выполнения этой команды."
					Logger.Printf("Пользователь %s (%d) попытался выполнить команду /addpoem без прав администратора", update.Message.From.UserName, update.Message.From.ID)
					break
				}
				args := strings.SplitN(update.Message.CommandArguments(), " ", 2)
				if len(args) < 2 {
					msg.Text = "Использование: /addpoem <эмоция> <текст>"
					Logger.Printf("Администратор %s (%d) выполнил команду /addpoem с некорректными аргументами", update.Message.From.UserName, update.Message.From.ID)
					break
				}
				emotion := args[0]
				poemText := args[1]
				AddPoem(emotion, poemText)
				SavePoemsToFile("poems.json")
				msg.Text = "Стих успешно добавлен!"
				Logger.Printf("Администратор %s (%d) добавил стих для эмоции '%s': %s", update.Message.From.UserName, update.Message.From.ID, emotion, poemText)
			case "removepoem":
				if update.Message.From.ID != adminID {
					msg.Text = "У вас нет прав для выполнения этой команды."
					Logger.Printf("Пользователь %s (%d) попытался выполнить команду /removepoem без прав администратора", update.Message.From.UserName, update.Message.From.ID)
					break
				}
				args := strings.SplitN(update.Message.CommandArguments(), " ", 2)
				if len(args) < 2 {
					msg.Text = "Использование: /removepoem <эмоция> <текст>"
					Logger.Printf("Администратор %s (%d) выполнил команду /removepoem с некорректными аргументами", update.Message.From.UserName, update.Message.From.ID)
					break
				}
				emotion := args[0]
				poemText := args[1]
				removed := RemovePoem(emotion, poemText)
				if removed {
					SavePoemsToFile("poems.json")
					msg.Text = "Стих успешно удалён!"
					Logger.Printf("Администратор %s (%d) удалил стих для эмоции '%s': %s", update.Message.From.UserName, update.Message.From.ID, emotion, poemText)
				} else {
					msg.Text = "Стих не найден."
					Logger.Printf("Администратор %s (%d) попытался удалить несуществующий стих для эмоции '%s': %s", update.Message.From.UserName, update.Message.From.ID, emotion, poemText)
				}
			case "listpoems":
				args := strings.SplitN(update.Message.CommandArguments(), " ", 1)
				emotionFilter := ""
				if len(args) > 0 {
					emotionFilter = args[0]
				}
				poemsList := ListAllPoems(emotionFilter)
				if poemsList == "" {
					msg.Text = "Стихи не найдены."
					Logger.Printf("Пользователь %s (%d) выполнил команду /listpoems для эмоции '%s', но стихи не найдены", update.Message.From.UserName, update.Message.From.ID, emotionFilter)
				} else {
					parts := SplitLongMessage(poemsList)
					for _, part := range parts {
						msg.Text = part
						bot.Send(msg)
					}
					Logger.Printf("Пользователь %s (%d) выполнил команду /listpoems для эмоции '%s'", update.Message.From.UserName, update.Message.From.ID, emotionFilter)
				}
			case "help":
				msg.Text = `Вот список доступных команд:
				/start - Начать работу с ботом.
				/random - Получить случайный стих.
				/listpoems [эмоция] - Показать все стихи или только для указанной эмоции (например, /listpoems 😊).
				Отправьте эмоцию (например, 😊, 😞, 😂, ❤️, 🥔, 🩸) - Получить случайный стих для этой эмоции.
	
				Админские команды:
				/addpoem <эмоция> <текст> - Добавить новый стих (только для администратора).
				/removepoem <эмоция> <текст> - Удалить стих (только для администратора).`
				Logger.Printf("Пользователь %s (%d) выполнил команду /help", update.Message.From.UserName, update.Message.From.ID)
			default:
				msg.Text = "Неизвестная команда. Используйте /help для просмотра списка команд."
				Logger.Printf("Пользователь %s (%d) отправил неизвестную команду: %s", update.Message.From.UserName, update.Message.From.ID, update.Message.Text)
			}
			bot.Send(msg)
		}

		// Обработка эмодзи и других текстовых сообщений
		switch update.Message.Text {
		case "😊", "😞", "😂", "❤️", "🥔", "🩸":
			emotion := update.Message.Text
			Logger.Printf("Получена эмоция: %s", emotion)
			if poem, exists := GetRandomPoemByEmotion(emotion); exists {
				msg.Text = poem
				Logger.Printf("Пользователь %s (%d) получил стих для эмоции '%s'", update.Message.From.UserName, update.Message.From.ID, emotion)
			} else {
				msg.Text = "Я не знаю такой эмоции 😕\nПопробуй одну из этих: 😊, 😞, 😂, ❤️, 🥔, 🩸."
				Logger.Printf("Пользователь %s (%d) отправил неизвестную эмоцию: %s", update.Message.From.UserName, update.Message.From.ID, emotion)
			}
			bot.Send(msg)
		}
	}
}
