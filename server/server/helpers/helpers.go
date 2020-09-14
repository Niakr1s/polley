package helpers

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"polley/controllers"
	"strconv"
)

// IsVoteAllowed checks for vote allowed.
// filter should be one of 'ip' or 'cookie', otherwise it will always return true.
func IsVoteAllowed(uuid string, filter string, r *http.Request, ipsController controllers.IPsController) bool {
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

// WriteError writes error to response.
func WriteError(w http.ResponseWriter, err error, code int) {
	type ResponsePayload struct {
		Error string `json:"error"`
	}
	w.WriteHeader(code)
	responsePayload := ResponsePayload{err.Error()}
	json.NewEncoder(w).Encode(responsePayload)
}

// GetFirstIntValueFromQuery gets an integer value from query, if exists.
func GetFirstIntValueFromQuery(query url.Values, key string) (int, error) {
	values, ok := query[key]
	if !ok || len(values) == 0 {
		return 0, fmt.Errorf("no value exists for key=%s", key)
	}
	return strconv.Atoi(values[0])
}
