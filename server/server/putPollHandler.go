package server

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"

	"github.com/gorilla/mux"
)

type putPollHandlerRequest struct {
	ChoiceTexts []string `json:"choices"`
}

func (s *Server) putPollHandler(w http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["uuid"]

	request := &putPollHandlerRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	poll, err := s.pollDB.Read(uuid)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}
	if poll.IsExpired() {
		writeError(w, errors.New("poll is expired"), http.StatusBadRequest)
		return
	}
	if !s.isVoteAllowed(uuid, poll.Filter, r) {
		writeError(w, errors.New("vote isn't allowed"), http.StatusBadRequest)
		return
	}

	for _, choiceText := range request.ChoiceTexts {
		err = s.pollDB.Increment(uuid, choiceText)
		if err != nil {
			writeError(w, err, http.StatusBadRequest)
			return
		}
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	s.ipsDB.AddIPForPoll(uuid, ip)
}
