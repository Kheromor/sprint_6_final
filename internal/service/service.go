package service

import (
	"strings"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/pkg/morse"
)

func DetectAndConvert(data string) (string, error) {
	data = strings.TrimSpace(data)
	if data == "" {
		return "", nil
	}

	for _, char := range data {
		if (char >= 'а' && char <= 'я') || (char >= 'А' && char <= 'Я') {
			return morse.ToMorse(data), nil
		}
	}

	normalized := strings.Join(strings.Fields(data), " ")
	return morse.ToText(normalized), nil
}
