package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/handlers"
	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/server"
	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/service"
	//"github.com/Yandex-Practicum/go1fl-sprint6-final/pkg/morse"
)

func main() {
	// Создаем логгер
	logger := log.New(os.Stdout, "[SERVER] ", log.LstdFlags|log.Lshortfile)

	// Создаем функцию конвертации
	converter := func(data string) (string, error) {
		return service.DetectAndConvert(data)
	}

	// Создаем хендлеры
	indexHandler := handlers.IndexHandler(logger)
	uploadHandler := handlers.UploadHandler(logger, converter)

	// Создаем сервер
	srv := server.NewServer(logger, indexHandler, uploadHandler)

	// Запускаем сервер
	logger.Printf("Starting server on %s", srv.Server.Addr)
	if err := srv.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Server failed to start: %v", err)
	}
}
