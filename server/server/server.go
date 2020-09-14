package server

import (
	"net/http"
	"os"
	"polley/controllers"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Server is used to handle poll requests.
type Server struct {
	pollController controllers.PollController
	ipsController  controllers.IPsController
	router         *mux.Router
}

// New constructs new server.
func New(pollController controllers.PollController, ipsController controllers.IPsController) *Server {
	server := &Server{
		pollController: pollController,
		ipsController:  ipsController,
		router:         mux.NewRouter(),
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
