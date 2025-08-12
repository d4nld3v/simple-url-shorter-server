package services

import (
	"crypto/md5"
	"encoding/base64"
	"net/url"
	"strings"
	"time"

	"github.com/d4nld3v/url-shortener-go/internal/repository"
)

func shortenURL(url *url.URL) string {

	hash := md5.Sum([]byte(url.String()))

	encoded := base64.URLEncoding.EncodeToString(hash[:])

	shortID := encoded[:8]
	shortID = strings.ReplaceAll(shortID, "-", "x")
	shortID = strings.ReplaceAll(shortID, "_", "y")

	return shortID
}

func ConvertToShorterUrl(u *url.URL) (*repository.URL, error) {

	shortID := shortenURL(u)

	convertedURL := repository.NewUrl(shortID, u, 0, time.Now())

	err := repository.SaveShortenedURL(convertedURL)
	if err != nil {
		// Handle error (e.g., log it)
		return nil, err
	}

	return convertedURL, nil
}

func GetShortenedURL(shortID string) (*repository.URL, error) {

	// Assuming a function exists in the repository to fetch the URL by its short ID
	url, err := repository.GetURLByShortID(shortID)
	if err != nil {
		return nil, err
	}

	if url == nil {
		return nil, nil // Not found
	}

	// Increment the click count
	url.IncrementClicks()
	err = repository.UpdateURL(url)
	if err != nil {
		return nil, err
	}

	return url, nil
}
