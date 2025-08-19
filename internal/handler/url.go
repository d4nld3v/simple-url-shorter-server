package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/d4nld3v/url-shortener-go/internal/services"
	"github.com/d4nld3v/url-shortener-go/pkg/middleware"
)

func RegisterUrlRoutes(mux *http.ServeMux, rl *middleware.RateLimiter) {
	mux.Handle("/shorten", rl.Limit(http.HandlerFunc(shortenURLHandler)))
	mux.Handle("/short/", rl.Limit(http.HandlerFunc(redirectHandler)))
}

func shortenURLHandler(w http.ResponseWriter, r *http.Request) {
	handlers := map[string]http.HandlerFunc{
		"POST": handlePostShortenURL,
	}

	if handler, ok := handlers[r.Method]; ok {
		handler(w, r)
	} else {
		SendErrorResponse(w, r, http.StatusMethodNotAllowed, ErrCodeMethodNotAllowed,
			"Method not allowed", "Only POST method is supported for this endpoint")
	}
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	handlers := map[string]http.HandlerFunc{
		"GET": redirect,
	}

	if handler, ok := handlers[r.Method]; ok {
		handler(w, r)
	} else {
		SendErrorResponse(w, r, http.StatusMethodNotAllowed, ErrCodeMethodNotAllowed,
			"Method not allowed", "Only GET method is supported for this endpoint")
	}
}

func handlePostShortenURL(w http.ResponseWriter, r *http.Request) {
	middleware.SetSecurityHeaders(w)

	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		SendErrorResponse(w, r, http.StatusUnsupportedMediaType, ErrCodeInvalidInput,
			"Invalid Content-Type", "Content-Type must be application/json")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024) // 1MB max

	var requestBody struct {
		Url string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&requestBody)
	if err != nil {
		SendErrorResponse(w, r, http.StatusBadRequest, ErrCodeInvalidInput,
			"Invalid JSON format", err.Error())
		return
	}

	if strings.TrimSpace(requestBody.Url) == "" {
		SendErrorResponse(w, r, http.StatusBadRequest, ErrCodeInvalidInput,
			"URL is required", "The 'url' field cannot be empty")
		return
	}

	if len(requestBody.Url) > 2048 {
		SendErrorResponse(w, r, http.StatusBadRequest, ErrCodeURLTooLong,
			"URL too long", "URL must be less than 2048 characters")
		return
	}

	if !utf8.ValidString(requestBody.Url) {
		SendErrorResponse(w, r, http.StatusBadRequest, ErrCodeInvalidInput,
			"Invalid URL encoding", "URL contains invalid UTF-8 characters")
		return
	}

	// Validar URL
	u, err := services.IsValidURL(requestBody.Url)
	if err != nil {
		SendErrorResponse(w, r, http.StatusBadRequest, ErrCodeInvalidURL,
			"Invalid URL", err.Error())
		return
	}

	// Crear URL acortada
	shortenedUrl, err := services.ConvertToShorterUrl(u)
	if err != nil {
		SendErrorResponse(w, r, http.StatusInternalServerError, ErrCodeInternalError,
			"Failed to shorten URL", err.Error())
		return
	}

	// Construir URL completa
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	host := r.Host
	if host == "" {
		host = "localhost:8080"
	}

	shortenedURL := scheme + "://" + host + "/short/" + shortenedUrl.GetShortID()

	responseData := ShortenURLResponse{
		OriginalURL:  requestBody.Url,
		ShortenedID:  shortenedUrl.GetShortID(),
		ShortenedURL: shortenedURL,
		CreatedAt:    shortenedUrl.GetCreatedAt().Format(time.RFC3339),
	}

	SendCreatedResponse(w, responseData, "URL shortened successfully")
}

func redirect(w http.ResponseWriter, r *http.Request) {
	// Extraer shortID removiendo el prefijo "/short/"
	shortID := r.URL.Path[7:] // "/short/" tiene 7 caracteres

	if err := services.ValidateShortID(shortID); err != nil {
		SendErrorResponse(w, r, http.StatusBadRequest, ErrCodeInvalidShortID,
			"Invalid short ID", err.Error())
		return
	}
	u, err := services.GetShortenedURL(shortID)
	if err != nil {
		SendErrorResponse(w, r, http.StatusInternalServerError, ErrCodeInternalError,
			"Failed to retrieve URL", err.Error())
		return
	}

	if u == nil {
		SendErrorResponse(w, r, http.StatusNotFound, ErrCodeURLNotFound,
			"URL not found", "The requested short URL does not exist")
		return
	}

	http.Redirect(w, r, u.GetOriginalURL(), http.StatusFound)
}
