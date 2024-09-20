package api

import (
	"encoding/json"
	"net/http"

	"github.com/segmentio/kafka-go"
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
		log.Fatal("Error marshalling documents for index request to json after adding adding user_id to original struct from user") // TODO: structured logging with more info
		return
	}

	_, err = s.kafkaConn.WriteMessages(
		kafka.Message{Value: []byte(jsonIndexRequest)},
	)
	if err != nil {
		utils.WriteJSON(w, r, http.StatusInternalServerError, false, "Internal server error", nil)
		log.Fatalf("Error marshalling documents for index request to json after adding adding user_id to original struct from user") // TODO: handling errors without aborting the app
		return
	}

	utils.WriteJSON(w, r, http.StatusOK, true, "Successfully sent documents for indexing", nil)
}
