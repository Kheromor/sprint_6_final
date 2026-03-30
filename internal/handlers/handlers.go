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
)

func IndexHandler(logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		// Получаем путь к исполняемому файлу
		execPath, err := os.Executable()
		if err != nil {
			logger.Printf("Ошибка при получении пути к исполняемому файлу: %v", err)
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
			return
		}

		// Строим путь к index.html относительно исполняемого файла
		htmlPath := filepath.Join(filepath.Dir(execPath), "index.html")

		// Если файл не найден, пробуем искать в текущей директории
		if _, err := os.Stat(htmlPath); os.IsNotExist(err) {
			htmlPath = "index.html"
		}

		// Читаем файл index.html
		htmlContent, err := os.ReadFile(htmlPath)
		if err != nil {
			logger.Printf("Ошибка чтения index.html из пути %s: %v", htmlPath, err)

			// Возвращаем простую HTML-форму, если файл не найден
			fallbackHTML := `<!DOCTYPE html>
<html>
<head>
    <title>Преобразователь азбуки Морзе</title>
</head>
<body>
    <h1>Преобразователь азбуки Морзе</h1>
    <form action="/upload" method="post" enctype="multipart/form-data">
        <input type="file" name="file" required>
        <button type="submit">Перевести</button>
    </form>
</body>
</html>`
			w.Header().Set("Тип содержимого", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fallbackHTML))
			return
		}

		// Устанавливаем заголовок
		w.Header().Set("Тип содержимого", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(htmlContent)
	}
}

// UploadHandler обрабатывает загрузку файлов
func UploadHandler(logger *log.Logger, converter func(string) (string, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
			return
		}

		// Парсим форму
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			logger.Printf("Ошибка синтаксиса формы: %v", err)
			http.Error(w, "Не удаётся разобрать форму", http.StatusBadRequest)
			return
		}

		// Получаем файл из формы
		file, handler, err := r.FormFile("file")
		if err != nil {
			logger.Printf("Ошибка при извлечении файла: %v", err)
			http.Error(w, "Не удается получить файл", http.StatusBadRequest)
			return
		}
		defer file.Close()

		logger.Printf("Полученный файл: %s", handler.Filename)

		// Читаем данные из файла
		fileContent, err := io.ReadAll(file)
		if err != nil {
			logger.Printf("Ошибка при чтении файла: %v", err)
			http.Error(w, "Не удается прочитать файл", http.StatusInternalServerError)
			return
		}

		logger.Printf("Длина содержимого файла: %d байтов", len(fileContent))
		logger.Printf("Предварительный просмотр содержимого файла: %s", string(fileContent)[:min(100, len(fileContent))])

		// Конвертируем данные
		convertedContent, err := converter(string(fileContent))
		if err != nil {
			logger.Printf("Ошибка при преобразовании содержимого: %v", err)
			http.Error(w, "Не удается преобразовать содержимое", http.StatusInternalServerError)
			return
		}

		// Создаем локальный файл с результатом
		// Генерируем уникальное имя файла
		timestamp := time.Now().UTC().Format("20060102_150405.000")
		ext := filepath.Ext(handler.Filename)
		if ext == "" {
			ext = ".txt"
		}

		// Создаем имя файла без специальных символов
		safeFilename := strings.ReplaceAll(handler.Filename, " ", "_")
		safeFilename = strings.ReplaceAll(safeFilename, "/", "_")
		safeFilename = strings.ReplaceAll(safeFilename, "\\", "_")

		outputFilename := fmt.Sprintf("translated_%s_%s%s",
			timestamp,
			safeFilename,
			ext)

		// Создаем выходной файл
		outputFile, err := os.Create(outputFilename)
		if err != nil {
			logger.Printf("Ошибка при создании выходного файла: %v", err)
			http.Error(w, "Не удалось создать выходной файл", http.StatusInternalServerError)
			return
		}
		defer outputFile.Close()

		// Записываем конвертированное содержимое
		_, err = outputFile.WriteString(convertedContent)
		if err != nil {
			logger.Printf("Ошибка записи в выходной файл: %v", err)
			http.Error(w, "Не удается выполнить запись в выходной файл", http.StatusInternalServerError)
			return
		}

		logger.Printf("Успешно преобразован и сохранен в: %s", outputFilename)

		// Возвращаем результат конвертации
		w.Header().Set("Тип содержимого", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		// Формируем ответ с информацией о сохраненном файле и результатом
		response := fmt.Sprintf("Файл успешно переведён!\nСохранён как: %s\n\nПеревод:\n%s",
			outputFilename,
			convertedContent)

		w.Write([]byte(response))
	}
}

// Вспомогательная функция
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
