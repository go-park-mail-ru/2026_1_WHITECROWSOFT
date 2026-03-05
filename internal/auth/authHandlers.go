package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/dto"
	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/storage"
	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/pkg/helpers"
	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/pkg/jwt"

	"github.com/go-playground/validator/v10"
)

const (
	CookieName    = "NoterianCookieJWT"
	CookieTimeJWT = time.Hour
)

var (
	ErrInvalidInput     = errors.New("invalid input")
	ErrInvalidToken     = errors.New("invalid token")
	ErrInternal         = errors.New("internal server error")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrMethodNotAllowed = errors.New("method not allowed")
	ErrTokenCreation    = errors.New("failed to create token")
	ErrBadCredentials   = errors.New("incorrect username or password")
	validate            = validator.New()
	isSecure            = os.Getenv("IS_SECURE") == "true"
)

type Handler struct {
	jwtSecret string
	users     *storage.UserSet
}

func (a *Handler) Secret() string {
	// NOTE: I added this method to get middleware.Auth working,
	// even though I proposed making jwtSecret private
	// in the first placee. How do we get around this? -Andrew
	return a.jwtSecret
}

type UserResponse struct {
	ID    string `json:"id"`
	Login string `json:"login"`
}

func NewHandler(secret string, users *storage.UserSet) *Handler {
	return &Handler{
		jwtSecret: secret,
		users:     users,
	}
}

func (a *Handler) SignupUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" || r.Body == nil {
		helpers.JSONErrorResponse(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
	}
	defer r.Body.Close()

	var signUpUser dto.SignUpUser

	if err := json.NewDecoder(r.Body).Decode(&signUpUser); err != nil {
		helpers.JSONErrorResponse(w, http.StatusBadRequest, ErrInvalidInput)
		return
	}

	if err := validate.Struct(signUpUser); err != nil {
		helpers.JSONErrorResponse(w, http.StatusBadRequest, ErrInvalidInput)
		return
	}

	user, err := a.users.CreateUser(signUpUser.Login, signUpUser.Password)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrUserExist):
			helpers.JSONErrorResponse(w, http.StatusConflict, storage.ErrUserExist)
		default:
			helpers.JSONErrorResponse(w, http.StatusInternalServerError, ErrInternal)
		}
		return
	}

	tokenStr, err := jwt.GenerateToken(user.ID.String(), CookieTimeJWT, a.jwtSecret)
	if err != nil {
		helpers.JSONErrorResponse(w, http.StatusInternalServerError, ErrTokenCreation)
		return
	}

	cookie := &http.Cookie{
		Name:     CookieName,
		Value:    tokenStr,
		HttpOnly: true,
		Secure:   isSecure,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(CookieTimeJWT),
		Path:     "/",
	}
	http.SetCookie(w, cookie)

	resp := UserResponse{
		ID:    user.ID.String(),
		Login: user.Username,
	}

	helpers.JSONResponse(w, http.StatusOK, resp)
}

func (a *Handler) SigninUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" || r.Body == nil {
		helpers.JSONErrorResponse(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
	}
	defer r.Body.Close()

	var signInUser dto.SignInUser

	if err := json.NewDecoder(r.Body).Decode(&signInUser); err != nil {
		helpers.JSONErrorResponse(w, http.StatusBadRequest, ErrInvalidInput)
		return
	}

	if err := validate.Struct(signInUser); err != nil {
		helpers.JSONErrorResponse(w, http.StatusBadRequest, ErrInvalidInput)
		return
	}

	user, err := a.users.ValidateUser(signInUser.Login, signInUser.Password)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrUserNotExist):
			helpers.JSONErrorResponse(w, http.StatusUnauthorized, ErrBadCredentials)
		default:
			helpers.JSONErrorResponse(w, http.StatusInternalServerError, ErrInternal)
		}
		return
	}

	tokenStr, err := jwt.GenerateToken(user.ID.String(), CookieTimeJWT, a.jwtSecret)
	if err != nil {
		helpers.JSONResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to create token",
		})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    tokenStr,
		HttpOnly: true,
		Secure:   isSecure,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(CookieTimeJWT.Seconds()),
		Path:     "/",
	})

	helpers.JSONResponse(w, http.StatusOK, UserResponse{
		ID:    user.ID.String(),
		Login: user.Username,
	})
}

func (a *Handler) LogOutUser(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    "",
		HttpOnly: true,
		Secure:   isSecure,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
		Path:     "/",
	})

	w.WriteHeader(http.StatusNoContent)
}
