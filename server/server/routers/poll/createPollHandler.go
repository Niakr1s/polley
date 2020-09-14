package poll

import (
	"encoding/json"
	"net/http"
	"polley/models"
	"polley/server/helpers"
	"polley/server/storage.go"
)

type createPollRequest struct {
	Choices  []string `json:"choices" validate:"min=2"`
	Name     string   `json:"name" validate:"required,max=100"`
	Settings struct {
		AllowMultiple  int    `json:"allowMultiple" validate:"min=1"`
		TimeoutMinutes int    `json:"timeoutMinutes" validate:"min=1,max=120"`
		Filter         string `json:"filter" validate:"oneof=none ip cookie"`
	} `json:"settings"`
}

func (r createPollRequest) toPoll() *models.Poll {
	return models.NewPoll(models.PollArgs{
		TimeLimitMinutes: r.Settings.TimeoutMinutes,
		Choices:          r.Choices,
		Name:             r.Name,
		AllowMultiple:    r.Settings.AllowMultiple,
		Filter:           r.Settings.Filter,
	})
}

type createPollResponse struct {
	UUID string `json:"uuid"`
}

func createPoll(storage *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		request := createPollRequest{}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			helpers.WriteError(w, err, http.StatusBadRequest)
			return
		}
		err = helpers.Validate.Struct(request)
		if err != nil {
			helpers.WriteError(w, err, http.StatusBadRequest)
			return
		}

		poll := request.toPoll()

		err = storage.Polls.Create(poll)
		if err != nil {
			helpers.WriteError(w, err, http.StatusInternalServerError)
			return
		}

		response := createPollResponse{UUID: poll.UUID}
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			helpers.WriteError(w, err, http.StatusInternalServerError)
			return
		}
	}
}
