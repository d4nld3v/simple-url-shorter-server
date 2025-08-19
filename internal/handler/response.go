package handler

import (
	"encoding/json"
	"net/http"
	"time"
)

type ErrorResponse struct {
	Error     ErrorDetail `json:"error"`
	Timestamp string      `json:"timestamp"`
	Path      string      `json:"path"`
	Method    string      `json:"method"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

type SuccessResponse struct {
	Data      interface{} `json:"data"`
	Message   string      `json:"message"`
	Timestamp string      `json:"timestamp"`
}

type ShortenURLResponse struct {
	OriginalURL  string `json:"original_url"`
	ShortenedID  string `json:"shortened_id"`
	ShortenedURL string `json:"shortened_url"`
	CreatedAt    string `json:"created_at"`
}

const (
	ErrCodeInvalidInput      = "INVALID_INPUT"
	ErrCodeURLNotFound       = "URL_NOT_FOUND"
	ErrCodeInternalError     = "INTERNAL_ERROR"
	ErrCodeValidationFailed  = "VALIDATION_FAILED"
	ErrCodeRateLimitExceeded = "RATE_LIMIT_EXCEEDED"
	ErrCodeMethodNotAllowed  = "METHOD_NOT_ALLOWED"
	ErrCodeInvalidURL        = "INVALID_URL"
	ErrCodeURLTooLong        = "URL_TOO_LONG"
	ErrCodeInvalidShortID    = "INVALID_SHORT_ID"
)

func SendErrorResponse(w http.ResponseWriter, r *http.Request, statusCode int, errorCode, message, details string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorResponse := ErrorResponse{
		Error: ErrorDetail{
			Code:    errorCode,
			Message: message,
			Details: details,
		},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Path:      r.URL.Path,
		Method:    r.Method,
	}

	json.NewEncoder(w).Encode(errorResponse)
}

func SendSuccessResponse(w http.ResponseWriter, data interface{}, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	successResponse := SuccessResponse{
		Data:      data,
		Message:   message,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	json.NewEncoder(w).Encode(successResponse)
}

func SendCreatedResponse(w http.ResponseWriter, data interface{}, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	successResponse := SuccessResponse{
		Data:      data,
		Message:   message,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	json.NewEncoder(w).Encode(successResponse)
}
