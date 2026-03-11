package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserSet(t *testing.T) {
	s := NewUserSet()
	login := "test123"
	pass := "password123"

	t.Run("createUser success", func(t *testing.T) {
		user, err := s.CreateUser(login, pass)
		require.NoError(t, err)
		assert.Equal(t, login, user.Username)
		assert.NotEmpty(t, user.ID)
	})

	t.Run("createUser alreadyExists", func(t *testing.T) {
		_, err := s.CreateUser(login, "password")
		assert.ErrorIs(t, err, ErrUserExist)
	})

	t.Run("validateUser success", func(t *testing.T) {
		user, err := s.ValidateUser(login, pass)
		require.NoError(t, err)
		assert.Equal(t, login, user.Username)
	})

	t.Run("validateUser wrongPassword", func(t *testing.T) {
		_, err := s.ValidateUser(login, "wrongPassword")
		assert.ErrorIs(t, err, ErrUserNotExist)
	})

	t.Run("validateUser notfound", func(t *testing.T) {
		_, err := s.ValidateUser("newUuser", pass)
		assert.ErrorIs(t, err, ErrUserNotExist)
	})
}
