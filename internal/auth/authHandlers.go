package authHandlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"
	"wcs/internal/dto"

	"github.com/go-playground/validator/v10"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

const (
	CookieName    = "NoterianCookieJWT"
	CookieTimeJWT = time.Hour
)

var (
	ErrUserExists         = errors.New("user already exists")
	ErrUserNotExists      = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid username or password")
	validate              = validator.New()
)

type AuthHandler struct {
	JWTSecret string
	UserSet   *UserSet
}

type User struct {
	ID       string
	Login    string
	Password string
}

type UserSet struct {
	users map[string]*User
	mu    sync.RWMutex
}

type UserResponse struct {
	ID    string `json:"id"`
	Login string `json:"login"`
	Token string `json:"token"`
}

func NewUserSet() *UserSet {
	return &UserSet{
		users: make(map[string]*User),
		mu:    sync.RWMutex{},
	}
}

func (s *UserSet) CreateUser(login, password string) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[login]; exists {
		return nil, ErrUserExists
	}

	user := &User{
		ID:       uuid.New().String(),
		Login:    login,
		Password: password,
	}

	s.users[login] = user
	return user, nil
}

func (s *UserSet) ValidateUser(login, password string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[login]
	if !exists {
		return nil, ErrUserNotExists
	}

	if password != user.Password {
		return nil, ErrInvalidCredentials
	}
	return user, nil
}

func WriteResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (a *AuthHandler) SignupUser(w http.ResponseWriter, r *http.Request) {
	var signUpUser dto.SignUpUser

	if err := json.NewDecoder(r.Body).Decode(&signUpUser); err != nil {
		WriteResponse(w, http.StatusBadRequest, map[string]string{
			"error": "invalid input",
		})
		return
	}
	defer r.Body.Close()

	if err := validate.Struct(signUpUser); err != nil {
		WriteResponse(w, http.StatusBadRequest, map[string]string{
			"error":   "validation failed",
			"details": err.Error(),
		})
		return
	}

	user, err := a.UserSet.CreateUser(signUpUser.Login, signUpUser.Password)
	if err != nil {
		WriteResponse(w, http.StatusConflict, map[string]string{
			"error":   "failed to create user",
			"details": err.Error(),
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(CookieTimeJWT).Unix(),
	})

	tokenStr, err := token.SignedString([]byte(a.JWTSecret))
	if err != nil {
		WriteResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to create token",
		})
		return
	}

	cookie := &http.Cookie{
		Name:  CookieName,
		Value: tokenStr,
		//HttpOnly: true,
		//Secure:  true,
		Expires: time.Now().Add(CookieTimeJWT),
		Path:    "/",
	}
	http.SetCookie(w, cookie)

	resp := UserResponse{
		ID:    user.ID,
		Login: user.Login,
		Token: tokenStr,
	}
	WriteResponse(w, http.StatusOK, resp)
}

func (a *AuthHandler) SigninUser(w http.ResponseWriter, r *http.Request) {
	var signInUser dto.SignInUser

	if err := json.NewDecoder(r.Body).Decode(&signInUser); err != nil {
		WriteResponse(w, http.StatusBadRequest, map[string]string{
			"error": "invalid input",
		})
		return
	}
	defer r.Body.Close()

	if err := validate.Struct(signInUser); err != nil {
		WriteResponse(w, http.StatusBadRequest, map[string]string{
			"error":   "validation failed",
			"details": err.Error(),
		})
		return
	}

	user, err := a.UserSet.ValidateUser(signInUser.Login, signInUser.Password)
	if err != nil {
		WriteResponse(w, http.StatusUnauthorized, map[string]string{
			"error":   "invalid username or password",
			"details": err.Error(),
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour).Unix(),
	})

	tokenStr, err := token.SignedString([]byte(a.JWTSecret))
	if err != nil {
		WriteResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to create token",
		})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  CookieName,
		Value: tokenStr,
		//HttpOnly: true,
		//Secure: true,
		Expires: time.Now().Add(CookieTimeJWT),
		Path:    "/",
	})

	WriteResponse(w, http.StatusOK, UserResponse{
		ID:    user.ID,
		Login: user.Login,
		Token: tokenStr,
	})
}
