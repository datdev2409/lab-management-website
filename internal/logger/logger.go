package logger

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type loggerCtxKey struct{}

func Init() *zap.Logger {
	logOptions := []zap.Option{zap.AddCaller(), zap.AddStacktrace(zap.PanicLevel)}
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	if os.Getenv("GO_ENV") == "local" {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	l := zap.Must(cfg.Build(logOptions...))
	l.Info("Logger initialized", zap.String("env", os.Getenv("GO_ENV")))
	return l
}

func FromCtx(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(loggerCtxKey{}).(*zap.Logger); ok {
		return l
	}
	return zap.L()
}

func WithCtx(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey{}, l)
}
