package main

import (
	"log"

	"github.com/d4nld3v/url-shortener-go/internal/config"
	"github.com/d4nld3v/url-shortener-go/internal/server"
)

func main() {

	cfg := config.Load()

	srv := server.New(cfg)
	log.Printf("Starting URL Shortener Service on %s", cfg.Addr)

	if err := srv.Start(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	defer config.CloseDatabase()
	log.Println("Server stopped gracefully")
}
