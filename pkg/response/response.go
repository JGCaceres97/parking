package response

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(data)
}

func ErrorJSON(w http.ResponseWriter, err error, status int) {
	JSON(w, status, map[string]any{
		"status": status,
		"error":  err.Error(),
	})
}
