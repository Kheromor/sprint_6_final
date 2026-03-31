package handlers

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/service"
)

// IndexHandler возвращает HTML-форму
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Парсим шаблон из корня проекта
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading template: %v", err), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
		return
	}
}

// UploadHandler обрабатывает загрузку файла
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Парсим форму (максимум 10 MB)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing form: %v", err), http.StatusInternalServerError)
		return
	}

	// Получаем файл из формы - используем имя поля "myFile"
	file, header, err := r.FormFile("myFile")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting file: %v", err), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Читаем содержимое файла
	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading file: %v", err), http.StatusInternalServerError)
		return
	}

	// Проверяем, что файл не пустой
	if len(data) == 0 {
		http.Error(w, "File is empty", http.StatusBadRequest)
		return
	}

	// Конвертируем через service
	converted, err := service.DetectAndConvert(string(data))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error converting: %v", err), http.StatusInternalServerError)
		return
	}

	// Генерируем имя файла согласно ТЗ: используем time.Now().UTC().String()
	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".txt"
	}

	// Используем time.Now().UTC().String() как указано в ТЗ
	// Заменяем двоеточия и пробелы на подчеркивания для безопасности файловой системы
	timestamp := time.Now().UTC().String()
	// Делаем имя файла безопасным для всех ОС
	timestamp = strings.ReplaceAll(timestamp, ":", "_")
	timestamp = strings.ReplaceAll(timestamp, " ", "_")
	filename := fmt.Sprintf("%s%s", timestamp, ext)

	// Создаем локальный файл
	outFile, err := os.Create(filename)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating file: %v", err), http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	// Записываем результат конвертации
	_, err = outFile.WriteString(converted)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error writing file: %v", err), http.StatusInternalServerError)
		return
	}

	// Возвращаем результат конвертации
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(converted))
}
