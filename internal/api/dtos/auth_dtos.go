package dtos

import "github.com/JGCaceres97/parking/internal/core/domain"

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Role      domain.Role `json:"role"`
	ExpiresIn int64       `json:"expires_in"`
	TokenType string      `json:"token_type"`
	Token     string      `json:"token"`
}
