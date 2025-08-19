package handler

import (
	"encoding/json"
	"net/http"
	"time"
)

// ErrorResponse estructura estándar para errores
type ErrorResponse struct {
	Error     ErrorDetail `json:"error"`
	Timestamp string      `json:"timestamp"`
	Path      string      `json:"path"`
	Method    string      `json:"method"`
}

// ErrorDetail contiene los detalles del error
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// SuccessResponse estructura estándar para respuestas exitosas
type SuccessResponse struct {
	Data      interface{} `json:"data"`
	Message   string      `json:"message"`
	Timestamp string      `json:"timestamp"`
}

// ShortenURLResponse respuesta específica para URL acortada
type ShortenURLResponse struct {
	OriginalURL  string `json:"original_url"`
	ShortenedID  string `json:"shortened_id"`
	ShortenedURL string `json:"shortened_url"`
	CreatedAt    string `json:"created_at"`
}

// ErrorCodes constantes para códigos de error
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

// SendErrorResponse envía una respuesta de error estandarizada
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

// SendSuccessResponse envía una respuesta exitosa estandarizada
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

// SendCreatedResponse envía una respuesta 201 Created
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
