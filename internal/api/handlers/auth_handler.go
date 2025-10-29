package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/JGCaceres97/parking/config"
	"github.com/JGCaceres97/parking/internal/api/dtos"
	"github.com/JGCaceres97/parking/internal/ports"
	"github.com/JGCaceres97/parking/pkg/response"
)

type AuthHandler struct {
	service ports.AuthService
}

func NewAuthHandler(service ports.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), config.HandlerTimeout)
	defer cancel()

	var req dtos.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ErrorJSON(w, response.ErrInvalidJSON, http.StatusBadRequest)
		return
	}

	authResponse, err := h.service.Login(ctx, req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			response.ErrorJSON(w, response.ErrTimeout, http.StatusServiceUnavailable)
			return
		}

		if errors.Is(err, ports.ErrInvalidCredentials) {
			response.ErrorJSON(w, response.ErrInvalidCredentials, http.StatusUnauthorized)
			return
		}

		if errors.Is(err, ports.ErrUserBlocked) {
			response.ErrorJSON(w, response.ErrUserBlocked, http.StatusForbidden)
			return
		}

		response.ErrorJSON(w, response.ErrInternalError, http.StatusInternalServerError)
		return
	}

	response.JSON(w, http.StatusOK, authResponse)
}
