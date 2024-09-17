package api

import (
	"net/http"
	
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	listenAddr	string
	router		*mux.Router
}

func NewServer(listenAddr string) *Server {
	log.Debug("Initializing server")

	server := Server{
		listenAddr: listenAddr,
		router: 	mux.NewRouter(),
	}
	
	server.initialiseRoutes()
	return &server
}

func (s *Server) initialiseRoutes() {
	log.Debug("Initializing routes")

	s.router.HandleFunc("/ping", s.Ping).Methods("GET")
}

func (s *Server) Start() error {
	log.Infof("Starting listening on %s", s.listenAddr)
	return http.ListenAndServe(s.listenAddr, s.router)
}