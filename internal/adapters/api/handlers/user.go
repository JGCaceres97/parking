package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/JGCaceres97/parking/internal/adapters/api/dto"
	"github.com/JGCaceres97/parking/internal/adapters/api/middlewares"
	"github.com/JGCaceres97/parking/internal/application/user"
	"github.com/JGCaceres97/parking/internal/domain"
	"github.com/JGCaceres97/parking/pkg/response"
)

type userHandler struct {
	service user.Service
}

func NewUserHandler(service user.Service) *userHandler {
	return &userHandler{service: service}
}

func (h *userHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	userID, err := middlewares.GetUserIDFromContext(r.Context())
	if err != nil {
		response.ErrorJSON(w, response.ErrInternalError, http.StatusInternalServerError)
		return
	}

	users, err := h.service.ListAll(r.Context(), userID)
	if err != nil {
		response.ErrorJSON(w, response.ErrInternalError, http.StatusInternalServerError)
		return
	}

	response.JSON(w, http.StatusOK, users)
}

func (h *userHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
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

	user, err := h.service.Create(r.Context(), newUser)
	if err != nil {
		if errors.Is(err, domain.ErrUsernameAlreadyExists) {
			response.ErrorJSON(w, err, http.StatusConflict)
			return
		}

		response.ErrorJSON(w, response.ErrInternalError, http.StatusInternalServerError)
		return
	}

	response.JSON(w, http.StatusCreated, user)
}

func (h *userHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	authUserID, err := middlewares.GetUserIDFromContext(r.Context())
	if err != nil {
		response.ErrorJSON(w, response.ErrInternalError, http.StatusInternalServerError)
		return
	}

	userID := chi.URLParam(r, "userID")
	if userID == "" {
		response.ErrorJSON(w, response.ErrInvalidID, http.StatusBadRequest)
		return
	}

	var req dto.UpdateUserRequest
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

	user, err := h.service.Update(r.Context(), userID, updatedUser)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			response.ErrorJSON(w, err, http.StatusNotFound)
			return
		}

		if errors.Is(err, domain.ErrUsernameAlreadyExists) || errors.Is(err, domain.ErrAdminProtected) {
			response.ErrorJSON(w, err, http.StatusConflict)
			return
		}

		response.ErrorJSON(w, response.ErrInternalError, http.StatusInternalServerError)
		return
	}

	response.JSON(w, http.StatusOK, user)
}

func (h *userHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	authUserID, err := middlewares.GetUserIDFromContext(r.Context())
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

	if err := h.service.Delete(r.Context(), userID); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			response.ErrorJSON(w, err, http.StatusNotFound)
			return
		}

		if errors.Is(err, domain.ErrAdminProtected) {
			response.ErrorJSON(w, err, http.StatusConflict)
			return
		}

		response.ErrorJSON(w, response.ErrInternalError, http.StatusInternalServerError)
		return
	}

	response.JSON(w, http.StatusOK, nil)
}

func (h *userHandler) UpdateUsername(w http.ResponseWriter, r *http.Request) {
	userID, err := middlewares.GetUserIDFromContext(r.Context())
	if err != nil {
		response.ErrorJSON(w, response.ErrInternalError, http.StatusInternalServerError)
		return
	}

	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ErrorJSON(w, response.ErrInvalidJSON, http.StatusBadRequest)
		return
	}

	user, err := h.service.UpdateUsername(r.Context(), userID, req.Username)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			response.ErrorJSON(w, err, http.StatusNotFound)
			return
		}

		if errors.Is(err, domain.ErrUsernameAlreadyExists) || errors.Is(err, domain.ErrAdminProtected) {
			response.ErrorJSON(w, err, http.StatusConflict)
			return
		}

		response.ErrorJSON(w, response.ErrInternalError, http.StatusInternalServerError)
		return
	}

	response.JSON(w, http.StatusOK, user)
}

func (h *userHandler) ToggleActiveStatus(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	if userID == "" {
		response.ErrorJSON(w, response.ErrInvalidID, http.StatusBadRequest)
		return
	}

	var req dto.ToggleActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ErrorJSON(w, response.ErrInvalidJSON, http.StatusBadRequest)
		return
	}

	user, err := h.service.ToggleActive(r.Context(), userID, req.IsActive)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			response.ErrorJSON(w, err, http.StatusNotFound)
			return
		}

		if errors.Is(err, domain.ErrAdminProtected) {
			response.ErrorJSON(w, err, http.StatusConflict)
			return
		}

		response.ErrorJSON(w, response.ErrInternalError, http.StatusInternalServerError)
		return
	}

	response.JSON(w, http.StatusOK, user)
}
