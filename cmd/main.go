package main

import (
	"fmt"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Use: go run ./cmd <url>")
		return
	}

	url := os.Args[1]

	fmt.Println("Original URL:", url)
}

// funcion para acortar la url original
// generar una id unica mas corta para la url
// almacenar la url original y la id unica en una base de datos
// redirigir a la url original cuando se accede a la id unica
// funcion para consultar la url original a partir de la id unica
