package server

import (
	"log"
	"net/http"
	"time"
)

// Server структура сервера
type Server struct {
	Logger *log.Logger
	Server *http.Server
}

// NewServer создает новый HTTP сервер с настроенными хендлерами
func NewServer(logger *log.Logger, indexHandler http.HandlerFunc, uploadHandler http.HandlerFunc) *Server {
	// Создаем роутер
	mux := http.NewServeMux()
	
	// Регистрируем хендлеры
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/upload", uploadHandler)
	
	// Создаем HTTP сервер
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