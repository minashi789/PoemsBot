package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Poem представляет структуру стиха
type Poem struct {
	Text string
}

// EmotionPoemsMap карта эмоций и соответствующих им стихов
var EmotionPoemsMap = make(map[string][]Poem)

// getRandomGenerator возвращает локальный генератор случайных чисел
func getRandomGenerator() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

// GetRandomPoemByEmotion возвращает случайный стих для заданной эмоции
func GetRandomPoemByEmotion(emotion string) (string, bool) {
	poems, exists := EmotionPoemsMap[emotion]
	if !exists || len(poems) == 0 {
		Logger.Printf("Нет стихов для эмоции '%s'", emotion)
		return "", false
	}
	randomIndex := getRandomGenerator().Intn(len(poems))
	Logger.Printf("Выбран случайный стих для эмоции '%s': %s", emotion, poems[randomIndex].Text)
	return poems[randomIndex].Text, true
}

// GetAllPoems возвращает все стихи из всех эмоций
func GetAllPoems() []Poem {
	var allPoems []Poem
	for _, poems := range EmotionPoemsMap {
		allPoems = append(allPoems, poems...)
	}
	return allPoems
}

// GetRandomPoem возвращает случайный стих из всех доступных
func GetRandomPoem() string {
	allPoems := GetAllPoems()
	if len(allPoems) == 0 {
		Logger.Printf("Нет доступных стихов для выбора случайного стиха.")
		return ""
	}
	randomIndex := getRandomGenerator().Intn(len(allPoems))
	Logger.Printf("Выбран случайный стих: %s", allPoems[randomIndex].Text)
	return allPoems[randomIndex].Text
}

// AddPoem добавляет новый стих в карту
func AddPoem(emotion, text string) {
	EmotionPoemsMap[emotion] = append(EmotionPoemsMap[emotion], Poem{Text: text})
}

// RemovePoem удаляет стих из карты
func RemovePoem(emotion, text string) bool {
	poems, exists := EmotionPoemsMap[emotion]
	if !exists || len(poems) == 0 {
		return false
	}

	// Ищем стих в списке
	for i, poem := range poems {
		if poem.Text == text {
			// Удаляем стих
			EmotionPoemsMap[emotion] = append(poems[:i], poems[i+1:]...)
			return true
		}
	}

	return false
}

// SavePoemsToFile сохраняет стихи в JSON-файл
func SavePoemsToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(EmotionPoemsMap)
}

func LoadPoemsFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			Logger.Printf("Файл %s не найден. Создаём пустую карту.", filename)
			EmotionPoemsMap = make(map[string][]Poem)
			return nil
		}
		Logger.Printf("Ошибка открытия файла %s: %v", filename, err)
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&EmotionPoemsMap); err != nil {
		Logger.Printf("Ошибка декодирования JSON из файла %s: %v", filename, err)
		return err
	}

	Logger.Printf("Загружено %d эмоций и %d стихов", len(EmotionPoemsMap), len(GetAllPoems()))

	Logger.Printf("Стихи успешно загружены из файла %s", filename)
	return nil

}

// SplitLongMessage разбивает длинный текст на части по 4096 символов
func SplitLongMessage(text string) []string {
	const maxMessageLength = 4096
	var parts []string

	for len(text) > maxMessageLength {
		parts = append(parts, text[:maxMessageLength])
		text = text[maxMessageLength:]
	}

	if len(text) > 0 {
		parts = append(parts, text)
	}

	return parts
}

// ListAllPoems возвращает все стихи с эмоциями в виде строки
func ListAllPoems(emotionFilter string) string {
	var result strings.Builder

	for emotion, poems := range EmotionPoemsMap {
		// Если задан фильтр по эмоции, пропускаем другие эмоции
		if emotionFilter != "" && emotion != emotionFilter {
			continue
		}

		if len(poems) == 0 {
			continue
		}

		// Добавляем эмоцию (жирным шрифтом)
		result.WriteString(fmt.Sprintf("**%s**\n", emotion))

		// Добавляем стихи
		for _, poem := range poems {
			result.WriteString(fmt.Sprintf("- %s\n", poem.Text))
		}

		// Разделитель между эмоциями
		result.WriteString("\n")
	}

	return result.String()
}
