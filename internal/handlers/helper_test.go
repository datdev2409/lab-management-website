package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteJSON(t *testing.T) {
	tests := []struct {
		name       string
		status     int
		data       interface{}
		wantStatus int
	}{
		{
			name:       "write simple object",
			status:     http.StatusOK,
			data:       map[string]string{"message": "success"},
			wantStatus: http.StatusOK,
		},
		{
			name:       "write array",
			status:     http.StatusOK,
			data:       []string{"item1", "item2"},
			wantStatus: http.StatusOK,
		},
		{
			name:       "write with created status",
			status:     http.StatusCreated,
			data:       map[string]int{"id": 123},
			wantStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			err := WriteJSON(w, tt.status, tt.data)
			
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			
			// Verify JSON is valid
			var result interface{}
			err = json.NewDecoder(w.Body).Decode(&result)
			require.NoError(t, err)
		})
	}
}

func TestSetFlashCookie(t *testing.T) {
	w := httptest.NewRecorder()
	value := "Success message"
	
	SetFlashCookie(w, value)
	
	cookies := w.Result().Cookies()
	require.Len(t, cookies, 1)
	
	cookie := cookies[0]
	assert.Equal(t, "flash", cookie.Name)
	assert.Equal(t, value, cookie.Value)
	assert.Equal(t, "/", cookie.Path)
	assert.Equal(t, 5, cookie.MaxAge)
}

func TestGetAndDeleteFlashCookie(t *testing.T) {
	t.Run("cookie exists", func(t *testing.T) {
		// Create request with flash cookie
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{
			Name:  "flash",
			Value: "Test message",
		})
		
		w := httptest.NewRecorder()
		value := GetAndDeleteFlashCookie(w, req)
		
		assert.Equal(t, "Test message", value)
		
		// Verify cookie is deleted (MaxAge = -1)
		cookies := w.Result().Cookies()
		require.Len(t, cookies, 1)
		assert.Equal(t, -1, cookies[0].MaxAge)
	})

	t.Run("cookie does not exist", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		
		value := GetAndDeleteFlashCookie(w, req)
		assert.Empty(t, value)
	})
}

func TestParseInputName(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		separator string
		want1     string
		want2     string
	}{
		{
			name:      "valid input with separator",
			input:     "field_name_123",
			separator: "_",
			want1:     "field",
			want2:     "name_123",
		},
		{
			name:      "input without separator",
			input:     "fieldname",
			separator: "_",
			want1:     "",
			want2:     "",
		},
		{
			name:      "input with multiple separators",
			input:     "field:name:value",
			separator: ":",
			want1:     "field",
			want2:     "name:value",
		},
		{
			name:      "empty input",
			input:     "",
			separator: "_",
			want1:     "",
			want2:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1, got2 := ParseInputName(tt.input, tt.separator)
			assert.Equal(t, tt.want1, got1)
			assert.Equal(t, tt.want2, got2)
		})
	}
}

func TestSafeAccessSliceIndex(t *testing.T) {
	tests := []struct {
		name  string
		slice []string
		index int
		want  string
	}{
		{
			name:  "valid index",
			slice: []string{"a", "b", "c"},
			index: 1,
			want:  "b",
		},
		{
			name:  "first index",
			slice: []string{"a", "b", "c"},
			index: 0,
			want:  "a",
		},
		{
			name:  "last index",
			slice: []string{"a", "b", "c"},
			index: 2,
			want:  "c",
		},
		{
			name:  "negative index",
			slice: []string{"a", "b", "c"},
			index: -1,
			want:  "",
		},
		{
			name:  "index out of bounds",
			slice: []string{"a", "b", "c"},
			index: 5,
			want:  "",
		},
		{
			name:  "empty slice",
			slice: []string{},
			index: 0,
			want:  "",
		},
		{
			name:  "nil slice",
			slice: nil,
			index: 0,
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SafeAccessSliceIndex(tt.slice, tt.index)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestHTMXRedirect(t *testing.T) {
	w := httptest.NewRecorder()
	path := "/dashboard"
	
	HTMXRedirect(w, path)
	
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, path, w.Header().Get("HX-Redirect"))
}

func TestParseDateInVietnamTimezone(t *testing.T) {
	t.Run("valid date with specific time", func(t *testing.T) {
		result, err := ParseDateInVietnamTimezone("2024-01-15", 10, 30, 45, 0)
		
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 2024, result.Year())
		assert.Equal(t, time.January, result.Month())
		assert.Equal(t, 15, result.Day())
	})

	t.Run("invalid date format", func(t *testing.T) {
		_, err := ParseDateInVietnamTimezone("invalid-date", 0, 0, 0, 0)
		assert.Error(t, err)
	})

	t.Run("midnight time", func(t *testing.T) {
		result, err := ParseDateInVietnamTimezone("2024-12-25", 0, 0, 0, 0)
		
		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("end of day time", func(t *testing.T) {
		result, err := ParseDateInVietnamTimezone("2024-12-25", 23, 59, 59, 999999999)
		
		require.NoError(t, err)
		assert.NotNil(t, result)
	})
}

func TestParseStartOfDayInVietnamTimezone(t *testing.T) {
	t.Run("valid date", func(t *testing.T) {
		result, err := ParseStartOfDayInVietnamTimezone("2024-06-15")
		
		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("invalid date", func(t *testing.T) {
		_, err := ParseStartOfDayInVietnamTimezone("not-a-date")
		assert.Error(t, err)
	})
}

func TestParseEndOfDayInVietnamTimezone(t *testing.T) {
	t.Run("valid date", func(t *testing.T) {
		result, err := ParseEndOfDayInVietnamTimezone("2024-06-15")
		
		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("invalid date", func(t *testing.T) {
		_, err := ParseEndOfDayInVietnamTimezone("invalid")
		assert.Error(t, err)
	})

	t.Run("start and end of same day difference", func(t *testing.T) {
		start, err1 := ParseStartOfDayInVietnamTimezone("2024-01-01")
		end, err2 := ParseEndOfDayInVietnamTimezone("2024-01-01")
		
		require.NoError(t, err1)
		require.NoError(t, err2)
		
		// End should be after start
		assert.True(t, end.After(*start))
	})
}
