package handler

import (
	"encoding/json"
	"net/http"
	"os"
)

// DebugCredentials returns test phone/code when DEBUG_MODE=true (для pre-fill формы)
func DebugCredentials(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if os.Getenv("DEBUG_MODE") != "true" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"phone": "77000000000",
		"code":  "0000",
	})
}
