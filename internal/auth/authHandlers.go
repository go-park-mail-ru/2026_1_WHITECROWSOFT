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

const CookieName = "MinoCookieJWT"

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

var ErrUserExists = errors.New("user already exists")
var validate = validator.New()

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

func (a *AuthHandler) SignupUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var input dto.SignUpUser
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid input"})
		return
	}
	defer r.Body.Close()

	if err := validate.Struct(input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "validation failed",
			"details": err.Error(),
		})
		return
	}

	user, err := a.UserSet.CreateUser(input.Login, input.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to create user"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"login":   user.Login,
		"exp":     time.Now().Add(time.Hour).Unix(),
	})

	tokenStr, err := token.SignedString([]byte(a.JWTSecret))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to create user"})
		return
	}

	cookie := &http.Cookie{
		Name:  CookieName,
		Value: tokenStr,
		//HttpOnly: true,
		Secure:  true,
		Expires: time.Now().Add(time.Hour),
		Path:    "/",
	}
	http.SetCookie(w, cookie)

	resp := UserResponse{
		ID:    user.ID,
		Login: user.Login,
		Token: tokenStr,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
