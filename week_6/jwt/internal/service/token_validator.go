package service

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/mbakhodurov/examples/week_6/jwt/internal/model"
)

func (s *JWTService) validateRefreshToken(tokenString string) (*model.Claims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(refreshTokenSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	// Проверяем тип токена
	if tokenType, ok := claims["type"].(string); !ok || tokenType != "refresh" {
		return nil, ErrInvalidToken
	}

	userID, ok := claims["user_id"].(float64) // JWT парсит числа как float64
	if !ok {
		return nil, ErrInvalidToken
	}

	username, ok := claims["username"].(string)
	if !ok {
		return nil, ErrInvalidToken
	}

	return &model.Claims{
		UserID:   int64(userID),
		Username: username,
	}, nil
}
