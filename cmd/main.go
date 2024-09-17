package main

import (
	"github.com/xavesen/search-api/internal/api"
	log "github.com/sirupsen/logrus"
)

func main() {
	server := api.NewServer("127.0.0.1:8897")

	log.Fatal(server.Start())
}