package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/datdev2409/lab-admin-go/internal/logger"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
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
	response := map[string]interface{}{
		"status": "success",
		"data":   data,
	}
	json.NewEncoder(w).Encode(response)
}

// RespondJSONWithPagination writes a JSON response with pagination for list endpoints
func RespondJSONWithPagination(w http.ResponseWriter, status int, data interface{}, pagination interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := map[string]interface{}{
		"status":     "success",
		"data":       data,
		"pagination": pagination,
	}
	json.NewEncoder(w).Encode(response)
}

// RespondError writes a JSON error response with the given status code and error
func RespondError(ctx context.Context, w http.ResponseWriter, status int, err error) {
	log := logger.FromCtx(ctx)
	log.Error("Handling error", zap.Error(err))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := map[string]interface{}{
		"status":    "error",
		"requestId": middleware.GetReqID(ctx),
		"error":     err.Error(),
	}
	json.NewEncoder(w).Encode(response)
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
			if appErr, ok := err.(*AppError); ok {
				RespondError(r.Context(), w, appErr.StatusCode, appErr)
			} else {
				RespondError(r.Context(), w, http.StatusInternalServerError, err)
			}
		}
	}
}
