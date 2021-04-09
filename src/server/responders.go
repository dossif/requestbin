package server

import (
	"encoding/json"
	"net/http"
)

// Return json-formatted HTTP error
func jsonError(w http.ResponseWriter, err string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	resp := struct {
		Code  int    `json:"code"`
		Error string `json:"error"`
	}{
		Code:  code,
		Error: err,
	}
	_ = json.NewEncoder(w).Encode(resp)
}
