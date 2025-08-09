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
		valid, err := IsValidURL(test.url)

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
