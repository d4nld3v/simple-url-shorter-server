package main

import (
	"fmt"
	"os"

	"github.com/d4nld3v/url-shortener-go/services"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Use: go run ./cmd <url>")
		return
	}

	url := os.Args[1]

	// validar la url
	valid, err := services.IsValidURL(url)
	if err != nil {
		fmt.Println("Error validating URL:", err)
		return
	}

	if valid {
		fmt.Println("URL is valid")
	}
	//generar una id unica para la url ( acortar la url )
	//almacenar la url original y la id unica en una base de datos
	//redirigir a la url original cuando se accede a la id unica

	fmt.Println("Original URL:", url)
}
