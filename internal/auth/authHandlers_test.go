package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/config"
	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/storage"
)

func setupTestHandler() *Handler {
	jwtConfig := config.JWTConfig{
		Secret: "test-secret-key-for-testing-only",
		Secure: false,
	}

	return NewHandler(jwtConfig, storage.NewUserSet())
}

func TestSignupUser(t *testing.T) {
	authHandler := setupTestHandler()

	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedLogin  string
		isError        bool
	}{
		{
			name:           "success",
			requestBody:    `{"login": "test123", "password": "Password123456"}`,
			expectedStatus: http.StatusOK,
			isError:        false,
			expectedLogin:  "test123",
		},
		{
			name:           "empty password",
			requestBody:    `{"login": "test123", "password": ""}`,
			expectedStatus: http.StatusBadRequest,
			isError:        true,
		},
		{
			name:           "empty login",
			requestBody:    `{"login": "", "password": "Password123456"}`,
			expectedStatus: http.StatusBadRequest,
			isError:        true,
		},
		{
			name:           "invalid json",
			requestBody:    `{"login": "admin", "password": `,
			expectedStatus: http.StatusBadRequest,
			isError:        true,
		},
		{
			name:           "two identical users",
			requestBody:    `{"login": "test123", "password": "Password123456"}`,
			expectedStatus: http.StatusConflict,
			isError:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBufferString(tt.requestBody))
			w := httptest.NewRecorder()

			authHandler.SignupUser(w, r)
			if w.Code != tt.expectedStatus {
				t.Errorf("Answer code: get %d, expected %d", w.Code, tt.expectedStatus)
			}

			if !tt.isError {
				var resp UserResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("Error in parsing JSON: %v", err)
				}

				if resp.Login != tt.expectedLogin {
					t.Errorf("Login: get %s, expected %s", resp.Login, tt.expectedLogin)
				}

				cookies := w.Result().Cookies()
				if len(cookies) == 0 || cookies[0].Name != authHandler.jwtConfig.CookieName {
					t.Error("JWT Cookie was not set")
				}
			} else {
				var errResp map[string]string
				json.NewDecoder(w.Body).Decode(&errResp)

				if _, ok := errResp["error"]; !ok {
					t.Error("Expected error field")
				}
			}
		})
	}
}
