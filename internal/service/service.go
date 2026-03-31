package service

import (
	"strings"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/pkg/morse"
)

// DetectAndConvert определяет тип входных данных (текст или код Морзе)
// и конвертирует их в противоположный формат.
func DetectAndConvert(input string) (string, error) {
	// Удаляем пробелы в начале и конце
	input = strings.TrimSpace(input)

	if input == "" {
		return "", nil
	}

	// Проверяем, является ли строка кодом Морзе
	// Код Морзе состоит из: . - / пробелы и символы новой строки
	isMorse := true
	hasMorseChars := false
	for _, ch := range input {
		if ch == '.' || ch == '-' {
			hasMorseChars = true
		}
		if ch != '.' && ch != '-' && ch != ' ' && ch != '/' && ch != '\n' && ch != '\r' && ch != '\t' {
			isMorse = false
			break
		}
	}

	// Если это код Морзе (содержит . или - и нет других символов)
	if isMorse && hasMorseChars {
		result := morse.ToText(input)
		return result, nil
	}

	// Иначе конвертируем текст в код Морзе
	result := morse.ToMorse(input)
	return result, nil
}
