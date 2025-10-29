package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/JGCaceres97/parking/config"
	"github.com/JGCaceres97/parking/internal/api/dtos"
	"github.com/JGCaceres97/parking/internal/api/middlewares"
	"github.com/JGCaceres97/parking/internal/ports"
	"github.com/JGCaceres97/parking/pkg/response"
	"github.com/go-chi/chi/v5"
)

type ParkingHandler struct {
	service ports.ParkingService
}

func NewParkingHandler(service ports.ParkingService) *ParkingHandler {
	return &ParkingHandler{service: service}
}

func (h *ParkingHandler) RecordEntry(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), config.HandlerTimeout)
	defer cancel()

	userID, err := middlewares.GetUserIDFromContext(ctx)
	if err != nil {
		response.ErrorJSON(w, response.ErrInternalError, http.StatusInternalServerError)
		return
	}

	var req dtos.EntryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ErrorJSON(w, response.ErrInvalidJSON, http.StatusBadRequest)
		return
	}

	if req.VehicleTypeID == "" || req.LicensePlate == "" {
		response.ErrorJSON(w, response.ErrPlateAndTypeRequired, http.StatusBadRequest)
		return
	}

	record, err := h.service.RecordEntry(ctx, userID, req.VehicleTypeID, req.LicensePlate)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			response.ErrorJSON(w, response.ErrTimeout, http.StatusServiceUnavailable)
			return
		}

		if errors.Is(err, ports.ErrActiveParkingExists) {
			response.ErrorJSON(w, err, http.StatusConflict)
			return
		}

		if errors.Is(err, ports.ErrVehicleTypeNotFound) {
			response.ErrorJSON(w, err, http.StatusNotFound)
			return
		}

		response.ErrorJSON(w, response.ErrInternalError, http.StatusInternalServerError)
		return
	}

	response.JSON(w, http.StatusCreated, record)
}

func (h *ParkingHandler) RecordExit(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), config.HandlerTimeout)
	defer cancel()

	userID, err := middlewares.GetUserIDFromContext(ctx)
	if err != nil {
		response.ErrorJSON(w, response.ErrInternalError, http.StatusInternalServerError)
		return
	}

	var req dtos.ExitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ErrorJSON(w, response.ErrInvalidJSON, http.StatusBadRequest)
		return
	}

	if req.LicensePlate == "" {
		response.ErrorJSON(w, response.ErrPlateRequired, http.StatusBadRequest)
		return
	}

	exitRecord, err := h.service.RecordExit(ctx, userID, req.LicensePlate)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			response.ErrorJSON(w, response.ErrTimeout, http.StatusServiceUnavailable)
			return
		}

		if errors.Is(err, ports.ErrActiveParkingNotFound) {
			response.ErrorJSON(w, err, http.StatusNotFound)
			return
		}

		response.ErrorJSON(w, response.ErrInternalError, http.StatusInternalServerError)
		return
	}

	response.JSON(w, http.StatusCreated, exitRecord)
}

func (h *ParkingHandler) GetRecordByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), config.HandlerTimeout)
	defer cancel()

	recordID := chi.URLParam(r, "id")
	if recordID == "" {
		response.ErrorJSON(w, response.ErrRegistryIDRequired, http.StatusBadRequest)
		return
	}

	record, err := h.service.GetRecordByID(ctx, recordID)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			response.ErrorJSON(w, response.ErrTimeout, http.StatusServiceUnavailable)
			return
		}

		if errors.Is(err, ports.ErrParkingRecordNotFound) {
			response.ErrorJSON(w, err, http.StatusNotFound)
			return
		}

		response.ErrorJSON(w, response.ErrInternalError, http.StatusInternalServerError)
		return
	}

	response.JSON(w, http.StatusOK, record)
}

func (h *ParkingHandler) GetCurrentlyParked(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), config.HandlerTimeout)
	defer cancel()

	records, err := h.service.GetCurrentlyParked(ctx)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			response.ErrorJSON(w, response.ErrTimeout, http.StatusServiceUnavailable)
			return
		}

		response.ErrorJSON(w, response.ErrInternalError, http.StatusInternalServerError)
		return
	}

	response.JSON(w, http.StatusOK, records)
}

func (h *ParkingHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), config.HandlerTimeout)
	defer cancel()

	records, err := h.service.GetHistory(ctx)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			response.ErrorJSON(w, response.ErrTimeout, http.StatusServiceUnavailable)
			return
		}

		response.ErrorJSON(w, response.ErrInternalError, http.StatusInternalServerError)
		return
	}

	response.JSON(w, http.StatusOK, records)
}
