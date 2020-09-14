package poll

import (
	"encoding/json"
	"net/http"
	"polley/models"
	"polley/server/helpers"
	"polley/server/storage.go"

	"github.com/gorilla/mux"
)

type getPollResponse struct {
	*models.Poll
	VoteAllowed bool `json:"voteAllowed"`
}

func newGetPollResponse(storage *storage.Storage, r *http.Request, poll *models.Poll) getPollResponse {
	return getPollResponse{poll, helpers.IsVoteAllowed(storage, poll, r)}
}

func getPoll(storage *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uuid := mux.Vars(r)["uuid"]

		poll, err := storage.Polls.Read(uuid)
		if err != nil {
			helpers.WriteError(w, err, http.StatusBadRequest)
			return
		}

		response := newGetPollResponse(storage, r, poll)

		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			helpers.WriteError(w, err, http.StatusInternalServerError)
			return
		}
	}
}
