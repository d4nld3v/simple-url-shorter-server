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

	shortenedURL := services.NewShortenedURL(url)
	if shortenedURL == nil {
		fmt.Println("Invalid URL")
		return
	}

	fmt.Println("Shortened id:", shortenedURL.GetShortID())
}

// store the original URL and the unique ID in a database
// redirect to the original URL when accessing the unique ID
