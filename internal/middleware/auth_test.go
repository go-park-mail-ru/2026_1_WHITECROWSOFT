package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/config"
	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/types"
	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/pkg/jwt"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	cfg := config.JWTConfig{
		CookieName:    "test_session_token",
		Secret:        "a_highly_guarded_secret_326",
		CookieTimeJWT: time.Duration(time.Hour),
		Secure:        true,
	}

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(types.UserIDKey)
		assert.NotNil(t, userID, "UserID not found in context")
		w.WriteHeader(http.StatusOK)
	})

	t.Run("401 if cookie missing", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		handler := Auth(nextHandler, cfg)
		handler.ServeHTTP(rr, req)

		assert.Equalf(t, rr.Code, http.StatusUnauthorized,
			"expected 401, got %d", rr.Code,
		)
	})

	t.Run("proceed if valid", func(t *testing.T) {
		validToken, err := jwt.GenerateToken("river_wyles", cfg.CookieTimeJWT, cfg.Secret)
		assert.NoErrorf(t, err, "couldn't generate token: %q", err)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{Name: cfg.CookieName, Value: validToken})

		rr := httptest.NewRecorder()

		handler := Auth(nextHandler, cfg)
		handler.ServeHTTP(rr, req)

		assert.Equalf(t, rr.Code, http.StatusOK,
			"expected 200, got %d", rr.Code,
		)
	})
}
