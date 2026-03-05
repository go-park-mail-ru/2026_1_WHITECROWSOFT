package logger

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

type ctxKey string

const (
	RequestIDKey ctxKey = "request_id"
)

type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func (e HTTPError) Error() string {
	return e.Message
}

func Init() *slog.Logger {
	logLevel := os.Getenv("LOG_LEVEL")

	var level slog.Level
	switch strings.ToUpper(logLevel) {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     level,
		AddSource: level == slog.LevelDebug,
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger
}

func WithRequest(ctx context.Context, args ...any) (*slog.Logger, []any) {
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		args = append(args, "request_id", requestID)
	}
	return slog.Default(), args
}

func Debug(ctx context.Context, msg string, args ...any) {
	log, updatedArgs := WithRequest(ctx, args...)
	log.DebugContext(ctx, msg, updatedArgs...)
}

func Info(ctx context.Context, msg string, args ...any) {
	log, updatedArgs := WithRequest(ctx, args...)
	log.InfoContext(ctx, msg, updatedArgs...)
}

func Warn(ctx context.Context, msg string, args ...any) {
	log, updatedArgs := WithRequest(ctx, args...)
	log.WarnContext(ctx, msg, updatedArgs...)
}

func Error(ctx context.Context, msg string, args ...any) {
	log, updatedArgs := WithRequest(ctx, args...)
	log.ErrorContext(ctx, msg, updatedArgs...)
}
