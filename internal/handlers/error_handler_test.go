package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppError_Error(t *testing.T) {
	err := &AppError{
		StatusCode: 404,
		Message:    "Resource not found",
	}
	
	assert.Equal(t, "Resource not found", err.Error())
}

func TestRespondJSON(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]string{"name": "test"}
	
	RespondJSON(w, http.StatusOK, data)
	
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	
	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)
	
	assert.Equal(t, "success", response["status"])
	assert.NotNil(t, response["data"])
}

func TestRespondJSONWithPagination(t *testing.T) {
	w := httptest.NewRecorder()
	data := []string{"item1", "item2"}
	pagination := map[string]int{"page": 1, "perPage": 10}
	
	RespondJSONWithPagination(w, http.StatusOK, data, pagination)
	
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	
	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)
	
	assert.Equal(t, "success", response["status"])
	assert.NotNil(t, response["data"])
	assert.NotNil(t, response["pagination"])
}

func TestRespondError(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	testErr := errors.New("test error")
	
	RespondError(req.Context(), w, http.StatusBadRequest, testErr)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	
	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)
	
	assert.Equal(t, "error", response["status"])
	assert.Equal(t, "test error", response["error"])
}

func TestNotFoundError(t *testing.T) {
	t.Run("with custom message", func(t *testing.T) {
		err := NotFoundError("Custom not found")
		
		appErr, ok := err.(*AppError)
		require.True(t, ok)
		assert.Equal(t, http.StatusNotFound, appErr.StatusCode)
		assert.Equal(t, "Custom not found", appErr.Message)
	})
	
	t.Run("with empty message", func(t *testing.T) {
		err := NotFoundError("")
		
		appErr, ok := err.(*AppError)
		require.True(t, ok)
		assert.Equal(t, http.StatusNotFound, appErr.StatusCode)
		assert.Equal(t, "Resource not found", appErr.Message)
	})
}

func TestUnauthorizedError(t *testing.T) {
	t.Run("with custom message", func(t *testing.T) {
		err := UnauthorizedError("Custom unauthorized")
		
		appErr, ok := err.(*AppError)
		require.True(t, ok)
		assert.Equal(t, http.StatusUnauthorized, appErr.StatusCode)
		assert.Equal(t, "Custom unauthorized", appErr.Message)
	})
	
	t.Run("with empty message", func(t *testing.T) {
		err := UnauthorizedError("")
		
		appErr, ok := err.(*AppError)
		require.True(t, ok)
		assert.Equal(t, http.StatusUnauthorized, appErr.StatusCode)
		assert.Equal(t, "Unauthorized", appErr.Message)
	})
}

func TestBadRequestError(t *testing.T) {
	t.Run("with custom message", func(t *testing.T) {
		err := BadRequestError("Invalid input")
		
		appErr, ok := err.(*AppError)
		require.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, appErr.StatusCode)
		assert.Equal(t, "Invalid input", appErr.Message)
	})
	
	t.Run("with empty message", func(t *testing.T) {
		err := BadRequestError("")
		
		appErr, ok := err.(*AppError)
		require.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, appErr.StatusCode)
		assert.Equal(t, "Bad request", appErr.Message)
	})
}

func TestForbiddenError(t *testing.T) {
	t.Run("with custom message", func(t *testing.T) {
		err := ForbiddenError("Access denied")
		
		appErr, ok := err.(*AppError)
		require.True(t, ok)
		assert.Equal(t, http.StatusForbidden, appErr.StatusCode)
		assert.Equal(t, "Access denied", appErr.Message)
	})
	
	t.Run("with empty message", func(t *testing.T) {
		err := ForbiddenError("")
		
		appErr, ok := err.(*AppError)
		require.True(t, ok)
		assert.Equal(t, http.StatusForbidden, appErr.StatusCode)
		assert.Equal(t, "Forbidden", appErr.Message)
	})
}

func TestInternalServerError(t *testing.T) {
	t.Run("with custom message", func(t *testing.T) {
		err := InternalServerError("Database error")
		
		appErr, ok := err.(*AppError)
		require.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, appErr.StatusCode)
		assert.Equal(t, "Database error", appErr.Message)
	})
	
	t.Run("with empty message", func(t *testing.T) {
		err := InternalServerError("")
		
		appErr, ok := err.(*AppError)
		require.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, appErr.StatusCode)
		assert.Equal(t, "Internal server error", appErr.Message)
	})
}

func TestMake(t *testing.T) {
	t.Run("handler returns no error", func(t *testing.T) {
		handler := Make(func(w http.ResponseWriter, r *http.Request) error {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success"))
			return nil
		})
		
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		
		handler(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "success", w.Body.String())
	})
	
	t.Run("handler returns AppError", func(t *testing.T) {
		handler := Make(func(w http.ResponseWriter, r *http.Request) error {
			return &AppError{
				StatusCode: http.StatusNotFound,
				Message:    "Not found",
			}
		})
		
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		
		handler(w, req)
		
		assert.Equal(t, http.StatusNotFound, w.Code)
		
		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, "error", response["status"])
		assert.Equal(t, "Not found", response["error"])
	})
	
	t.Run("handler returns generic error", func(t *testing.T) {
		handler := Make(func(w http.ResponseWriter, r *http.Request) error {
			return errors.New("unexpected error")
		})
		
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		
		handler(w, req)
		
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		
		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, "error", response["status"])
		assert.Equal(t, "unexpected error", response["error"])
	})
}
