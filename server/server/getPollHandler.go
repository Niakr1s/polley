package server

import (
	"encoding/json"
	"net"
	"net/http"
	"polley/models"

	"github.com/gorilla/mux"
)

type getPollHandlerResponse struct {
	*models.Poll
	VoteAllowed bool `json:"voteAllowed"`
}

func (s *Server) getPollHandler(w http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["uuid"]

	poll, err := s.pollDB.Read(uuid)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	response := getPollHandlerResponse{poll, s.isVoteAllowed(uuid, poll.Filter, r)}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
}

func (s *Server) isVoteAllowed(uuid string, filter string, r *http.Request) bool {
	switch filter {
	case "ip":
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return true
		}
		return s.ipsDB.IsVoteAllowedForIP(uuid, ip)

	case "cookie":
		_, err := r.Cookie(uuid)
		if err != nil {
			return true
		}
		return false

	}
	return true
}
