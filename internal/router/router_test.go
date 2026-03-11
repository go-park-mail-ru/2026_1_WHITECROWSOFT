package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestRouter_Integration(t *testing.T) {
	jwtCfg := config.JWTConfig{
		Secret:        "test-secret",
		CookieName:    "test-cookie-name",
		CookieTimeJWT: 3600,
		Secure:        false,
	}

	serverCfg := config.ServerConfig{
		Port:            "8000",
		ShutdownTimeout: 5,
	}

	cfg := config.Config{
		JWT:    jwtCfg,
		Server: serverCfg,
	}
	r := New(&cfg)

	tests := []struct {
		name           string
		method         string
		url            string
		expectedStatus int
	}{
		{"ping", http.MethodGet, "/ping", http.StatusOK},
		{"signup_method_check", http.MethodGet, "/signup", http.StatusMethodNotAllowed},
		{"protected_no_auth", http.MethodGet, "/notes", http.StatusUnauthorized},
		{"not_found", http.MethodGet, "/unknown", http.StatusNotFound},
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
