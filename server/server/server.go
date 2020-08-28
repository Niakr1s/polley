package server

import (
	"net/http"
	"os"
	"polley/db"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// Server is used to handle poll requests.
type Server struct {
	pollDB db.PollDB
	ipsDB  db.IPsDB
	router *mux.Router
}

// New constructs new server.
func New(pollDB db.PollDB, ipsDB db.IPsDB) *Server {
	server := &Server{
		pollDB: pollDB,
		ipsDB:  ipsDB,
		router: mux.NewRouter(),
	}
	server.configureHandlers()

	return server
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handlers.LoggingHandler(os.Stdout, s.router).ServeHTTP(w, r)
}

func (s *Server) configureHandlers() {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
		AllowedMethods:   []string{"POST", "GET", "PUT"},
	})
	s.router.Use(c.Handler)

	s.router.HandleFunc("/createPoll", s.createPollHandler).Methods("POST", "OPTIONS")
	s.router.HandleFunc("/poll/{uuid}", s.getPollHandler).Methods("GET", "OPTIONS")
	s.router.HandleFunc("/poll/{uuid}", s.putPollHandler).Methods("PUT", "OPTIONS")
	s.router.HandleFunc("/getNPolls", s.getNPollsHandler).Methods("GET", "OPTIONS")
}
