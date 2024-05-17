package user_service

import (
	jwt2 "github.com/wardrobe-gateway/internal/jwt"
)

type service struct {
	jwtService jwt2.Service
}

func NewUserService(jwt jwt2.Service) *service {
	return &service{
		jwtService: jwt,
	}
}
