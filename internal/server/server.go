package server

import (
	"log"
	"net/http"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/handlers"
)

// Server представляет HTTP сервер
type Server struct {
	Logger *log.Logger
	Server *http.Server
}

// NewServer создает новый сервер
func NewServer(logger *log.Logger) *Server {
	// Создаем роутер
	mux := http.NewServeMux()

	// Регистрируем хендлеры
	mux.HandleFunc("/", handlers.IndexHandler)
	mux.HandleFunc("/upload", handlers.UploadHandler)

	// Настраиваем HTTP сервер
	httpServer := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &Server{
		Logger: logger,
		Server: httpServer,
	}
}
