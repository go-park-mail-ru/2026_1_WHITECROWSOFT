package authHandlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"sync"
	"time"
	"wcs/internal/dto"
	"wcs/internal/models"
	"wcs/pkg/helpers"
	"wcs/pkg/jwt"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	CookieName    = "NoterianCookieJWT"
	CookieTimeJWT = time.Hour
)

var (
	ErrUserExists    = errors.New("user already exists")
	ErrUserNotExists = errors.New("user not found")
	validate         = validator.New()
	isSecure         = os.Getenv("IS_SECURE") == "true"
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
		return nil, ErrUserExists
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
		return nil, ErrUserNotExists
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, ErrUserNotExists
	}

	return user, nil
}

func (a *AuthHandler) SignupUser(w http.ResponseWriter, r *http.Request) {
	var signUpUser dto.SignUpUser

	if err := json.NewDecoder(r.Body).Decode(&signUpUser); err != nil {
		helpers.JSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": "invalid input",
		})
		return
	}
	defer r.Body.Close()

	if err := validate.Struct(signUpUser); err != nil {
		helpers.JSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": "validation failed",
		})
		return
	}

	user, err := a.userSet.CreateUser(signUpUser.Login, signUpUser.Password)
	if err != nil {
		switch {
		case errors.Is(err, ErrUserExists):
			helpers.JSONResponse(w, http.StatusConflict, map[string]string{
				"error": "user already exists",
			})
		default:
			helpers.JSONResponse(w, http.StatusInternalServerError, map[string]string{
				"error": "internal server error",
			})
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

	cookie := &http.Cookie{
		Name:     CookieName,
		Value:    tokenStr,
		HttpOnly: true,
		Secure:   isSecure,
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
		helpers.JSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": "invalid input",
		})
		return
	}
	defer r.Body.Close()

	if err := validate.Struct(signInUser); err != nil {
		helpers.JSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": "validation failed",
		})
		return
	}

	user, err := a.userSet.ValidateUser(signInUser.Login, signInUser.Password)
	if err != nil {
		switch {
		case errors.Is(err, ErrUserNotExists):
			helpers.JSONResponse(w, http.StatusUnauthorized, map[string]string{
				"error": "incorrect username or password",
			})
		default:
			helpers.JSONResponse(w, http.StatusInternalServerError, map[string]string{
				"error": "internal server error",
			})
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
		Expires:  time.Now().Add(-CookieTimeJWT),
		Path:     "/",
	})

	w.WriteHeader(http.StatusNoContent)
}

func (a *AuthHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookieJWT, err := r.Cookie(CookieName)
		if err != nil {
			helpers.JSONResponse(w, http.StatusUnauthorized, map[string]string{
				"error": "unathorized",
			})
			return
		}

		tokenPayload, err := jwt.ValidateToken(cookieJWT.Value, a.jwtSecret)
		if err != nil {
			helpers.JSONResponse(w, http.StatusUnauthorized, map[string]string{
				"error": "invalid token",
			})
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", tokenPayload.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *AuthHandler) TestProtectedEndpoint(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		helpers.JSONResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "user_id not found in context",
		})
		return
	}

	response := map[string]string{
		"msg":     "This is a protected endpoint",
		"user_id": userID,
		"time":    time.Now().Format(time.RFC3339),
	}

	helpers.JSONResponse(w, http.StatusOK, response)
}
