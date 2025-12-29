package handlers

import (
	"context"
	"net/http"

	"github.com/JGCaceres97/parking/internal/application/vehicle_type"
	"github.com/JGCaceres97/parking/internal/infrastructure/config"
	"github.com/JGCaceres97/parking/pkg/response"
)

type vehicleTypeHandler struct {
	service vehicle_type.Service
}

func NewVehicleTypeHandler(service vehicle_type.Service) *vehicleTypeHandler {
	return &vehicleTypeHandler{service: service}
}

func (h *vehicleTypeHandler) ListAll(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), config.HandlerTimeout)
	defer cancel()

	vts, err := h.service.ListAll(ctx)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			response.ErrorJSON(w, response.ErrTimeout, http.StatusServiceUnavailable)
			return
		}

		response.ErrorJSON(w, response.ErrInternalError, http.StatusInternalServerError)
		return
	}

	response.JSON(w, http.StatusOK, vts)
}
