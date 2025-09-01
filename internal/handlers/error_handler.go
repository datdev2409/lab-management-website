package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// AppError is a custom error type for business logic errors
// It includes an HTTP status code and a message
// Usage: return &AppError{StatusCode: 404, Message: "Not found"}
type AppError struct {
	StatusCode int
	Message    string
}

func (e *AppError) Error() string {
	return e.Message
}

type HandlerFuncReturnError = func(w http.ResponseWriter, r *http.Request) error

// RespondJSON writes a JSON response with the given status code and data
func RespondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// RespondError writes a JSON error response with the given status code and error
func RespondError(w http.ResponseWriter, status int, err error) {
	RespondJSON(w, status, map[string]string{"error": err.Error()})
}

// Utility functions to quickly raise common AppError types
func NotFoundError(message string) error {
	if message == "" {
		message = "Resource not found"
	}
	return &AppError{StatusCode: http.StatusNotFound, Message: message}
}

func UnauthorizedError(message string) error {
	if message == "" {
		message = "Unauthorized"
	}
	return &AppError{StatusCode: http.StatusUnauthorized, Message: message}
}

func BadRequestError(message string) error {
	if message == "" {
		message = "Bad request"
	}
	return &AppError{StatusCode: http.StatusBadRequest, Message: message}
}

func ForbiddenError(message string) error {
	if message == "" {
		message = "Forbidden"
	}
	return &AppError{StatusCode: http.StatusForbidden, Message: message}
}

func InternalServerError(message string) error {
	if message == "" {
		message = "Internal server error"
	}
	return &AppError{StatusCode: http.StatusInternalServerError, Message: message}
}

// Make wraps a handler and handles AppError responses
func Make(fn HandlerFuncReturnError) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			slog.Error("error", "error", err)
			if appErr, ok := err.(*AppError); ok {
				RespondError(w, appErr.StatusCode, appErr)
			} else {
				RespondError(w, http.StatusInternalServerError, err)
			}
		}
	}
}
