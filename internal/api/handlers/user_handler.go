package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/JGCaceres97/parking/config"
	"github.com/JGCaceres97/parking/internal/api/dtos"
	"github.com/JGCaceres97/parking/internal/api/middlewares"
	"github.com/JGCaceres97/parking/internal/core/domain"
	"github.com/JGCaceres97/parking/internal/ports"
	"github.com/JGCaceres97/parking/pkg/response"
	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	service ports.UserService
}

func NewUserHandler(service ports.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), config.HandlerTimeout)
	defer cancel()

	userID, err := middlewares.GetUserIDFromContext(ctx)
	if err != nil {
		response.ErrorJSON(w, response.ErrInternalError, http.StatusInternalServerError)
		return
	}

	users, err := h.service.ListAll(ctx, userID)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			response.ErrorJSON(w, response.ErrTimeout, http.StatusServiceUnavailable)
			return
		}

		response.ErrorJSON(w, response.ErrInternalError, http.StatusInternalServerError)
		return
	}

	response.JSON(w, http.StatusOK, users)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), config.HandlerTimeout)
	defer cancel()

	var req dtos.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ErrorJSON(w, response.ErrInvalidJSON, http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" || req.Role == "" {
		response.ErrorJSON(w, response.ErrUserCreateValidation, http.StatusBadRequest)
		return
	}

	if req.Role != domain.RoleAdmin && req.Role != domain.RoleCommon {
		response.ErrorJSON(w, response.ErrInvalidRole, http.StatusBadRequest)
		return
	}

	newUser := &domain.User{
		Username: req.Username,
		Password: req.Password,
		Role:     req.Role,
		IsActive: req.IsActive,
	}

	user, err := h.service.Create(ctx, newUser)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			response.ErrorJSON(w, response.ErrTimeout, http.StatusServiceUnavailable)
			return
		}

		if errors.Is(err, ports.ErrUsernameExists) {
			response.ErrorJSON(w, err, http.StatusConflict)
			return
		}

		response.ErrorJSON(w, response.ErrInternalError, http.StatusInternalServerError)
		return
	}

	response.JSON(w, http.StatusCreated, user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), config.HandlerTimeout)
	defer cancel()

	authUserID, err := middlewares.GetUserIDFromContext(ctx)
	if err != nil {
		response.ErrorJSON(w, response.ErrInternalError, http.StatusInternalServerError)
		return
	}

	userID := chi.URLParam(r, "userID")
	if userID == "" {
		response.ErrorJSON(w, response.ErrInvalidID, http.StatusBadRequest)
		return
	}

	var req dtos.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ErrorJSON(w, response.ErrInvalidJSON, http.StatusBadRequest)
		return
	}

	if userID == authUserID && req.Role != "" {
		response.ErrorJSON(w, response.ErrChangeOwnRole, http.StatusForbidden)
		return
	}

	if req.Username == "" && req.Role == "" {
		response.ErrorJSON(w, response.ErrUpdateValidation, http.StatusBadRequest)
		return
	}

	if req.Role != "" && req.Role != domain.RoleAdmin && req.Role != domain.RoleCommon {
		response.ErrorJSON(w, response.ErrInvalidRole, http.StatusBadRequest)
		return
	}

	updatedUser := &domain.User{
		ID:       userID,
		Username: req.Username,
		Role:     req.Role,
		IsActive: req.IsActive,
	}

	user, err := h.service.Update(ctx, userID, updatedUser)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			response.ErrorJSON(w, response.ErrTimeout, http.StatusServiceUnavailable)
			return
		}

		if errors.Is(err, ports.ErrUserNotFound) {
			response.ErrorJSON(w, err, http.StatusNotFound)
			return
		}

		if errors.Is(err, ports.ErrUsernameExists) || errors.Is(err, ports.ErrAdminOperation) {
			response.ErrorJSON(w, err, http.StatusConflict)
			return
		}

		response.ErrorJSON(w, response.ErrInternalError, http.StatusInternalServerError)
		return
	}

	response.JSON(w, http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), config.HandlerTimeout)
	defer cancel()

	authUserID, err := middlewares.GetUserIDFromContext(ctx)
	if err != nil {
		response.ErrorJSON(w, response.ErrInternalError, http.StatusInternalServerError)
		return
	}

	userID := chi.URLParam(r, "userID")
	if userID == "" {
		response.ErrorJSON(w, response.ErrInvalidID, http.StatusBadRequest)
		return
	}

	if userID == authUserID {
		response.ErrorJSON(w, response.ErrOwnDelete, http.StatusForbidden)
		return
	}

	if err := h.service.Delete(ctx, userID); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			response.ErrorJSON(w, response.ErrTimeout, http.StatusServiceUnavailable)
			return
		}

		if errors.Is(err, ports.ErrUserNotFound) {
			response.ErrorJSON(w, err, http.StatusNotFound)
			return
		}

		if errors.Is(err, ports.ErrAdminOperation) {
			response.ErrorJSON(w, err, http.StatusConflict)
			return
		}

		response.ErrorJSON(w, response.ErrInternalError, http.StatusInternalServerError)
		return
	}

	response.JSON(w, http.StatusOK, nil)
}

func (h *UserHandler) UpdateUsername(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), config.HandlerTimeout)
	defer cancel()

	userID, err := middlewares.GetUserIDFromContext(ctx)
	if err != nil {
		response.ErrorJSON(w, response.ErrInternalError, http.StatusInternalServerError)
		return
	}

	var req dtos.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ErrorJSON(w, response.ErrInvalidJSON, http.StatusBadRequest)
		return
	}

	user, err := h.service.UpdateUsername(ctx, userID, req.Username)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			response.ErrorJSON(w, response.ErrTimeout, http.StatusServiceUnavailable)
			return
		}

		if errors.Is(err, ports.ErrUserNotFound) {
			response.ErrorJSON(w, err, http.StatusNotFound)
			return
		}

		if errors.Is(err, ports.ErrUsernameExists) || errors.Is(err, ports.ErrAdminOperation) {
			response.ErrorJSON(w, err, http.StatusConflict)
			return
		}

		response.ErrorJSON(w, response.ErrInternalError, http.StatusInternalServerError)
		return
	}

	response.JSON(w, http.StatusOK, user)
}

func (h *UserHandler) ToggleActiveStatus(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), config.HandlerTimeout)
	defer cancel()

	userID := chi.URLParam(r, "userID")
	if userID == "" {
		response.ErrorJSON(w, response.ErrInvalidID, http.StatusBadRequest)
		return
	}

	var req dtos.ToggleActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ErrorJSON(w, response.ErrInvalidJSON, http.StatusBadRequest)
		return
	}

	user, err := h.service.ToggleActive(ctx, userID, req.IsActive)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			response.ErrorJSON(w, response.ErrTimeout, http.StatusServiceUnavailable)
			return
		}

		if errors.Is(err, ports.ErrUserNotFound) {
			response.ErrorJSON(w, err, http.StatusNotFound)
			return
		}

		if errors.Is(err, ports.ErrAdminOperation) {
			response.ErrorJSON(w, err, http.StatusConflict)
			return
		}

		response.ErrorJSON(w, response.ErrInternalError, http.StatusInternalServerError)
		return
	}

	response.JSON(w, http.StatusOK, user)
}
