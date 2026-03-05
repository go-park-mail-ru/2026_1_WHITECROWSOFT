package middleware

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/auth"
	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/pkg/helpers"
	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/pkg/jwt"
)

func Auth(next http.Handler, a *auth.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookieJWT, err := r.Cookie(auth.CookieName)
		if err != nil {
			helpers.JSONErrorResponse(w, http.StatusUnauthorized, auth.ErrUnauthorized)
			return
		}

		tokenPayload, err := jwt.ValidateToken(cookieJWT.Value, a.Secret())
		if err != nil {
			helpers.JSONErrorResponse(w, http.StatusUnauthorized, auth.ErrInvalidToken)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", tokenPayload.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
