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
		AllowMultiple  int    `json:"allowMultiple" validate:"min=1"`
		TimeoutMinutes int    `json:"timeoutMinutes" validate:"min=1,max=120"`
		Filter         string `json:"filter"`
	} `json:"settings"`
}

func (r createPollHandlerRequest) toPoll() *models.Poll {
	poll := models.NewPoll(models.PollArgs{
		TimeLimitMinutes: r.Settings.TimeoutMinutes,
		Choices:          r.Choices,
		Name:             r.Name,
		AllowMultiple:    r.Settings.AllowMultiple,
		Filter:           r.Settings.Filter,
	})
	return poll
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

	poll := request.toPoll()

	err = s.pollController.Create(poll)
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
