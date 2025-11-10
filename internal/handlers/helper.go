package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/go-playground/validator/v10"
)

func Render(ctx context.Context, w http.ResponseWriter, comp templ.Component) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return comp.Render(ctx, w)
}

func RenderMultiComponents(ctx context.Context, w http.ResponseWriter, comps []templ.Component) error {
	strBuffer := bytes.NewBufferString("")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	for _, comp := range comps {
		comp.Render(ctx, strBuffer)
	}
	_, err := w.Write(strBuffer.Bytes())
	return err
}

func HTMXRedirect(w http.ResponseWriter, path string) {
	w.Header().Set("HX-Redirect", path)
	w.WriteHeader(http.StatusFound)
}

func SetFlashCookie(w http.ResponseWriter, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:   "flash",
		Value:  value,
		Path:   "/",
		MaxAge: 5,
	})
}

func GetAndDeleteFlashCookie(w http.ResponseWriter, r *http.Request) string {
	var value string
	if cookie, err := r.Cookie("flash"); err == nil {
		value = cookie.Value
		// Clear cookie
		http.SetCookie(w, &http.Cookie{
			Name:   "flash",
			Path:   "/",
			MaxAge: -1,
		})
	}
	return value
}

func ParseInputName(name string, sep string) (string, string) {
	parts := strings.SplitN(name, sep, 2)
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}

func SafeAccessSliceIndex(slice []string, index int) string {
	if index < 0 || index >= len(slice) {
		return ""
	}
	return slice[index]
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

// ParseDateInVietnamTimezone parses a date string (format: "2006-01-02") and converts it to a specific time
// in Vietnamese timezone, then returns it as UTC time.
// The hour, minute, second, and nanosecond parameters allow specifying the exact time of day.
func ParseDateInVietnamTimezone(dateStr string, hour, minute, second, nanosecond int) (*time.Time, error) {
	// Parse the date string
	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, err
	}

	// Load Vietnam timezone
	vietnamLocation, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		vietnamLocation = time.UTC // Fallback to UTC
	}

	// Create time in Vietnam timezone with specified hour, minute, second, nanosecond
	vietnamTime := time.Date(
		parsedDate.Year(),
		parsedDate.Month(),
		parsedDate.Day(),
		hour,
		minute,
		second,
		nanosecond,
		vietnamLocation,
	)

	// Convert to UTC
	utcTime := vietnamTime.UTC()
	return &utcTime, nil
}

// ParseStartOfDayInVietnamTimezone parses a date string and returns the start of day (00:00:00) in UTC
func ParseStartOfDayInVietnamTimezone(dateStr string) (*time.Time, error) {
	return ParseDateInVietnamTimezone(dateStr, 0, 0, 0, 0)
}

// ParseEndOfDayInVietnamTimezone parses a date string and returns the end of day (23:59:59.999999999) in UTC
func ParseEndOfDayInVietnamTimezone(dateStr string) (*time.Time, error) {
	return ParseDateInVietnamTimezone(dateStr, 23, 59, 59, 999999999)
}

func BindAndValidate(r *http.Request, v *validator.Validate, dst interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return &AppError{StatusCode: http.StatusBadRequest, Message: INVALID_REQUEST_PAYLOAD_ERROR}
	}
	if v == nil {
		return nil
	}
	if err := v.Struct(dst); err != nil {
		// map validator errors to a readable message or return AppError
		return &AppError{StatusCode: http.StatusBadRequest, Message: err.Error()}
	}
	return nil
}

func ParseListParams(r *http.Request, defaultPageSize int) models.GenericQueryOptions {
	q := r.URL.Query()

	page := 1
	if p := q.Get("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}

	pageSize := defaultPageSize
	if ps := q.Get("page_size"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && v > 0 {
			pageSize = v
		}
	}

	sortBy := q.Get("sort_by")
	if sortBy == "" {
		sortBy = "created_at"
	}

	sortOrder := q.Get("sort_order")
	if sortOrder == "" {
		sortOrder = "desc"
	}

	return models.GenericQueryOptions{
		Page:      page,
		PageSize:  pageSize,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}
}
