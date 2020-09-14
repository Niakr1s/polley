package poll

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"polley/models"
	"polley/server/helpers"
	"polley/server/storage.go"

	"github.com/gorilla/mux"
)

func votePoll(storage *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type putPollHandlerRequest struct {
			ChoiceTexts []string `json:"choices"`
		}

		uuid := mux.Vars(r)["uuid"]

		request := &putPollHandlerRequest{}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			helpers.WriteError(w, err, http.StatusBadRequest)
			return
		}

		poll, err := storage.Polls.Read(uuid)
		if err != nil {
			helpers.WriteError(w, err, http.StatusBadRequest)
			return
		}
		if poll.IsExpired() {
			helpers.WriteError(w, errors.New("poll is expired"), http.StatusBadRequest)
			return
		}
		if !helpers.IsVoteAllowed(uuid, poll.Filter, r, storage.Ips) {
			helpers.WriteError(w, errors.New("vote isn't allowed"), http.StatusBadRequest)
			return
		}

		err = storage.Polls.Increment(uuid, request.ChoiceTexts)
		if err != nil {
			helpers.WriteError(w, err, http.StatusBadRequest)
			return
		}

		err = storeVotedClient(storage, poll, w, r)
		if err != nil {
			log.Printf("captureFilterInfo error: %v", err)
			return
		}
	}
}

func storeVotedClient(storage *storage.Storage, poll *models.Poll, w http.ResponseWriter, r *http.Request) error {
	switch poll.Filter {
	case "ip":
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return err
		}
		return storage.Ips.AddIPForPoll(poll.UUID, ip)

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
