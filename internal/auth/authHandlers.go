package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/config"
	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/dto"
	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/models"
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
	ErrInternal         = errors.New("internal server error")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrMethodNotAllowed = errors.New("method not allowed")
	ErrBadCredentials   = errors.New("Неверный логин или пароль!")
	validate            = validator.New()
	isSecure            = os.Getenv("IS_SECURE") == "true"
)

type UserRepository interface {
	CreateUser(login, password string) (*models.Account, error)
	ValidateUser(login, password string) (*models.Account, error)
}

type Handler struct {
	jwtConfig config.JWTConfig
	users     UserRepository
}

type UserResponse struct {
	ID    string `json:"id"`
	Login string `json:"login"`
}

func NewHandler(jwtConfig config.JWTConfig, users UserRepository) *Handler {
	return &Handler{
		jwtConfig: jwtConfig,
		users:     users,
	}
}

func getFromBody[T dto.SignInUser | dto.SignUpUser](r *http.Request, u *T) error {
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		return err
	}
	return validate.Struct(u)
}

func (a *Handler) saveUserCookie(w http.ResponseWriter, user *models.Account) {
	tokenStr, err := jwt.GenerateToken(user.ID.String(), CookieTimeJWT, a.jwtConfig.Secret)
	if err != nil {
		helpers.JSONErrorResponse(w, http.StatusInternalServerError, jwt.ErrTokenCreation)
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

func (a *Handler) SignupUser(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		helpers.JSONErrorResponse(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
	}
	defer r.Body.Close()

	var signUpUser dto.SignUpUser

	if err := getFromBody(r, &signUpUser); err != nil {
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

	a.saveUserCookie(w, user)
}

func (a *Handler) SigninUser(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		helpers.JSONErrorResponse(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
	}
	defer r.Body.Close()

	var signInUser dto.SignInUser

	if err := getFromBody(r, &signInUser); err != nil {
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

	a.saveUserCookie(w, user)
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
