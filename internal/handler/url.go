package handler

import (
	"encoding/json"
	"net/http"

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
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	handlers := map[string]http.HandlerFunc{

		"GET": redirect,
	}

	if handler, ok := handlers[r.Method]; ok {
		handler(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handlePostShortenURL(w http.ResponseWriter, r *http.Request) {

	var requestBody struct {
		Url string `json:"url"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if requestBody.Url == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	u, err := services.IsValidURL(requestBody.Url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortenedUrl, err := services.ConvertToShorterUrl(u)

	if err != nil {
		http.Error(w, "Failed to shorten URL: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Construir la URL completa del acortador con /short/ prefix
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	host := r.Host
	if host == "" {
		host = "localhost:8080" // fallback por defecto
	}

	shortenedURL := scheme + "://" + host + "/short/" + shortenedUrl.GetShortID()

	json.NewEncoder(w).Encode(map[string]string{
		"original_url":  requestBody.Url,
		"shortened_id":  shortenedUrl.GetShortID(),
		"shortened_url": shortenedURL,
		"message":       "URL shortened successfully",
	})
}

func redirect(w http.ResponseWriter, r *http.Request) {

	// Extraer solo el shortID removiendo el prefijo "/short/"
	shortID := r.URL.Path[7:] // "/short/" tiene 7 caracteres

	println("Redirecting for short ID:", shortID)

	if len(shortID) > 8 {
		http.Error(w, "Short ID too long", http.StatusBadRequest)
		return
	}

	if shortID == "" {
		http.Error(w, "Short ID is required", http.StatusBadRequest)
		return
	}

	u, err := services.GetShortenedURL(shortID)

	if err != nil || u == nil {
		http.Error(w, "Failed to retrieve shortened URL: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, u.GetOriginalURL(), http.StatusFound)
}
