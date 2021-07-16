package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

// TokenManager provides logic for JWT Access & Refresh tokens generation and parsing.
type TokenManager interface {
	NewAccessToken() (string, error)
	Parse(accessToken string) (string, error)
	NewRefreshToken() (string, error)
	ExtractToken(context.Context) (string, error)
}

type Manager struct {
	signingKey string
}

func NewManager(signingKey string) (*Manager, error) {
	if signingKey == "" {
		return nil, errors.New("empty signing key")
	}

	return &Manager{signingKey: signingKey}, nil
}

func (m *Manager) NewAccessToken(userId string) (string, error) {
	const ttl = time.Minute * time.Duration(15)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(ttl).Unix(),
		Subject:   userId,
	})

	return token.SignedString([]byte(m.signingKey))
}

func (m *Manager) NewRefreshToken(userId string) (string, error) {
	const ttl = time.Hour * time.Duration(24)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(ttl).Unix(),
		Subject:   userId,
	})

	return token.SignedString([]byte(m.signingKey))
}

func (m *Manager) Parse(accessToken string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.signingKey), nil
	})
	if err != nil {
		return nil, err
	} else if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("error getting user claims from token")
	}

	return claims, nil
}

func (m *Manager) GenerateTokenPair(userId string) (map[string]string, error) {
	accessToken, err := m.NewAccessToken(userId)
	if err != nil {
		return nil, err
	}
	refreshToken, err := m.NewRefreshToken(userId)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil
}

func (m *Manager) NewAccessTokenWithClaims(userId string, accessUuid string) (string, error) {
	const ttl = time.Minute * time.Duration(15)

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["access_uuid"] = accessUuid
	claims["sub"] = userId
	claims["exp"] = time.Now().Add(ttl).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(m.signingKey))
}

func (m *Manager) NewRefreshTokenWithClaims(userId string, refreshUuid string) (string, error) {
	const ttl = time.Hour * time.Duration(24)

	claims := jwt.MapClaims{}
	claims["refresh_uuid"] = refreshUuid
	claims["sub"] = userId
	claims["exp"] = time.Now().Add(ttl).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(m.signingKey))
}

func (m *Manager) GenerateTokenPairWithClaims(userId string, accessUuid string, refreshUuid string) (map[string]string, error) {
	accessToken, err := m.NewAccessTokenWithClaims(userId, accessUuid)
	if err != nil {
		return nil, err
	}
	refreshToken, err := m.NewRefreshTokenWithClaims(userId, refreshUuid)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil
}
