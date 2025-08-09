package services

import (
	"crypto/md5"
	"encoding/base64"
	"net/url"
	"strings"
)

type ShortenedURL struct {
	url       *url.URL
	shortenID string
}

func NewShortenedURL(originalURL string) (*ShortenedURL, error) {

	u, err := IsValidURL(originalURL)
	if err != nil || u == nil {
		return nil, err
	}

	shortenID := ShortenURL(u)

	return &ShortenedURL{
		url:       u,
		shortenID: shortenID,
	}, nil
}

func ShortenURL(url *url.URL) string {
	// Crear hash MD5 de la URL completa
	hash := md5.Sum([]byte(url.String()))

	// Codificar en base64 para obtener caracteres alfanuméricos
	encoded := base64.URLEncoding.EncodeToString(hash[:])

	shortID := encoded[:8]
	shortID = strings.ReplaceAll(shortID, "-", "x")
	shortID = strings.ReplaceAll(shortID, "_", "y")

	return shortID
}

func (s *ShortenedURL) GetOriginalURL() string {
	if s.url == nil {
		return ""
	}
	return s.url.String()
}

func (s *ShortenedURL) GetShortID() string {
	return s.shortenID
}
