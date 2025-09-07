package utils

import (
	"context"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

// Global logger instance
var logger *zap.Logger

// InitLogger initializes the global zap logger
func InitLogger(development bool) error {
	var err error
	if development {
		logger, err = zap.NewDevelopment(zap.AddCallerSkip(1))
	} else {
		logger, err = zap.NewProduction(zap.AddCallerSkip(1))
	}
	if err != nil {
		return err
	}
	return nil
}

// GetLogger returns the global logger instance
func GetLogger() *zap.Logger {
	if logger == nil {
		// Fallback to development logger if not initialized
		logger, _ = zap.NewDevelopment(zap.AddCallerSkip(1))
	}
	return logger
}

// LoggerWithRequestID creates a logger that includes the request ID from context
func LoggerWithRequestID(ctx context.Context) *zap.Logger {
	reqID := middleware.GetReqID(ctx)
	baseLogger := GetLogger()
	if reqID != "" {
		return baseLogger.With(zap.String("requestId", reqID))
	}
	return baseLogger
}

// LogError logs an error with request context
func LogError(ctx context.Context, msg string, fields ...zap.Field) {
	logger := LoggerWithRequestID(ctx)
	logger.Error(msg, fields...)
}

// LogInfo logs an info message with request context
func LogInfo(ctx context.Context, msg string, fields ...zap.Field) {
	logger := LoggerWithRequestID(ctx)
	logger.Info(msg, fields...)
}

// LogWarn logs a warning message with request context
func LogWarn(ctx context.Context, msg string, fields ...zap.Field) {
	logger := LoggerWithRequestID(ctx)
	logger.Warn(msg, fields...)
}

// LogDebug logs a debug message with request context
func LogDebug(ctx context.Context, msg string, fields ...zap.Field) {
	logger := LoggerWithRequestID(ctx)
	logger.Debug(msg, fields...)
}

// Cleanup flushes any buffered log entries
func Cleanup() {
	if logger != nil {
		logger.Sync()
	}
}
