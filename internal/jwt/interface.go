package jwt

import "github.com/golang-jwt/jwt/v5"

type Service interface {
	New(claims jwt.Claims) (string, error)
	Parse(tokenString string) (*jwt.Token, error)
}
