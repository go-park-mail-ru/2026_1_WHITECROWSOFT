package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/config"
	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/dto"
	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTest() (*Handler, *storage.UserSet) {
	userStorage := storage.NewUserSet()
	handler := NewHandler(
		config.JWTConfig{
			Secret:        "secret-for-testing",
			CookieName:    "test-coockie-name",
			CookieTimeJWT: 3600,
			Secure:        false,
		}, userStorage,
	)
	return handler, userStorage
}

func TestSignupUser(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           interface{}
		setupStorage   func(*storage.UserSet)
		expectedStatus int
		expectedError  string
		checkCookie    bool
	}{
		{
			name:   "success signup",
			method: http.MethodPost,
			body: dto.SignUpUser{
				Login:    "test123",
				Password: "password123",
			},
			setupStorage:   func(s *storage.UserSet) {},
			expectedStatus: http.StatusOK,
			expectedError:  "",
			checkCookie:    true,
		},
		{
			method: http.MethodPost,
			body: dto.SignUpUser{
				Login:    "test123",
				Password: "",
			},
			setupStorage:   func(s *storage.UserSet) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  ErrInvalidInput.Error(),
			checkCookie:    false,
		},
		{
			name:   "empty login",
			method: http.MethodPost,
			body: dto.SignUpUser{
				Login:    "",
				Password: "password123",
			},
			setupStorage:   func(s *storage.UserSet) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  ErrInvalidInput.Error(),
			checkCookie:    false,
		},
		{
			name:           "invalid json",
			method:         http.MethodPost,
			body:           "invalid json string",
			setupStorage:   func(s *storage.UserSet) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  ErrInvalidInput.Error(),
			checkCookie:    false,
		},
		{
			name:   "user already exists",
			method: http.MethodPost,
			body: dto.SignUpUser{
				Login:    "test123",
				Password: "password123",
			},
			setupStorage: func(s *storage.UserSet) {
				_, err := s.CreateUser("test123", "password123")
				require.NoError(t, err)
			},
			expectedStatus: http.StatusConflict,
			expectedError:  storage.ErrUserExist.Error(),
			checkCookie:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authHandler, userStorage := setupTest()
			tt.setupStorage(userStorage)

			jsonBody, err := json.Marshal(tt.body)
			require.NoError(t, err)
			r := httptest.NewRequest(tt.method, "/signup", bytes.NewBuffer(jsonBody))
			w := httptest.NewRecorder()

			authHandler.SignupUser(w, r)
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus != http.StatusOK {
				var errResponse map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &errResponse)
				require.NoError(t, err, "Actual body: %s", w.Body.String())
				assert.Contains(t, errResponse["error"], tt.expectedError)
				return
			}

			var response UserResponse
			err = json.Unmarshal(w.Body.Bytes(), &response)

			require.NoError(t, err)
			assert.Equal(t, tt.body.(dto.SignUpUser).Login, response.Login)
			assert.NotEmpty(t, response.ID)

			if tt.checkCookie {
				cookies := w.Result().Cookies()
				require.Len(t, cookies, 1)
				assert.Equal(t, authHandler.jwtConfig.CookieName, cookies[0].Name)
				assert.True(t, cookies[0].HttpOnly)
				assert.Equal(t, http.SameSiteStrictMode, cookies[0].SameSite)
				assert.NotEmpty(t, cookies[0].Value)
			}

			user, err := userStorage.ValidateUser(tt.body.(dto.SignUpUser).Login, tt.body.(dto.SignUpUser).Password)
			require.NoError(t, err)
			assert.Equal(t, tt.body.(dto.SignUpUser).Login, user.Username)
		})
	}
}

