package handler

import (
	"encoding/json"
	"net/http"
)

func jsonResponse(w http.ResponseWriter, body any, code int) error {
	if code != 0 {
		w.WriteHeader(code)
	}
	w.Header().Set("Content-Type", "application/json")
	if body == nil {
		return nil
	}
	return json.NewEncoder(w).Encode(body)
}
