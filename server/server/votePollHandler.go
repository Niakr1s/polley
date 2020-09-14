package server

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"polley/models"

	"github.com/gorilla/mux"
)

func (s *Server) votePollHandler(w http.ResponseWriter, r *http.Request) {
	type putPollHandlerRequest struct {
		ChoiceTexts []string `json:"choices"`
	}

	uuid := mux.Vars(r)["uuid"]

	request := &putPollHandlerRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	poll, err := s.pollController.Read(uuid)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}
	if poll.IsExpired() {
		writeError(w, errors.New("poll is expired"), http.StatusBadRequest)
		return
	}
	if !isVoteAllowed(uuid, poll.Filter, r, s.ipsController) {
		writeError(w, errors.New("vote isn't allowed"), http.StatusBadRequest)
		return
	}

	err = s.pollController.Increment(uuid, request.ChoiceTexts)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	err = s.storeVotedClient(poll, w, r)
	if err != nil {
		log.Printf("captureFilterInfo error: %v", err)
		return
	}
}

func (s *Server) storeVotedClient(poll *models.Poll, w http.ResponseWriter, r *http.Request) error {
	switch poll.Filter {
	case "ip":
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return err
		}
		return s.ipsController.AddIPForPoll(poll.UUID, ip)

	case "cookie":
		cookie := &http.Cookie{
			Name:    poll.UUID,
			Expires: poll.Expires,
		}
		http.SetCookie(w, cookie)
		log.Printf("added cookie %v", cookie)
		return nil
	}
	return nil
}
