package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/service"
)

// IndexHandler возвращает HTML форму из файла index.html
func IndexHandler(logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		// Пытаемся найти index.html
		var htmlContent []byte
		var err error

		paths := []string{
			"index.html",
			"../index.html",
			"./index.html",
			"../../index.html",
		}

		for _, path := range paths {
			htmlContent, err = os.ReadFile(path)
			if err == nil {
				logger.Printf("Найден index.html по пути: %s", path)
				break
			}
		}

		if err != nil {
			logger.Printf("Не удалось найти index.html, использую встроенную форму")
			htmlContent = []byte(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Преобразователь азбуки Морзе</title>
</head>
<body>
    <h1>Преобразователь азбуки Морзе</h1>
    <form action="/upload" method="post" enctype="multipart/form-data">
        <input type="file" name="myFile" required>
        <br><br>
        <button type="submit">Перевести</button>
    </form>
</body>
</html>`)
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(htmlContent)
	}
}

// UploadHandler обрабатывает загрузку файлов и их конвертацию
func UploadHandler(logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
			return
		}

		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			logger.Printf("Ошибка парсинга формы: %v", err)
			http.Error(w, "Не удаётся разобрать форму", http.StatusBadRequest)
			return
		}

		// Используем имя поля "myFile" из index.html
		file, fileHeader, err := r.FormFile("myFile")
		if err != nil {
			logger.Printf("Ошибка получения файла: %v", err)
			http.Error(w, "Не удается получить файл", http.StatusBadRequest)
			return
		}
		defer file.Close()

		logger.Printf("Получен файл: %s", fileHeader.Filename)

		fileContent, err := io.ReadAll(file)
		if err != nil {
			logger.Printf("Ошибка чтения файла: %v", err)
			http.Error(w, "Не удается прочитать файл", http.StatusInternalServerError)
			return
		}

		convertedContent, err := service.DetectAndConvert(string(fileContent))
		if err != nil {
			logger.Printf("Ошибка конвертации: %v", err)
			http.Error(w, "Не удается преобразовать содержимое", http.StatusInternalServerError)
			return
		}

		// Используем time.Now().UTC().String() как требуется в задании
		timestamp := time.Now().UTC().String()
		timestamp = strings.ReplaceAll(timestamp, ":", "-")
		timestamp = strings.ReplaceAll(timestamp, " ", "_")

		ext := filepath.Ext(fileHeader.Filename)
		if ext == "" {
			ext = ".txt"
		}

		safeFilename := strings.ReplaceAll(fileHeader.Filename, " ", "_")
		safeFilename = strings.ReplaceAll(safeFilename, "/", "_")
		safeFilename = strings.ReplaceAll(safeFilename, "\\", "_")

		outputFilename := fmt.Sprintf("translated_%s_%s%s", timestamp, safeFilename, ext)

		err = os.WriteFile(outputFilename, []byte(convertedContent), 0644)
		if err != nil {
			logger.Printf("Ошибка сохранения файла: %v", err)
			http.Error(w, "Не удалось сохранить файл", http.StatusInternalServerError)
			return
		}

		logger.Printf("Файл сохранён как: %s", outputFilename)

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		response := fmt.Sprintf("Файл успешно переведён!\nСохранён как: %s\n\nПеревод:\n%s", outputFilename, convertedContent)
		w.Write([]byte(response))
	}
}
