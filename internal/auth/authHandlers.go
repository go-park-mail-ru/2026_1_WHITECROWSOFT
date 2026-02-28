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
