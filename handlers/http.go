package handlers

import (
	"encoding/json"
	"net/http"
)

// WriteResponse returns a json response with the given status
func WriteResponse(w http.ResponseWriter, status int, r interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(r)
}
