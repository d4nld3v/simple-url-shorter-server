package main

import (
	"fmt"
	"os"

	db "github.com/d4nld3v/url-shortener-go/database"
	"github.com/d4nld3v/url-shortener-go/models"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Use: go run ./cmd <url>")
		return
	}

	url := os.Args[1]

	fmt.Println("Original URL:", url)

	shortenedURL, err := models.NewShortenedURL(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Shortened id:", shortenedURL.GetShortID())

	err = db.InitDatabase()
	if err != nil {
		fmt.Println("Error initializing database:", err)
		return
	}

	fmt.Println("Database initialized successfully")

	defer db.CloseDatabase()

}

// store the original URL and the unique ID in a database
// redirect to the original URL when accessing the unique ID

// store the original URL and the unique ID in a database
// redirect to the original URL when accessing the unique ID
