package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/types"
	"github.com/google/uuid"
)

func generateRequestID() string {
	return uuid.New().String()
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		requestID := generateRequestID()

		ctx := context.WithValue(r.Context(), types.RequestIDKey, requestID)
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
