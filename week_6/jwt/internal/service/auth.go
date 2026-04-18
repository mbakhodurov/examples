package service

import (
	"time"

	"github.com/mbakhodurov/examples/week_6/jwt/internal/model"
	"golang.org/x/crypto/bcrypt"
)

func (s *JWTService) Login(username, password string) (*model.TokenPair, error) {
	user, exists := s.users[username]
	if !exists {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return s.generateTokenPair(user)
}

// GetAccessToken - получение нового access токена по refresh токену
func (s *JWTService) GetAccessToken(refreshToken string) (string, time.Time, error) {
	claims, err := s.validateRefreshToken(refreshToken)
	if err != nil {
		return "", time.Time{}, err
	}

	user, exists := s.users[claims.Username]
	if !exists {
		return "", time.Time{}, ErrInvalidToken
	}

	return s.generateAccessToken(user)
}

// GetRefreshToken - получение нового refresh токена
func (s *JWTService) GetRefreshToken(refreshToken string) (string, time.Time, error) {
	claims, err := s.validateRefreshToken(refreshToken)
	if err != nil {
		return "", time.Time{}, err
	}

	user, exists := s.users[claims.Username]
	if !exists {
		return "", time.Time{}, ErrInvalidToken
	}

	return s.generateRefreshToken(user)
}
