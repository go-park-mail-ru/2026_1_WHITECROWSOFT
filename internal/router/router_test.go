package router

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestRouter_Integration(t *testing.T) {
	r := New()

	tests := []struct {
		name           string
		method         string
		url            string
		expectedStatus int
	}{
		{"ping", "GET", "/ping", http.StatusOK},
		{"signup_method_check", "GET", "/signup", http.StatusMethodNotAllowed},
		{"protected_no_auth", "GET", "/protected", http.StatusUnauthorized},
		{"not_found", "GET", "/unknown", http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestProtectedEndpoint_Logic(t *testing.T) {
	t.Run("success_with_id", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/protected", nil)
		ctx := context.WithValue(req.Context(), types.UserIDKey, "user-123")
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		TestProtectedEndpoint(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "user-123")
		assert.Contains(t, w.Body.String(), "protected endpoint")
	})

	t.Run("fail_no_id", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/protected", nil)
		w := httptest.NewRecorder()

		TestProtectedEndpoint(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "error")
	})
}
