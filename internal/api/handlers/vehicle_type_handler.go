package handlers

import (
	"context"
	"net/http"

	"github.com/JGCaceres97/parking/config"
	"github.com/JGCaceres97/parking/internal/ports"
	"github.com/JGCaceres97/parking/pkg/response"
)

type VehicleTypeHandler struct {
	service ports.VehicleTypeService
}

func NewVehicleTypeHandler(service ports.VehicleTypeService) *VehicleTypeHandler {
	return &VehicleTypeHandler{service: service}
}

func (h *VehicleTypeHandler) ListAll(w http.ResponseWriter, r *http.Request) {
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
