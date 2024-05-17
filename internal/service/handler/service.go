package handler

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	jwt2 "github.com/wardrobe-gateway/internal/jwt"
	"github.com/wardrobe-gateway/internal/model"
	"github.com/wardrobe-gateway/internal/service/user_service"
	"net/http"
	"time"
)

type service struct {
	userService user_service.Service
	jwtService  jwt2.Service
}

func NewHandler(userService user_service.Service, jwt jwt2.Service) *service {
	return &service{userService: userService,
		jwtService: jwt}
}

func (s *service) LoginUser(w http.ResponseWriter, r *http.Request) {
	var userResp model.User
	username := r.URL.Query().Get("username")
	if username == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	password := r.URL.Query().Get("password")
	if password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	url := fmt.Sprintf("http://localhost:8083/login?name=%s&password=%s", username, password)
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&userResp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	claims := &jwt.RegisteredClaims{
		Subject:   userResp.Name,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
	}

	token, err := s.jwtService.New(claims)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userResp.Password = ""
	loginResp := model.LoginResponse{
		Token: token,
		User:  userResp,
	}
	writeResponseJson(w, &loginResp)
}
