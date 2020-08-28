package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

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
