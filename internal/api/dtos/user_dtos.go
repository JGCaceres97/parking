package dtos

import "github.com/JGCaceres97/parking/internal/core/domain"

type CreateUserRequest struct {
	Username string      `json:"username"`
	Password string      `json:"password"`
	Role     domain.Role `json:"role"`
	IsActive bool        `json:"is_active"`
}

type UpdateUserRequest struct {
	Username string      `json:"username"`
	Role     domain.Role `json:"role"`
	IsActive bool        `json:"is_active"`
}

type ToggleActiveRequest struct {
	IsActive bool `json:"is_active"`
}
