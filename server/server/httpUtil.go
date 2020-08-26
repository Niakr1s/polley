package server

import (
	"encoding/json"
	"net/http"
)

func writeError(w http.ResponseWriter, err error, code int) {
	type ResponsePayload struct {
		Error string `json:"error"`
	}
	w.WriteHeader(code)
	responsePayload := ResponsePayload{err.Error()}
	json.NewEncoder(w).Encode(responsePayload)
}
