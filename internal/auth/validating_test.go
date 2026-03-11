package auth

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestLoginValidate(t *testing.T) {
	type LoginStruct struct {
		Login string `validate:"login"`
	}

	validate := validator.New()
	err := validate.RegisterValidation("login", validateLogin)
	assert.NoError(t, err, "Failed to register validation")

	tests := []struct {
		name    string
		login   string
		wantErr bool
	}{
		{
			name:    "simple login",
			login:   "john",
			wantErr: false,
		},
		{
			name:    "login with numbers",
			login:   "john123",
			wantErr: false,
		},
		{
			name:    "login with underscore",
			login:   "john_doe",
			wantErr: false,
		},
		{
			name:    "login with dot",
			login:   "john.doe",
			wantErr: false,
		},
		{
			name:    "login mixed",
			login:   "john.doe_123",
			wantErr: false,
		},
		{
			name:    "login single character",
			login:   "a",
			wantErr: false,
		},
		{
			name:    "login mixed case",
			login:   "JohnDoe",
			wantErr: false,
		},
		{
			name:    "login with russian characteres",
			login:   "Андрей_228",
			wantErr: false,
		},
		{
			name:    "Invalid with dash",
			login:   "john-doe",
			wantErr: true,
		},
		{
			name:    "Invalid with space",
			login:   "john doe",
			wantErr: true,
		},
		{
			name:    "login starts with underscore",
			login:   "_john",
			wantErr: true,
		},
		{
			name:    "login starts with dot",
			login:   ".john",
			wantErr: true,
		},
		{
			name:    "login with multiple underscores",
			login:   "__john",
			wantErr: true,
		},
		{
			name:    "login with multiple dots",
			login:   "..john",
			wantErr: true,
		},
		{
			name:    "login with underscore and dot",
			login:   "_.john",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := LoginStruct{Login: tt.login}
			err := validate.Struct(s)

			if tt.wantErr {
				assert.Error(t, err, "Expected validation error for login: %s", tt.login)
			} else {
				assert.NoError(t, err, "Expected no validation error for login: %s", tt.login)
			}
		})
	}
}

func TestPasswordValidate(t *testing.T) {
	type PasswordStruct struct {
		Password string `validate:"password"`
	}

	validate := validator.New()
	err := validate.RegisterValidation("password", validatePassword)
	assert.NoError(t, err, "Failed to register validation")

	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid - basic",
			password: "Password1",
			wantErr:  false,
		},
		{
			name:     "valid - with special chars",
			password: "P@ssw0rd!",
			wantErr:  false,
		},
		{
			name:     "valid - with russian",
			password: "Пароль123A",
			wantErr:  false,
		},
		{
			name:     "valid - minimum length",
			password: "Pass1",
			wantErr:  false,
		},
		{
			name:     "invalid - too short (4 chars)",
			password: "Ab1",
			wantErr:  true,
		},
		{
			name:     "invalid - empty",
			password: "",
			wantErr:  true,
		},
		{
			name:     "invalid - no uppercase",
			password: "password1",
			wantErr:  true,
		},
		{
			name:     "invalid - no uppercase russian",
			password: "пароль123",
			wantErr:  true,
		},
		{
			name:     "invalid - no digits",
			password: "Password",
			wantErr:  true,
		},
		{
			name:     "invalid - no digits russian",
			password: "ПарольA",
			wantErr:  true,
		},
		{
			name:     "invalid - no uppercase and no digits",
			password: "password",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := PasswordStruct{Password: tt.password}
			err := validate.Struct(s)

			if tt.wantErr {
				assert.Error(t, err, "Expected validation error for password: %s", tt.password)
			} else {
				assert.NoError(t, err, "Expected no validation error for password: %s", tt.password)
			}
		})
	}
}
