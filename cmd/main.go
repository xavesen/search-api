package main

import (
	"context"

	"github.com/xavesen/search-api/internal/api"
	log "github.com/sirupsen/logrus"
	"github.com/segmentio/kafka-go"
)

func main() {
	kafkaConn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", "test", 0)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	server := api.NewServer("127.0.0.1:8897", kafkaConn)

	log.Fatal(server.Start())
}