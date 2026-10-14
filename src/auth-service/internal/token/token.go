package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const ttl = 24 * time.Hour

var ErrInvalidToken = errors.New("invalid token")

type Manager interface {
	Issue(userID, email string) (string, error)
	Parse(tokenStr string) (userID, email string, err error)
}

type jwtManager struct {
	secret []byte
}

func New(secret string) Manager {
	return &jwtManager{secret: []byte(secret)}
}

func (m *jwtManager) Issue(userID, email string) (string, error) {
	claims := jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(ttl).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(m.secret)
}

func (m *jwtManager) Parse(tokenStr string) (string, string, error) {
	t, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return m.secret, nil
	})
	if err != nil || !t.Valid {
		return "", "", ErrInvalidToken
	}
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", ErrInvalidToken
	}
	sub, _ := claims["sub"].(string)
	email, _ := claims["email"].(string)
	if sub == "" {
		return "", "", ErrInvalidToken
	}
	return sub, email, nil
}


