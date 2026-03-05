package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/dto"
	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/models"
	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/pkg/helpers"
	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/pkg/jwt"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	CookieName    = "NoterianCookieJWT"
	CookieTimeJWT = time.Hour
)

var (
	ErrUserExist        = errors.New("user already exists")
	ErrUserNotExist     = errors.New("user not found")
	ErrInvalidInput     = errors.New("invalid input")
	ErrInvalidToken     = errors.New("invalid token")
	ErrInternal         = errors.New("internal server error")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrMethodNotAllowed = errors.New("method not allowed")
	validate            = validator.New()
	isSecure            = os.Getenv("IS_SECURE") == "true"
)

type Handler struct {
	jwtSecret string
	userSet   *UserSet
}

func (a *Handler) Secret() string {
	// NOTE: I added this method to get middleware.Auth working,
	// even though I proposed making jwtSecret private
	// in the first placee. How do we get around this? -Andrew
	return a.jwtSecret
}

type UserSet struct {
	users map[string]*models.User
	mu    sync.RWMutex
}

type UserResponse struct {
	ID    string `json:"id"`
	Login string `json:"login"`
	Token string `json:"token"`
}

func NewHandler(secret string, users *UserSet) *Handler {
	return &Handler{
		jwtSecret: secret,
		userSet:   users,
	}
}

func NewUserSet() *UserSet {
	return &UserSet{
		users: make(map[string]*models.User),
		mu:    sync.RWMutex{},
	}
}

func (s *UserSet) CreateUser(login, password string) (*models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[login]; exists {
		return nil, ErrUserExist
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:       uuid.New(),
		Username: login,
		Password: hashPassword,
	}

	s.users[login] = user
	return user, nil
}

func (s *UserSet) ValidateUser(login, password string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[login]
	if !exists {
		return nil, ErrUserNotExist
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, ErrUserNotExist
	}

	return user, nil
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

	user, err := a.userSet.CreateUser(signUpUser.Login, signUpUser.Password)
	if err != nil {
		switch {
		case errors.Is(err, ErrUserExist):
			helpers.JSONErrorResponse(w, http.StatusConflict, ErrUserExist)
		default:
			helpers.JSONErrorResponse(w, http.StatusInternalServerError, ErrInternal)
		}
		return
	}

	tokenStr, err := jwt.GenerateToken(user.ID.String(), CookieTimeJWT, a.jwtSecret)
	if err != nil {
		helpers.JSONErrorResponse(w, http.StatusInternalServerError, errors.New("failed to create token"))
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
		Token: tokenStr,
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

	user, err := a.userSet.ValidateUser(signInUser.Login, signInUser.Password)
	if err != nil {
		switch {
		case errors.Is(err, ErrUserNotExist):
			helpers.JSONErrorResponse(w, http.StatusUnauthorized, errors.New("incorrect username or password"))
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
		Expires:  time.Now().Add(CookieTimeJWT),
		Path:     "/",
	})

	helpers.JSONResponse(w, http.StatusOK, UserResponse{
		ID:    user.ID.String(),
		Login: user.Username,
		Token: tokenStr,
	})
}

func (a *Handler) LogOutUser(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    "",
		HttpOnly: true,
		Secure:   isSecure,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(-CookieTimeJWT),
		Path:     "/",
	})

	w.WriteHeader(http.StatusNoContent)
}
