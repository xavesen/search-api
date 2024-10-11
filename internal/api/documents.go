package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	log "github.com/sirupsen/logrus"
	"github.com/xavesen/search-api/internal/models"
	"github.com/xavesen/search-api/internal/storage"
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

	documentsIndexingRequest.UserId = "66d8420df6e5311a791e0a08" // TODO: hardcoded, change to real userId retrieved from token
	userHasAccess, err := s.userStorage.CheckUserIndexRights(context.TODO(), documentsIndexingRequest.UserId, documentsIndexingRequest.Index)
	if err != nil {
		utils.WriteJSON(w, r, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}

	indexExists, err := s.docStorage.IndexExists(context.TODO(), documentsIndexingRequest.Index)
	if err != nil {
		utils.WriteJSON(w, r, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}

	if !indexExists || !userHasAccess {
		utils.WriteJSON(w, r, http.StatusForbidden, false, "Index doesn't exist or you don't have access to it", nil)
		return
	}

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

	userId := "66d8420df6e5311a791e0a08" // TODO: hardcoded, change to real userId retrieved from token
	userHasAccess, err := s.userStorage.CheckUserIndexRights(context.TODO(), userId, searchRequest.Index)
	if err != nil {
		utils.WriteJSON(w, r, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}

	indexExists, err := s.docStorage.IndexExists(context.TODO(), searchRequest.Index)
	if err != nil {
		utils.WriteJSON(w, r, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}

	if !indexExists || !userHasAccess {
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

func (s *Server) createIndex(w http.ResponseWriter, r *http.Request) {
	var createIndexRequest *models.CreateIndexRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&createIndexRequest); err != nil {
		utils.WriteJSON(w, r, http.StatusBadRequest, false, "Invalid request payload", nil)
		return
	}

	userId := "6707e924ce91f1780b739e8c" // TODO: hardcoded, change to real userId retrieved from token

	// TODO: validate payload

	err := s.docStorage.NewIndex(context.TODO(), createIndexRequest.Index)
	if err != nil {
		var esError *types.ElasticsearchError
		if errors.As(err, &esError) && esError.Status == 400 && esError.ErrorCause.Type == storage.ErrResourceAlreadyExists {
			utils.WriteJSON(w, r, http.StatusConflict, false, "Index with such name already exists", nil)
		} else {
			utils.WriteJSON(w, r, http.StatusInternalServerError, false, "Internal server error", nil)
		}
		return
	}

	err = s.userStorage.AddIndexToUser(context.TODO(), userId, createIndexRequest.Index)
	if err != nil {
		// TODO: delete existing index on error here
		utils.WriteJSON(w, r, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}

	utils.WriteJSON(w, r, http.StatusOK, true, "", nil)
}