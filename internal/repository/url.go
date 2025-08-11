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

	err := config.InitDatabase()
	if err != nil {
		fmt.Println("Error initializing database:", err)
		return err
	}

	fmt.Println("Database initialized successfully")

	_, e := config.DB.Exec("INSERT INTO urls (original_url, shorten_id, clicks, created_at) VALUES (?, ?, ?, ?)",
		url.GetOriginalURL(), url.GetShortID(), url.GetClicks(), url.GetCreatedAt())

	defer config.CloseDatabase()
	return e
}
