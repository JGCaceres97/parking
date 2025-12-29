package dto

import "github.com/JGCaceres97/parking/internal/domain"

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Role      domain.Role `json:"role"`
	ExpiresIn int64       `json:"expires_in"`
	TokenType string      `json:"token_type"`
	Token     string      `json:"token"`
}
