package api

import (
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/xavesen/search-api/internal/queue"
)

type Server struct {
	listenAddr string
	router     *mux.Router
	queue      queue.Queue
}

func NewServer(listenAddr string, kafkaQueue *queue.KafkaQueue) *Server {
	log.Debug("Initializing server")

	server := Server{
		listenAddr: listenAddr,
		router:     mux.NewRouter(),
		queue:      kafkaQueue,
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
