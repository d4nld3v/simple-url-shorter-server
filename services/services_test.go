package services

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type test struct {
	url      string
	expected bool
}

type testWithError struct {
	url           string
	expectedValid bool
	expectedError string
}

func TestIsHttpURL(t *testing.T) {
	tests := []test{
		{"http://example.com", true},
		{"https://example.com", true},
		{"ftp://example.com", false},
		{"example.com", false},
	}

	for _, test := range tests {

		u, err := url.Parse(test.url)
		if err != nil {
			t.Errorf("Error parsing URL %s: %v", test.url, err)
			continue
		}

		result := isHttpURL(u)
		if result != test.expected {
			t.Errorf("isHttpURL(%s) = %v; want %v", u.String(), result, test.expected)
		}
	}
}

func TestIsPublicIP(t *testing.T) {
	tests := []test{
		{"https://google.com", true},   // Public domain
		{"https://github.com", true},   // Public domain
		{"https://localhost", false},   // Localhost
		{"https://127.0.0.1", false},   // Loopback IP
		{"https://192.168.1.1", false}, // Private IP
		{"https://10.0.0.1", false},    // Private IP
		{"https://172.16.0.1", false},  // Private IP
	}

	for _, test := range tests {
		u, err := url.Parse(test.url)
		if err != nil {
			t.Errorf("Error parsing URL %s: %v", test.url, err)
			continue
		}

		result := isPublicIP(u)
		if result != test.expected {
			t.Errorf("isPublicIP(%s) = %v; want %v", u.String(), result, test.expected)
		}
	}
}

func TestIsAvailable(t *testing.T) {
	//test server that responds with different status codes
	// to simulate different scenarios
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/success":
			w.WriteHeader(http.StatusOK)
		case "/redirect":
			w.WriteHeader(http.StatusMovedPermanently)
		case "/error":
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	tests := []test{
		{server.URL + "/success", true},                              // 200 OK
		{server.URL + "/redirect", true},                             // 301 Redirect
		{server.URL + "/error", false},                               // 404 Not Found
		{"https://invalid-url-that-does-not-exist-12345.com", false}, // Domain that does not exist
	}

	for _, test := range tests {
		u, err := url.Parse(test.url)
		if err != nil {
			t.Errorf("Error parsing URL %s: %v", test.url, err)
			continue
		}

		result := isAvailable(u)
		if result != test.expected {
			t.Errorf("isAvailable(%s) = %v; want %v", u.String(), result, test.expected)
		}
	}
}

func TestIsValidURL(t *testing.T) {
	tests := []testWithError{
		{"https://google.com", true, ""},
		{"http://github.com", true, ""},
		{"ftp://example.com", false, "url is not http or https"},
		{"", false, "url is null"},
		{"https://localhost", false, "url is not public IP"},
		{"https://127.0.0.1", false, "url is not public IP"},
		{"https://192.168.1.1", false, "url is not public IP"},
		{"invalid-url", false, "url is null"},
		{fmt.Sprintf("https://example.com/%s", string(make([]byte, 2100))), false, "url is too long"},
	}

	for _, test := range tests {
		u, err := IsValidURL(test.url)
		valid := u != nil

		if valid != test.expectedValid {
			t.Errorf("IsValidURL(%s) validity = %v; want %v", test.url, valid, test.expectedValid)
		}

		if test.expectedError != "" {
			if err == nil {
				t.Errorf("IsValidURL(%s) expected error containing '%s', but got no error", test.url, test.expectedError)
			} else if err.Error() == "" {
				t.Errorf("IsValidURL(%s) expected error containing '%s', but got empty error", test.url, test.expectedError)
			}
		} else if err != nil {
			t.Errorf("IsValidURL(%s) expected no error, but got: %v", test.url, err)
		}
	}
}

