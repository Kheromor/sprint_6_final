package service

import (
	"strings"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/pkg/morse"
)

// DetectAndConvert определяет тип входных данных (текст или код Морзе)
// и конвертирует их в противоположный формат
func DetectAndConvert(data string) (string, error) {
	// Удаляем пробелы в начале и конце строки для анализа
	trimmedData := strings.TrimSpace(data)
	
	if trimmedData == "" {
		return "", nil
	}
	
	// Проверяем, является ли строка кодом Морзе
	if isMorseCode(trimmedData) {
		// Конвертируем из Морзе в текст
		result := morse.ToText(data)
		return result, nil
	}
	
	// Иначе считаем, что это текст, и конвертируем в Морзе
	result := morse.ToMorse(data)
	return result, nil
}

// isMorseCode проверяет, является ли строка кодом Морзе
func isMorseCode(s string) bool {
    // Если в строке есть русские буквы - это текст
    for _, char := range s {
        if (char >= 'а' && char <= 'я') || (char >= 'А' && char <= 'Я') {
            return false
        }
    }
    return true
}