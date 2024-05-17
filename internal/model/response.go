package model

type LoginResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}
