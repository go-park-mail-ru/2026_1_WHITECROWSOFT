package authHandlers

import (
	"context"
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
	ErrUserExist    = errors.New("user already exists")
	ErrUserNotExist = errors.New("user not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrInvalidToken = errors.New("invalid token")
	ErrInternal     = errors.New("internal server error")
	ErrUnauthorized = errors.New("unauthorized")
	validate        = validator.New()
	isSecure        = os.Getenv("IS_SECURE") == "true"
)

type AuthHandler struct {
	jwtSecret string
	userSet   *UserSet
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

func NewAuthHandler(secret string, users *UserSet) *AuthHandler {
	return &AuthHandler{
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

func (a *AuthHandler) SignupUser(w http.ResponseWriter, r *http.Request) {
	var signUpUser dto.SignUpUser

	if err := json.NewDecoder(r.Body).Decode(&signUpUser); err != nil {
		helpers.JSONErrorResponse(w, http.StatusBadRequest, ErrInvalidInput)
		return
	}
	defer r.Body.Close()

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

func (a *AuthHandler) SigninUser(w http.ResponseWriter, r *http.Request) {
	var signInUser dto.SignInUser

	if err := json.NewDecoder(r.Body).Decode(&signInUser); err != nil {
		helpers.JSONErrorResponse(w, http.StatusBadRequest, ErrInvalidInput)
		return
	}
	defer r.Body.Close()

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

func (a *AuthHandler) LogOutUser(w http.ResponseWriter, r *http.Request) {
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

func (a *AuthHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookieJWT, err := r.Cookie(CookieName)
		if err != nil {
			helpers.JSONErrorResponse(w, http.StatusUnauthorized, ErrUnauthorized)
			return
		}

		tokenPayload, err := jwt.ValidateToken(cookieJWT.Value, a.jwtSecret)
		if err != nil {
			helpers.JSONErrorResponse(w, http.StatusUnauthorized, ErrInvalidToken)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", tokenPayload.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *AuthHandler) TestProtectedEndpoint(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		helpers.JSONErrorResponse(w, http.StatusInternalServerError, errors.New("user_id not found in context"))
		return
	}

	response := map[string]string{
		"msg":     "This is a protected endpoint",
		"user_id": userID,
		"time":    time.Now().Format(time.RFC3339),
	}

	helpers.JSONResponse(w, http.StatusOK, response)
}
