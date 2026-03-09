package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/auth"
	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/storage"
	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/types"
	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/pkg/jwt"
	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	secret := "new-secret"
	authHandler := auth.NewHandler(secret, storage.NewUserSet())

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(types.UserIDKey)
		assert.NotNil(t, userID)
		w.WriteHeader(http.StatusOK)
	})

	t.Run("no_cookie", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()

		Auth(nextHandler, authHandler).ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid_token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.AddCookie(&http.Cookie{Name: auth.CookieName, Value: "no_value_its_a_trick"})
		w := httptest.NewRecorder()

		Auth(nextHandler, authHandler).ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("success", func(t *testing.T) {
		token, _ := jwt.GenerateToken("testuser123", time.Minute, secret)
		req := httptest.NewRequest("GET", "/", nil)
		req.AddCookie(&http.Cookie{Name: auth.CookieName, Value: token})
		w := httptest.NewRecorder()

		Auth(nextHandler, authHandler).ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
