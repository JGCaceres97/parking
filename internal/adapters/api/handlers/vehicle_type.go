package handlers

import (
	"net/http"

	"github.com/JGCaceres97/parking/internal/application/vehicle_type"
	"github.com/JGCaceres97/parking/pkg/response"
)

type vehicleTypeHandler struct {
	service vehicle_type.Service
}

func NewVehicleTypeHandler(service vehicle_type.Service) *vehicleTypeHandler {
	return &vehicleTypeHandler{service: service}
}

func (h *vehicleTypeHandler) ListAll(w http.ResponseWriter, r *http.Request) {
	vts, err := h.service.ListAll(r.Context())
	if err != nil {
		response.ErrorJSON(w, response.ErrInternalError, http.StatusInternalServerError)
		return
	}

	response.JSON(w, http.StatusOK, vts)
}
