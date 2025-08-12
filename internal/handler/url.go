package handler

import (
	"encoding/json"
	"net/http"

	"github.com/d4nld3v/url-shortener-go/internal/services"
)

func RegisterUrlRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/shorten", shortenURLHandler)
}

func shortenURLHandler(w http.ResponseWriter, r *http.Request) {

	handlers := map[string]http.HandlerFunc{
		"POST": handlePostShortenURL,
		"GET":  handleGetShortenURL,
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

	json.NewEncoder(w).Encode(map[string]string{
		"original_url": requestBody.Url,
		"shortened_id": shortenedUrl.GetShortID(),
		"message":      "URL received successfully",
	})
}

func handleGetShortenURL(w http.ResponseWriter, r *http.Request) {

	shortID := r.URL.Query().Get("id")

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
