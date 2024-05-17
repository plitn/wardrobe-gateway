package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

type service struct {
	key []byte
}

func NewService(key string) *service {
	return &service{
		key: []byte(key),
	}
}

func (s *service) New(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.key)
}

func (s *service) Parse(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.key, nil
	})
}
