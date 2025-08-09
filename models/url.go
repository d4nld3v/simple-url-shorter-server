package models

import (
	"crypto/md5"
	"encoding/base64"
	"net/url"
	"strings"

	"github.com/d4nld3v/url-shortener-go/services"
)

type ShortenedURL struct {
	originalURL *url.URL
	shortenID   string
	//TODO: generate shortened URL
	//shortenedUrl string
}

func NewShortenedURL(rawurl string) (*ShortenedURL, error) {

	u, err := services.IsValidURL(rawurl)
	if err != nil || u == nil {
		return nil, err
	}

	shortenID := ShortenURL(u)

	return &ShortenedURL{
		originalURL: u,
		shortenID:   shortenID,
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
	if s.originalURL == nil {
		return ""
	}
	return s.originalURL.String()
}

func (s *ShortenedURL) GetShortID() string {
	return s.shortenID
}
