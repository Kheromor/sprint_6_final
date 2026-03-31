package main

import (
	"log"
	"os"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/server"
)

func main() {
	logger := log.New(os.Stdout, "[SERVER] ", log.LstdFlags)
	logger.Println("Starting server on :8080")

	if err := server.NewServer(logger).Server.ListenAndServe(); err != nil {
		logger.Fatal(err)
	}
}
