package handlers

import (
	"net/http"
	"time"

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
