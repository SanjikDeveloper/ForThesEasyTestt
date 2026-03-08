package http

import (
	"encoding/json"
	"net/http"
)

func errorResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": message}); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
