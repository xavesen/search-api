package storage

import (
	"context"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	log "github.com/sirupsen/logrus"
	"github.com/xavesen/search-api/internal/models"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
)

type ElasticSearchClient struct {
	Client 	*elasticsearch.TypedClient
}

func NewElasticSearchClient(addr []string, apiKey string) (*ElasticSearchClient, error) {
	// As it is an pet project addr will be http so no certs are considered
	cfg := elasticsearch.Config{
        Addresses: addr,
		APIKey: apiKey,
	}
	es, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		log.Errorf("Error creating elastic search client: %s", err)
		return nil, err
	}

	err = es.BaseClient.DiscoverNodes()
	if err != nil {
		log.Errorf("Error connecting to elastic search on %s: %s", strings.Join(addr, ", "), err)
		return nil, err
	}

	return &ElasticSearchClient{Client: es}, nil
}

func (es *ElasticSearchClient) SearchQuery(ctx context.Context, searchRequest *models.DocumentSearchRequest) {
	es.Client.Search().
	Index(searchRequest.Index).
	Request(
		&search.Request{
			Query: &types.Query{
				Match: map[string]types.MatchQuery{
					"title": {Query: searchRequest.Title},
					"text": {Query: searchRequest.Text},
				},
			},
		},
	).Do(ctx)
	
}