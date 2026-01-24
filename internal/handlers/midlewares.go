package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/datdev2409/lab-admin-go/internal/auth"
	"github.com/datdev2409/lab-admin-go/internal/logger"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func LoggingMiddleware(logObj *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestLogger := logObj.With(zap.String("trace_id", middleware.GetReqID(r.Context())))
			ctx := logger.WithCtx(r.Context(), requestLogger)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// extractToken attempts to extract JWT token from either cookie or Authorization header
func extractToken(r *http.Request) (string, error) {
	// First, try to get token from cookie
	cookie, err := r.Cookie("auth_token")
	if err == nil {
		return cookie.Value, nil
	}

	// If cookie not found, try Authorization header with Bearer scheme
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			return parts[1], nil
		}
	}

	return "", http.ErrNoCookie
}

func JWTAuthAPIEndpoint(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr, err := extractToken(r)
		if err != nil {
			http.Error(w, "Missing or invalid authentication token", http.StatusUnauthorized)
			return
		}

		userId, err := auth.ValidateJWT(tokenStr)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userId", userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func JWTAuthWebEndpoint(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr, err := extractToken(r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		userId, err := auth.ValidateJWT(tokenStr)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), "userId", userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func HTTPLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.FromCtx(r.Context())
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		start := time.Now()
		next.ServeHTTP(ww, r)
		duration := time.Since(start)
		statusCode := ww.Status()

		logLevel := zap.InfoLevel
		if statusCode >= 500 {
			logLevel = zap.ErrorLevel
		} else if statusCode >= 400 {
			logLevel = zap.WarnLevel
		}
		log.Log(logLevel, "HTTP Request", zap.String("method", r.Method), zap.String("url", r.URL.String()), zap.Int("status", ww.Status()), zap.Duration("duration", duration))
	})
}
