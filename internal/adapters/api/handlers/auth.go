package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/JGCaceres97/parking/internal/adapters/api/dto"
	"github.com/JGCaceres97/parking/internal/application/auth"
	"github.com/JGCaceres97/parking/internal/domain"
	"github.com/JGCaceres97/parking/pkg/response"
)

type authHandler struct {
	service auth.Service
}

func NewAuthHandler(service auth.Service) *authHandler {
	return &authHandler{service: service}
}

func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ErrorJSON(w, response.ErrInvalidJSON, http.StatusBadRequest)
		return
	}

	out, err := h.service.Login(
		r.Context(),
		auth.LoginInput{Username: req.Username, Password: req.Password})

	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			response.ErrorJSON(w, response.ErrInvalidCredentials, http.StatusUnauthorized)
			return
		}

		if errors.Is(err, domain.ErrUserInactive) {
			response.ErrorJSON(w, response.ErrUserBlocked, http.StatusForbidden)
			return
		}

		response.ErrorJSON(w, response.ErrInternalError, http.StatusInternalServerError)
		return
	}

	response.JSON(
		w,
		http.StatusOK,
		dto.LoginResponse{
			Token:     out.Token,
			TokenType: out.TokenType,
			ExpiresIn: out.ExpiresIn,
			Role:      out.Role},
	)
}
