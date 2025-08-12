package repository

import (
	"fmt"
	"net/url"
	"time"

	"github.com/d4nld3v/url-shortener-go/internal/config"
)

type URL struct {
	shortenID   string
	originalURL *url.URL
	clicks      int
	createdAt   time.Time
}

func NewUrl(shortenID string, originalURL *url.URL, clicks int, createdAt time.Time) *URL {
	return &URL{
		shortenID:   shortenID,
		originalURL: originalURL,
		clicks:      clicks,
		createdAt:   createdAt,
	}
}

func (s *URL) GetOriginalURL() string {
	if s.originalURL == nil {
		return ""
	}
	return s.originalURL.String()
}

func (s *URL) GetShortID() string {
	return s.shortenID
}

func (s *URL) GetClicks() int {
	return s.clicks
}

func (s *URL) IncrementClicks() {
	s.clicks++
}

func (s *URL) GetCreatedAt() time.Time {
	return s.createdAt
}

func SaveShortenedURL(url *URL) error {

	db, err := config.GetDB()
	if err != nil {
		fmt.Println("Error initializing database:", err)
		return err
	}

	fmt.Println("Save url to database successfully!")

	_, e := db.Exec("INSERT INTO urls (original_url, shorten_id, clicks, created_at) VALUES (?, ?, ?, ?)",
		url.GetOriginalURL(), url.GetShortID(), url.GetClicks(), url.GetCreatedAt())

	defer config.CloseDatabase()
	return e
}

func GetURLByShortID(shortenID string) (*URL, error) {
	db, err := config.GetDB()
	if err != nil {
		fmt.Println("Error initializing database:", err)
		return nil, err
	}

	fmt.Println("Get url by shorten ID from database successfully!")

	row := db.QueryRow("SELECT original_url, shorten_id, clicks, created_at FROM urls WHERE shorten_id = ?", shortenID)

	var originalURL string
	var clicks int
	var createdAt time.Time

	err = row.Scan(&originalURL, &shortenID, &clicks, &createdAt)
	if err != nil {
		return nil, err
	}

	parsedURL, parseErr := url.Parse(originalURL)
	if parseErr != nil {
		return nil, parseErr
	}

	return NewUrl(shortenID, parsedURL, clicks, createdAt), nil
}

func UpdateURL(url *URL) error {
	db, err := config.GetDB()
	if err != nil {
		fmt.Println("Error initializing database:", err)
		return err
	}

	fmt.Println("Increment clicks for url in database successfully!")

	_, e := db.Exec("UPDATE urls SET clicks = ? WHERE shorten_id = ?",
		url.GetClicks(), url.GetShortID())

	defer config.CloseDatabase()
	return e
}
