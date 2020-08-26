package server

import (
	"encoding/json"
	"net/http"
	"polley/models"
)

type createPollHandlerRequest struct {
	Choices  []string `json:"choices" validate:"min=2"`
	Name     string   `json:"name" validate:"required,max=100"`
	Settings struct {
		AllowMultiple  int `json:"allowMultiple" validate:"min=1"`
		TimeoutMinutes int `json:"timeoutMinutes" validate:"min=1,max=120"`
	} `json:"settings"`
}

func (r createPollHandlerRequest) toPoll() (*models.Poll, error) {
	poll, err := models.NewPoll(r.Settings.TimeoutMinutes, r.Choices)
	if err != nil {
		return nil, err
	}
	poll = poll.WithAllowMultiple(r.Settings.AllowMultiple).WithName(r.Name)
	return poll, nil
}

type createPollHandlerResponse struct {
	UUID string `json:"uuid"`
}

func (s *Server) createPollHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	request := createPollHandlerRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}
	err = validate.Struct(request)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	poll, err := request.toPoll()
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	err = s.pollDB.Create(poll)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	response := createPollHandlerResponse{UUID: poll.UUID}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
}
