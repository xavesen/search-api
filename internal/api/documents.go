package api

import (
	"context"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/xavesen/search-api/internal/models"
	"github.com/xavesen/search-api/internal/utils"
)

func (s *Server) indexDocuments(w http.ResponseWriter, r *http.Request) {
	var documentsIndexingRequest *models.DocumentsForIndexing

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&documentsIndexingRequest); err != nil {
		utils.WriteJSON(w, r, http.StatusBadRequest, false, "Invalid request payload", nil)
		return
	}
	// TODO: add payload validation

	indexExists, err := s.docStorage.IndexExists(context.TODO(), documentsIndexingRequest.Index)
	if err != nil {
		utils.WriteJSON(w, r, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}

	if !indexExists { // TODO: add user rights check
		utils.WriteJSON(w, r, http.StatusForbidden, false, "Index doesn't exist or you don't have access to it", nil)
		return
	}

	documentsIndexingRequest.UserId = "1" // TODO: hardcoded, change to real userId retrieved from token

	jsonIndexRequest, err := json.Marshal(documentsIndexingRequest)
	if err != nil {
		utils.WriteJSON(w, r, http.StatusInternalServerError, false, "Internal server error", nil)
		log.Error("Error marshalling documents for index request to json after adding adding user_id to original struct from user") // TODO: structured logging with more info
		return
	}

	err = s.queue.WriteMessage(context.TODO(), jsonIndexRequest)
	if err != nil {
		utils.WriteJSON(w, r, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}

	utils.WriteJSON(w, r, http.StatusOK, true, "", nil)
}

func (s *Server) searchDocuments(w http.ResponseWriter, r *http.Request) {
	var searchRequest *models.DocumentSearchRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&searchRequest); err != nil {
		utils.WriteJSON(w, r, http.StatusBadRequest, false, "Invalid request payload", nil)
		return
	}

	// TODO: validate payload
	
	indexExists, err := s.docStorage.IndexExists(context.TODO(), searchRequest.Index)
	if err != nil {
		utils.WriteJSON(w, r, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}

	if !indexExists { // TODO: add user rights check
		utils.WriteJSON(w, r, http.StatusForbidden, false, "Index doesn't exist or you don't have access to it", nil)
		return
	}

	documents, err := s.docStorage.SearchQuery(context.TODO(), searchRequest)
	if err != nil {
		utils.WriteJSON(w, r, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}

	utils.WriteJSON(w, r, http.StatusOK, true, "", documents)
}