package http

import (
	"encoding/json"
	"net/http"
)

func errorResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": message}); err != nil {
		// If we can't even send the error response, there's not much we can do
		// but at least we tried to handle the error.
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
