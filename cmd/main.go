package main

import (
	"log"
	"os"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/server"
)

func main() {
	// Создаем логгер
	logger := log.New(os.Stdout, "SERVER: ", log.Ldate|log.Ltime|log.Lshortfile)

	// Создаем сервер
	srv := server.NewServer(logger)

	// Запускаем сервер
	logger.Println("Starting server on http://localhost:8080")
	if err := srv.Server.ListenAndServe(); err != nil {
		logger.Fatal("Server failed:", err)
	}
}
