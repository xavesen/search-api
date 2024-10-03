package main

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/xavesen/search-api/internal/api"
	"github.com/xavesen/search-api/internal/queue"
	"github.com/xavesen/search-api/internal/storage"
)

func main() {
	kafkaQueue, err := queue.NewKafkaQueue(context.TODO(), []string{"localhost:9092"}, "test")
	if err != nil {
		os.Exit(1)
	}

	esClient, err := storage.NewElasticSearchClient([]string{"http://localhost:920/"}, "***") // TODO: put back api key and store all sensitive data in env file
	if err != nil {
		os.Exit(1)
	}

	server := api.NewServer("127.0.0.1:8897", kafkaQueue, esClient)

	log.Fatal(server.Start())
}