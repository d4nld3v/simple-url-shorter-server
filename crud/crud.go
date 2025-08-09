package crud

import (
	db "github.com/d4nld3v/url-shortener-go/database"
	"github.com/d4nld3v/url-shortener-go/models"
)

// save the original URL and its unique ID in the database

//get the original URL by its unique ID from the database

func SaveShortenedURL(url *models.ShortenedURL) error {
	_, err := db.DB.Exec("INSERT INTO urls (original_url, shorten_id) VALUES (?, ?)", url.GetOriginalURL(), url.GetShortID())
	return err
}
