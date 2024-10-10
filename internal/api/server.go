package api

import (
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/xavesen/search-api/internal/queue"
	"github.com/xavesen/search-api/internal/storage"
)

type Server struct {
	listenAddr string
	router     	*mux.Router
	queue      	queue.Queue
	docStorage	storage.DocumentStorage
	userStorage storage.UserStorage
}

func NewServer(listenAddr string, queue queue.Queue, documentStorage storage.DocumentStorage, userStorage storage.UserStorage) *Server {
	log.Debug("Initializing server")

	server := Server{
		listenAddr: listenAddr,
		router:     mux.NewRouter(),
		queue:      queue,
		docStorage: documentStorage,
		userStorage: userStorage,
	}

	server.initialiseRoutes()
	return &server
}

func (s *Server) initialiseRoutes() {
	log.Debug("Initializing routes")

	s.router.HandleFunc("/ping", s.Ping).Methods("GET")
	s.router.HandleFunc("/indexDocuments", s.indexDocuments).Methods("POST")
	s.router.HandleFunc("/searchDocuments", s.searchDocuments).Methods("POST")
	s.router.HandleFunc("/createIndex", s.createIndex).Methods("POST")
}

func (s *Server) Start() error {
	log.Infof("Starting listening on %s", s.listenAddr)
	return http.ListenAndServe(s.listenAddr, s.router)
}
