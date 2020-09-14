package server

import (
	"encoding/json"
	"net/http"
	"polley/models"

	"github.com/gorilla/mux"
)

type pollResponse struct {
	*models.Poll
	VoteAllowed bool `json:"voteAllowed"`
}

func newPollResponse(s *Server, r *http.Request, poll *models.Poll) pollResponse {
	return pollResponse{poll, isVoteAllowed(poll.UUID, poll.Filter, r, s.ipsDB)}
}

func (s *Server) getPollHandler(w http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["uuid"]

	poll, err := s.pollDB.Read(uuid)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	response := newPollResponse(s, r, poll)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
}
