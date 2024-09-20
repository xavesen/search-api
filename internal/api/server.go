package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	listenAddr string
	router     *mux.Router
	kafkaConn  *kafka.Conn
}

func NewServer(listenAddr string, kafkaConn *kafka.Conn) *Server {
	log.Debug("Initializing server")

	server := Server{
		listenAddr: listenAddr,
		router:     mux.NewRouter(),
		kafkaConn:  kafkaConn,
	}

	server.initialiseRoutes()
	return &server
}

func (s *Server) initialiseRoutes() {
	log.Debug("Initializing routes")

	s.router.HandleFunc("/ping", s.Ping).Methods("GET")
	s.router.HandleFunc("/indexDocuments", s.indexDocuments).Methods("POST")
}

func (s *Server) Start() error {
	log.Infof("Starting listening on %s", s.listenAddr)
	return http.ListenAndServe(s.listenAddr, s.router)
}
