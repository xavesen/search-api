package main

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/xavesen/search-api/internal/api"
	"github.com/xavesen/search-api/internal/queue"
)

func main() {
	kafkaQueue, err := queue.NewKafkaQueue(context.TODO(), []string{"localhost:9092"}, "test")
	if err != nil {
		os.Exit(1)
	}

	server := api.NewServer("127.0.0.1:8897", kafkaQueue)

	log.Fatal(server.Start())
}