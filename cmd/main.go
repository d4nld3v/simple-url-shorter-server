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

	fmt.Println("Original URL:", url)

	// check if the URL is valid
	valid, err := services.IsValidURL(url)
	if err != nil {
		fmt.Println("Error validating URL:", err)
		return
	}

	if valid {
		fmt.Println("URL is valid")
	}
	// generate a unique ID for the URL (shorten the URL)
	// store the original URL and the unique ID in a database
	// redirect to the original URL when accessing the unique ID

}
