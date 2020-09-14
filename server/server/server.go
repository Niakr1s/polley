package server

import (
	"net/http"
	"os"
	"polley/server/storage.go"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Server is used to handle poll requests.
type Server struct {
	storage *storage.Storage
	router  *mux.Router
}

// New constructs new server.
func New(storage *storage.Storage) *Server {
	server := &Server{
		storage: storage,
		router:  mux.NewRouter(),
	}
	server.configureHandlers()

	return server
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handlers.LoggingHandler(os.Stdout, s.router).ServeHTTP(w, r)
}

func (s *Server) configureHandlers() {
	s.router.HandleFunc("/createPoll", s.createPollHandler).Methods("POST")
	s.router.HandleFunc("/poll/{uuid}", s.getPollHandler).Methods("GET")
	s.router.HandleFunc("/poll/{uuid}", s.votePollHandler).Methods("PUT")
	s.router.HandleFunc("/getUUIDs", s.getUUIDsHandler).Methods("GET")
}