func TestNewShortenedURL(t *testing.T) {
	tests := []testWithError{
		{"https://google.com", true, ""},
		{"https://github.com/golang/go", true, ""},
		{"http://example.com", true, ""},
		{"ftp://example.com", false, "url is not http or https"},
		{"", false, "url is null"},
		{"https://localhost", false, "url is not public IP"},
		{"https://127.0.0.1", false, "url is not public IP"},
		{"https://192.168.1.1", false, "url is not public IP"},
		{"invalid-url", false, "url is null"},
		{fmt.Sprintf("https://example.com/%s", string(make([]byte, 2100))), false, "url is too long"},
	}

	for _, test := range tests {
		shortenedURL, err := NewShortenedURL(test.url)

		if test.expectedValid {
			// Si esperamos que sea válida
			if err != nil {
				t.Errorf("NewShortenedURL(%s) expected no error, but got: %v", test.url, err)
				continue
			}

			if shortenedURL == nil {
				t.Errorf("NewShortenedURL(%s) expected valid ShortenedURL, but got nil", test.url)
				continue
			}

			// Verificar que el ShortenID tenga el formato esperado (8 caracteres)
			if len(shortenedURL.GetShortID()) != 8 {
				t.Errorf("NewShortenedURL(%s) ShortenID length = %d; want 8", test.url, len(shortenedURL.GetShortID()))
			}

			// Verificar que la URL original se preserve correctamente
			if shortenedURL.GetOriginalURL() != test.url {
				t.Errorf("NewShortenedURL(%s) GetOriginalURL() = %s; want %s", test.url, shortenedURL.GetOriginalURL(), test.url)
			}

			// Verificar que el ShortenID no esté vacío
			if shortenedURL.GetShortID() == "" {
				t.Errorf("NewShortenedURL(%s) ShortenID is empty", test.url)
			}

		} else {
			// Si esperamos que sea inválida
			if err == nil {
				t.Errorf("NewShortenedURL(%s) expected error, but got none", test.url)
			}

			if shortenedURL != nil {
				t.Errorf("NewShortenedURL(%s) expected nil ShortenedURL, but got valid instance", test.url)
			}
		}
	}
}

func TestShortenedURLConsistency(t *testing.T) {
	// Test que URLs iguales generen el mismo ShortenID
	url1, err1 := NewShortenedURL("https://google.com")
	url2, err2 := NewShortenedURL("https://google.com")

	if err1 != nil || err2 != nil {
		t.Fatalf("Expected no errors, but got: err1=%v, err2=%v", err1, err2)
	}

	if url1 == nil || url2 == nil {
		t.Fatal("Expected valid ShortenedURL instances")
	}

	if url1.GetShortID() != url2.GetShortID() {
		t.Errorf("Same URLs should generate same ShortenID: %s != %s", url1.GetShortID(), url2.GetShortID())
	}

	// Test que URLs diferentes generen ShortenIDs diferentes
	url3, err3 := NewShortenedURL("https://github.com")
	if err3 != nil {
		t.Fatalf("Expected no error, but got: %v", err3)
	}

	if url3 == nil {
		t.Fatal("Expected valid ShortenedURL instance")
	}

	if url1.GetShortID() == url3.GetShortID() {
		t.Errorf("Different URLs should generate different ShortenIDs: %s == %s", url1.GetShortID(), url3.GetShortID())
	}
}

func TestShortenedURLMethods(t *testing.T) {
	shortenedURL, err := NewShortenedURL("https://www.google.com/search?q=golang")

	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	if shortenedURL == nil {
		t.Fatal("Expected valid ShortenedURL instance")
	}

	// Test GetOriginalURL
	expectedURL := "https://www.google.com/search?q=golang"
	if shortenedURL.GetOriginalURL() != expectedURL {
		t.Errorf("GetOriginalURL() = %s; want %s", shortenedURL.GetOriginalURL(), expectedURL)
	}

	// Test GetShortID
	shortID := shortenedURL.GetShortID()
	if shortID == "" {
		t.Error("GetShortID() returned empty string")
	}

	if len(shortID) != 8 {
		t.Errorf("GetShortID() length = %d; want 8", len(shortID))
	}
}
