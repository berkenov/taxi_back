package handler

import (
	"encoding/json"
	"net/http"
)

// HealthResponse represents health check response
type HealthResponse struct {
	Status string `json:"status"`
}

// Health returns API health status
func Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(HealthResponse{Status: "ok"})
}
