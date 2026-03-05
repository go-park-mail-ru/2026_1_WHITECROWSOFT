package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
			return nil, errors.New("invalid algorithm")
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	uid, ok := claims["user_id"].(string)
	if !ok {
		return &TokenPayload{}, errors.New("user_id not found")
	}
	return &TokenPayload{UserID: uid}, nil
}
