package server

import (
	"encoding/json"
	"net"
	"net/http"
	"polley/controllers"
)

func isVoteAllowed(uuid string, filter string, r *http.Request, ipsController controllers.IPsController) bool {
	switch filter {
	case "ip":
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return true
		}
		return ipsController.IsVoteAllowedForIP(uuid, ip)

	case "cookie":
		_, err := r.Cookie(uuid)
		if err != nil {
			return true
		}
		return false

	}
	return true
}

func writeError(w http.ResponseWriter, err error, code int) {
	type ResponsePayload struct {
		Error string `json:"error"`
	}
	w.WriteHeader(code)
	responsePayload := ResponsePayload{err.Error()}
	json.NewEncoder(w).Encode(responsePayload)
}
