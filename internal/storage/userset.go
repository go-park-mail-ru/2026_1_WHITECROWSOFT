package storage

import (
	"errors"
	"sync"

	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserExist    = errors.New("user already exists")
	ErrUserNotExist = errors.New("user not found")
)

type UserSet struct {
	users map[string]*models.User
	mu    sync.RWMutex
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