func TestSigninUser(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           interface{}
		setupStorage   func(*storage.UserSet)
		expectedStatus int
		expectedError  string
		checkCookie    bool
	}{
		{
			name:   "success signin",
			method: http.MethodPost,
			body: dto.SignUpUser{
				Login:    "test123",
				Password: "password123",
			},
			setupStorage: func(s *storage.UserSet) {
				_, err := s.CreateUser("test123", "password123")
				require.NoError(t, err)
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
			checkCookie:    true,
		},
		{
			name:   "user not found",
			method: http.MethodPost,
			body: dto.SignUpUser{
				Login:    "test123",
				Password: "password123",
			},
			setupStorage:   func(s *storage.UserSet) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  ErrBadCredentials.Error(),
			checkCookie:    false,
		},
		{
			name:   "wrong password",
			method: http.MethodPost,
			body: dto.SignInUser{
				Login:    "test123",
				Password: "password",
			},
			setupStorage: func(s *storage.UserSet) {
				_, err := s.CreateUser("testuser", "password123456")
				require.NoError(t, err)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  ErrBadCredentials.Error(),
			checkCookie:    false,
		},
		{
			name:   "without password",
			method: http.MethodPost,
			body: dto.SignUpUser{
				Login:    "test123",
				Password: "",
			},
			setupStorage:   func(s *storage.UserSet) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  ErrInvalidInput.Error(),
			checkCookie:    false,
		},
		{
			name:   "empty login",
			method: http.MethodPost,
			body: dto.SignUpUser{
				Login:    "",
				Password: "password123",
			},
			setupStorage:   func(s *storage.UserSet) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  ErrInvalidInput.Error(),
			checkCookie:    false,
		},
		{
			name:           "invalid json",
			method:         http.MethodPost,
			body:           "invalid json string",
			setupStorage:   func(s *storage.UserSet) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  ErrInvalidInput.Error(),
			checkCookie:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authHandler, userStorage := setupTest()
			tt.setupStorage(userStorage)

			jsonBody, err := json.Marshal(tt.body)
			require.NoError(t, err)
			r := httptest.NewRequest(tt.method, "/signin", bytes.NewBuffer(jsonBody))
			w := httptest.NewRecorder()

			authHandler.SigninUser(w, r)
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus != http.StatusOK {
				var errResponse map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &errResponse)
				require.NoError(t, err, "Actual body: %s", w.Body.String())
				assert.Contains(t, errResponse["error"], tt.expectedError)
				return
			}

			var response UserResponse
			err = json.Unmarshal(w.Body.Bytes(), &response)

			require.NoError(t, err)
			assert.Equal(t, tt.body.(dto.SignUpUser).Login, response.Login)
			assert.NotEmpty(t, response.ID)

			if tt.checkCookie {
				cookies := w.Result().Cookies()
				require.Len(t, cookies, 1)
				assert.Equal(t, authHandler.jwtConfig.CookieName, cookies[0].Name)
				assert.True(t, cookies[0].HttpOnly)
				assert.Equal(t, http.SameSiteStrictMode, cookies[0].SameSite)
				assert.NotEmpty(t, cookies[0].Value)
			}

			user, err := userStorage.ValidateUser(tt.body.(dto.SignUpUser).Login, tt.body.(dto.SignUpUser).Password)
			require.NoError(t, err)
			assert.Equal(t, tt.body.(dto.SignUpUser).Login, user.Username)
		})
	}
}

func TestLogout(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		setupStorage   func(*storage.UserSet)
		expectedStatus int
		checkCookie    bool
	}{
		{
			name:           "success logout",
			method:         http.MethodPost,
			setupStorage:   func(s *storage.UserSet) {},
			expectedStatus: http.StatusNoContent,
			checkCookie:    true,
		},
		{
			name:   "logout after being logged in",
			method: http.MethodPost,
			setupStorage: func(s *storage.UserSet) {
				_, err := s.CreateUser("test123", "password123")
				require.NoError(t, err)
			},
			expectedStatus: http.StatusNoContent,
			checkCookie:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authHandler, userStorage := setupTest()
			tt.setupStorage(userStorage)

			r := httptest.NewRequest(tt.method, "/logout", nil)
			w := httptest.NewRecorder()

			authHandler.LogOutUser(w, r)
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.checkCookie {
				cookies := w.Result().Cookies()
				require.Len(t, cookies, 1)
				assert.Equal(t, authHandler.jwtConfig.CookieName, cookies[0].Name)
				assert.True(t, cookies[0].HttpOnly)
				assert.Equal(t, http.SameSiteStrictMode, cookies[0].SameSite)
				assert.Equal(t, -1, cookies[0].MaxAge)
				assert.Equal(t, "", cookies[0].Value)
			}
		})
	}
}
