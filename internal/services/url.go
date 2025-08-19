package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/d4nld3v/url-shortener-go/internal/repository"
)

func shortenURL(url *url.URL) string {
	// Usar SHA-256 + timestamp + random para evitar colisiones
	timestamp := fmt.Sprintf("%d", time.Now().UnixNano())
	randomBytes := make([]byte, 8)
	rand.Read(randomBytes)

	input := url.String() + timestamp + string(randomBytes)
	hash := sha256.Sum256([]byte(input))

	encoded := base64.URLEncoding.EncodeToString(hash[:])

	shortID := encoded[:8]
	shortID = strings.ReplaceAll(shortID, "-", "x")
	shortID = strings.ReplaceAll(shortID, "_", "y")
	shortID = strings.ReplaceAll(shortID, "+", "z")

	return shortID
}

func ConvertToShorterUrl(u *url.URL) (*repository.URL, error) {

	normalizedURL := normalizeURL(u)

	fmt.Println("Normalized URL:", normalizedURL.String())

	var shortID string
	var convertedURL *repository.URL
	maxRetries := 5

	for i := 0; i < maxRetries; i++ {
		shortID = shortenURL(normalizedURL)

		fmt.Println("Attempt", i+1, "to generate short ID:", shortID)

		existing, _ := repository.GetURLByShortID(shortID)

		if existing == nil {
			// shortID disponible
			break
		}

		if i == maxRetries-1 {
			return nil, fmt.Errorf("failed to generate unique short ID after %d attempts", maxRetries)
		}
	}

	fmt.Println("Generated short ID:", shortID)

	convertedURL = repository.NewUrl(shortID, normalizedURL, 0, time.Now())

	err := repository.SaveShortenedURL(convertedURL)
	if err != nil {
		return nil, fmt.Errorf("failed to save shortened URL: %w", err)
	}

	return convertedURL, nil
}

// normalizeURL estandariza las URLs para evitar duplicados
func normalizeURL(u *url.URL) *url.URL {
	normalized := *u

	normalized.Scheme = strings.ToLower(normalized.Scheme)
	normalized.Host = strings.ToLower(normalized.Host)

	if (normalized.Scheme == "http" && strings.HasSuffix(normalized.Host, ":80")) ||
		(normalized.Scheme == "https" && strings.HasSuffix(normalized.Host, ":443")) {
		normalized.Host = normalized.Host[:strings.LastIndex(normalized.Host, ":")]
	}

	normalized.Fragment = ""

	if normalized.Path == "" {
		normalized.Path = "/"
	}

	return &normalized
}

func GetShortenedURL(shortID string) (*repository.URL, error) {
	if err := ValidateShortID(shortID); err != nil {
		return nil, fmt.Errorf("invalid short ID: %w", err)
	}

	url, err := repository.GetURLByShortID(shortID)
	if err != nil {
		return nil, fmt.Errorf("failed to get URL: %w", err)
	}

	if url == nil {
		return nil, nil // Not found
	}

	url.IncrementClicks()
	if err := repository.UpdateURL(url); err != nil {
		return nil, fmt.Errorf("failed to increment clicks: %w", err)
	}

	return url, nil
}

func ValidateShortID(shortID string) error {
	if strings.TrimSpace(shortID) == "" {
		return fmt.Errorf("short ID cannot be empty")
	}

	if len(shortID) < 1 || len(shortID) > 10 {
		return fmt.Errorf("short ID must be between 1 and 10 characters")
	}

	// Solo caracteres alfanuméricos y caracteres seguros
	for _, char := range shortID {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == 'x' || char == 'y' || char == 'z') {
			return fmt.Errorf("short ID contains invalid characters")
		}
	}

	return nil
}

func GetURLStats(shortID string) (*repository.URL, error) {
	if err := ValidateShortID(shortID); err != nil {
		return nil, fmt.Errorf("invalid short ID: %w", err)
	}

	stats, err := repository.GetURLStatsByShortID(shortID)
	if err != nil {
		return nil, fmt.Errorf("failed to get URL stats: %w", err)
	}

	if stats == nil {
		return nil, nil // Not found
	}

	return stats, nil
}
