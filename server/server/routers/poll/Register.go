package poll

import (
	"polley/server/storage.go"

	"github.com/gorilla/mux"
)

// RegisterPollRouter constructs new router, that handles polls.
func RegisterPollRouter(subRouter *mux.Router, storage *storage.Storage) {
	subRouter.HandleFunc("/createPoll", createPoll(storage)).Methods("POST")
	subRouter.HandleFunc("/poll/{uuid}", getPoll(storage)).Methods("GET")
	subRouter.HandleFunc("/poll/{uuid}", votePoll(storage)).Methods("PUT")
	subRouter.HandleFunc("/getUUIDs", getUUIDs(storage)).Methods("GET")
}
