package auth

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTManager interface {
	NewJWT(userId string) (string, time.Duration, error)
	ParseJWT(token string) (int, error)
}

type Manager struct {
	signingKey string
	ttl        time.Duration
}

func NewManager(signingKey string, ttl time.Duration) (*Manager, error) {
	if signingKey == "" {
		return nil, fmt.Errorf("signing key is empty")
	}

	if ttl == 0 {
		return nil, fmt.Errorf("ttl is zero")
	}

	return &Manager{signingKey: signingKey, ttl: ttl}, nil
}

func (m *Manager) NewJWT(userId string) (string, time.Duration, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(m.ttl).Unix(),
		Subject:   userId,
	})
	tokenString, err := token.SignedString([]byte(m.signingKey))

	return tokenString, m.ttl, err
}

func (m *Manager) ParseJWT(accessToken string) (int, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.signingKey), nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	userId, ok := claims["sub"].(string)
	if !ok {
		return 0, fmt.Errorf("invalid token")
	}

	return strconv.Atoi(userId)
}
