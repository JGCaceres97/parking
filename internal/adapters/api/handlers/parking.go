package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/JGCaceres97/parking/internal/adapters/api/dto"
	"github.com/JGCaceres97/parking/internal/adapters/api/middlewares"
	"github.com/JGCaceres97/parking/internal/application/parking"
	"github.com/JGCaceres97/parking/internal/domain"
	"github.com/JGCaceres97/parking/pkg/response"
)

type parkingHandler struct {
	service parking.Service
}

func NewParkingHandler(service parking.Service) *parkingHandler {
	return &parkingHandler{service: service}
}

func (h *parkingHandler) RecordEntry(w http.ResponseWriter, r *http.Request) {
	userID, err := middlewares.GetUserIDFromContext(r.Context())
	if err != nil {
		response.ErrorJSON(w, response.ErrInternalError, http.StatusInternalServerError)
		return
	}

	var req dto.EntryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ErrorJSON(w, response.ErrInvalidJSON, http.StatusBadRequest)
		return
	}

	if req.VehicleTypeID == "" || req.LicensePlate == "" {
		response.ErrorJSON(w, response.ErrPlateAndTypeRequired, http.StatusBadRequest)
		return
	}

	record, err := h.service.RecordEntry(r.Context(), userID, req.VehicleTypeID, req.LicensePlate)
	if err != nil {
		if errors.Is(err, domain.ErrActiveParkingAlreadyExists) {
			response.ErrorJSON(w, err, http.StatusConflict)
			return
		}

		if errors.Is(err, domain.ErrVehicleTypeNotFound) {
			response.ErrorJSON(w, err, http.StatusNotFound)
			return
		}

		response.ErrorJSON(w, response.ErrInternalError, http.StatusInternalServerError)
		return
	}

	response.JSON(w, http.StatusCreated, record)
}

func (h *parkingHandler) RecordExit(w http.ResponseWriter, r *http.Request) {
	userID, err := middlewares.GetUserIDFromContext(r.Context())
	if err != nil {
		response.ErrorJSON(w, response.ErrInternalError, http.StatusInternalServerError)
		return
	}

	var req dto.ExitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ErrorJSON(w, response.ErrInvalidJSON, http.StatusBadRequest)
		return
	}

	if req.LicensePlate == "" {
		response.ErrorJSON(w, response.ErrPlateRequired, http.StatusBadRequest)
		return
	}

	exitRecord, err := h.service.RecordExit(r.Context(), userID, req.LicensePlate)
	if err != nil {
		if errors.Is(err, domain.ErrActiveParkingNotFound) {
			response.ErrorJSON(w, err, http.StatusNotFound)
			return
		}

		response.ErrorJSON(w, response.ErrInternalError, http.StatusInternalServerError)
		return
	}

	response.JSON(w, http.StatusCreated, exitRecord)
}

func (h *parkingHandler) GetRecordByID(w http.ResponseWriter, r *http.Request) {
	recordID := chi.URLParam(r, "id")
	if recordID == "" {
		response.ErrorJSON(w, response.ErrRegistryIDRequired, http.StatusBadRequest)
		return
	}

	record, err := h.service.GetRecordByID(r.Context(), recordID)
	if err != nil {
		if errors.Is(err, domain.ErrParkingRecordNotFound) {
			response.ErrorJSON(w, err, http.StatusNotFound)
			return
		}

		response.ErrorJSON(w, response.ErrInternalError, http.StatusInternalServerError)
		return
	}

	response.JSON(w, http.StatusOK, record)
}

func (h *parkingHandler) GetCurrentlyParked(w http.ResponseWriter, r *http.Request) {
	records, err := h.service.GetCurrentlyParked(r.Context())
	if err != nil {
		response.ErrorJSON(w, response.ErrInternalError, http.StatusInternalServerError)
		return
	}

	response.JSON(w, http.StatusOK, records)
}

func (h *parkingHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	records, err := h.service.GetHistory(r.Context())
	if err != nil {
		response.ErrorJSON(w, response.ErrInternalError, http.StatusInternalServerError)
		return
	}

	response.JSON(w, http.StatusOK, records)
}
