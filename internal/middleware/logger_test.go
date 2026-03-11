package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestLoggerMiddleware(t *testing.T) {
	status := http.StatusTeapot
	msg := []byte("I'm a teapot; And to me, time is a place.")
	var requestId any
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId = r.Context().Value(types.RequestIDKey)

		w.WriteHeader(status)
		w.Write(msg)
	})

	handler := Logger(nextHandler)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, status, rr.Code, "incorrect code")
	assert.Equal(t, msg, rr.Body.Bytes(), "incorrect body")

	assert.NotNil(t, requestId, "requestID not present")
	uuidStr, ok := requestId.(string)
	assert.True(t, ok, "requestID not a string")
	_, err := uuid.Parse(uuidStr)
	assert.NoError(t, err, "couldn't parse requestID as UUID")
}
