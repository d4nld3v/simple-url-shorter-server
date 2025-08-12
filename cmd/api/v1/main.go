package main

import (
	"log"

	"github.com/d4nld3v/url-shortener-go/internal/config"
	"github.com/d4nld3v/url-shortener-go/internal/server"
)

func main() {

	//TODO:

	//* Implement get clicks from shortened URL

	// buscar una forma de inicializar la base de datos al iniciar el servidor
	// y no al guardar un nuevo URL. (BUSCAR ALGÚN PATRON DE DISEÑO)

	// manejar mejor los errores tanto en la base de datos como en la API ( buscar patron de diseño)

	// revisar si las consultas a la base de datos están bien implementadas ( evitando injection / preparedStatement)

	// crear un README.md con la documentación de la API

	cfg := config.Load()

	srv := server.New(cfg)

	log.Printf("Starting URL Shortener Service on %s", cfg.Addr)

	if err := srv.Start(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
