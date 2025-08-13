package repository

import (
	"context"
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, e := db.ExecContext(ctx, "INSERT INTO urls (original_url, shorten_id, clicks, created_at) VALUES (?, ?, ?, ?)",
		url.GetOriginalURL(), url.GetShortID(), url.GetClicks(), url.GetCreatedAt())

	if e != nil {
		fmt.Println("Error saving url to database:", e)
		return e
	}

	fmt.Println("Save url to database successfully!")

	return nil
}

func GetURLByShortID(shortenID string) (*URL, error) {
	db, err := config.GetDB()
	if err != nil {
		fmt.Println("Error initializing database:", err)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := db.QueryRowContext(ctx, "SELECT original_url, shorten_id, clicks, created_at FROM urls WHERE shorten_id = ?", shortenID)

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

	fmt.Println("Get url by shorten ID from database successfully!")

	return NewUrl(shortenID, parsedURL, clicks, createdAt), nil
}

func UpdateURL(url *URL) error {
	db, err := config.GetDB()
	if err != nil {
		fmt.Println("Error initializing database:", err)
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, e := db.ExecContext(ctx, "UPDATE urls SET clicks = ? WHERE shorten_id = ?",
		url.GetClicks(), url.GetShortID())

	if e != nil {
		fmt.Println("Error updating url in database:", e)
		return e
	}

	fmt.Println("Increment clicks for url in database successfully!")

	return nil
}
