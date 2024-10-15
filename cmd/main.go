package main

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/xavesen/search-api/internal/api"
	"github.com/xavesen/search-api/internal/config"
	"github.com/xavesen/search-api/internal/queue"
	"github.com/xavesen/search-api/internal/storage"
	"github.com/xavesen/search-api/internal/utils"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		os.Exit(1)
	}

	ctx := context.TODO()
	mongoStorage, err := storage.NewMongoStorage(ctx, config.DbAddr, config.Db, config.DbUser, config.DbPass)
	if err != nil {
		os.Exit(1)
	}

	kafkaQueue, err := queue.NewKafkaQueue(context.TODO(), config.KafkaAddrs, config.KafkaTopic)
	if err != nil {
		os.Exit(1)
	}

	esClient, err := storage.NewElasticSearchClient(config.ElasticSearchURLs, config.ElasticSearchKey)
	if err != nil {
		os.Exit(1)
	}

	tokenOp := &utils.JwtTokenOperator{}

	server := api.NewServer(config.ListenAddr, kafkaQueue, esClient, mongoStorage, config, tokenOp)

	log.Fatal(server.Start())
}