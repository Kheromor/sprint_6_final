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

	words := regexp.MustCompile(`\s{2,}`).Split(data, -1)

	var normalizedWords []string
	for _, word := range words {
		// Нормализуем пробелы внутри слова
		symbols := strings.Fields(word)
		if len(symbols) > 0 {
			normalizedWords = append(normalizedWords, strings.Join(symbols, " "))
		}
	}

	// Собираем с тремя пробелами между словами
	normalizedMorse := strings.Join(normalizedWords, "   ")

	return morse.ToText(normalizedMorse), nil
}
