package server

import (
	"log"
	"net/http"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/handlers"
)

// Server представляет структуру сервера
type Server struct {
	Logger *log.Logger
	Server *http.Server
}

// NewServer создает новый экземпляр сервера
func NewServer(logger *log.Logger) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.IndexHandler(logger))
	mux.HandleFunc("/upload", handlers.UploadHandler(logger))

	return &Server{
		Logger: logger,
		Server: &http.Server{
			Addr:         ":8080",
			Handler:      mux,
			ErrorLog:     logger,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Second,
		},
	}
}
