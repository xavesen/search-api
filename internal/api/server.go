package api

import (
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/xavesen/search-api/internal/config"
	"github.com/xavesen/search-api/internal/middleware"
	"github.com/xavesen/search-api/internal/queue"
	"github.com/xavesen/search-api/internal/storage"
	"github.com/xavesen/search-api/internal/utils"
)

type Server struct {
	listenAddr string
	router     	*mux.Router
	queue      	queue.Queue
	docStorage	storage.DocumentStorage
	userStorage storage.UserStorage
	config		*config.Config
	tokenOp 	utils.TokenOperator
}

func NewServer(listenAddr string, queue queue.Queue, documentStorage storage.DocumentStorage, userStorage storage.UserStorage, config *config.Config, tokenOp utils.TokenOperator) *Server {
	log.Debug("Initializing server")

	server := Server{
		listenAddr: listenAddr,
		router:     mux.NewRouter(),
		queue:      queue,
		docStorage: documentStorage,
		userStorage: userStorage,
		config: config,
		tokenOp: tokenOp,
	}

	server.initialiseRoutes()
	return &server
}

func (s *Server) initialiseRoutes() {
	log.Debug("Initializing routes")

	s.router.HandleFunc("/ping", s.Ping).Methods("GET")
	s.router.HandleFunc("/login", s.login).Methods("POST")
	s.router.HandleFunc("/refresh", s.refresh).Methods("POST")

	privateRouter := s.router.PathPrefix("/").Subrouter()
	amw := middleware.AuthMiddleware{
		JwtKey: s.config.JwtKey,
		TokenOp: s.tokenOp,
		TokenHeaderName: s.config.TokenHeaderName,
	}
	privateRouter.Use(amw.Authenticate)

	privateRouter.HandleFunc("/indexDocuments", s.indexDocuments).Methods("POST")
	privateRouter.HandleFunc("/searchDocuments", s.searchDocuments).Methods("POST")
	privateRouter.HandleFunc("/createIndex", s.createIndex).Methods("POST")
}

func (s *Server) Start() error {
	log.Infof("Starting listening on %s", s.listenAddr)
	return http.ListenAndServe(s.listenAddr, s.router)
}
