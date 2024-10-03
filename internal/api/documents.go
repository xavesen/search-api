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

	// TODO: check if index exists and user can write to it

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

	utils.WriteJSON(w, r, http.StatusOK, true, "Successfully sent documents for indexing", nil)
}

func (s *Server) searchDocuments(w http.ResponseWriter, r *http.Request) {
	var searchRequest *models.DocumentSearchRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&searchRequest); err != nil {
		utils.WriteJSON(w, r, http.StatusBadRequest, false, "Invalid request payload", nil)
		return
	}

	
}