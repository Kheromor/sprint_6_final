package service

import (
	"regexp"
	"strings"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/pkg/morse"
)

func DetectAndConvert(data string) (string, error) {
	data = strings.TrimSpace(data)
	if data == "" {
		return "", nil
	}

	// Проверяем наличие русских букв
	for _, char := range data {
		if (char >= 'а' && char <= 'я') || (char >= 'А' && char <= 'Я') {
			return morse.ToMorse(data), nil
		}
	}

	wordRegex := regexp.MustCompile(`\s{2,}`)
	words := wordRegex.Split(data, -1)

	var normalizedWords []string
	for _, word := range words {
		// Убираем лишние пробелы внутри слова
		symbols := strings.Fields(word)
		if len(symbols) > 0 {
			normalizedWords = append(normalizedWords, strings.Join(symbols, " "))
		}
	}

	// Собираем с тремя пробелами между словами (требование пакета morse)
	normalizedMorse := strings.Join(normalizedWords, "   ")

	// Конвертируем в текст
	return morse.ToText(normalizedMorse), nil
}
