package main

import (
	"log"
	"os"
)

var Logger *log.Logger

// InitLogger инициализирует логгер для записи в файл
func InitLogger(filename string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	Logger = log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)
	return nil
}
