package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken  = errors.New("invalid token")
	ErrTokenCreation = errors.New("failed to create token")
	ErrNoUserID      = errors.New("user_id not found")
	ErrBadAlgorithm  = errors.New("invalid algorithm")
)

type TokenPayload struct {
	UserID string
}

func GenerateToken(userID string, duration time.Duration, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(duration).Unix(),
	})
	return token.SignedString([]byte(secretKey))
}

func ValidateToken(tokenStr string, secretKey string) (*TokenPayload, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrBadAlgorithm
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, ErrInvalidToken
	}

	uid, ok := claims["user_id"].(string)
	if !ok {
		return &TokenPayload{}, ErrNoUserID
	}
	return &TokenPayload{UserID: uid}, nil
}
