package logger

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ctxKey string

const (
	requestIDKey ctxKey = "request_id"
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

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		requestID := generateRequestID()

		ctx := context.WithValue(r.Context(), requestIDKey, requestID)
		r = r.WithContext(ctx)

		slog.InfoContext(ctx, "request started",
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
			"request_id", requestID,
		)

		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(rw, r)

		slog.InfoContext(ctx, "request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rw.statusCode,
			"duration_ms", time.Since(start).Milliseconds(),
			"request_id", requestID,
		)
	})
}

func WithRequest(ctx context.Context, args ...any) (*slog.Logger, []any) {
	if requestID, ok := ctx.Value(requestIDKey).(string); ok {
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

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func generateRequestID() string {
	return uuid.New().String()
}
