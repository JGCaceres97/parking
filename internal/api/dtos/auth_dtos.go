package dtos

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	ExpiresIn int64  `json:"expires_in"`
	TokenType string `json:"token_type"`
	Token     string `json:"token"`
}
