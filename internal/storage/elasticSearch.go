package storage

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	log "github.com/sirupsen/logrus"
	"github.com/xavesen/search-api/internal/models"
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

func (es *ElasticSearchClient) SearchQuery(ctx context.Context, searchRequest *models.DocumentSearchRequest) ([]models.Document, error) {
	documents := []models.Document{}

	searchResult, err := es.Client.Search().
	Index(searchRequest.Index).
	Request(
		&search.Request{
			Query: &types.Query{
				QueryString: &types.QueryStringQuery{
					Query: searchRequest.Query,
				},
			},
		},
	).Do(ctx)
	if err != nil {
		log.Errorf("Error performing search request with query %s in index %s: %s", searchRequest.Query, searchRequest.Index, err)
		return []models.Document{}, err
	}

	for _, hit := range searchResult.Hits.Hits {
		var document models.Document
		err = json.Unmarshal(hit.Source_, &document)
		if err != nil {
			log.Errorf("Error unmarshalling hit from ES to document struct: %s", err)
			continue
		}
		documents = append(documents, document)
	}

	return documents, nil
}

func (es *ElasticSearchClient) IndexExists(ctx context.Context, indexName string) (bool, error) {
	exists, err := es.Client.Indices.Exists(indexName).Do(ctx)
	if err != nil {
		log.Errorf("Error performing index exists check in ES with index name %s: %s", indexName, err)
		return false, err
	}
	return exists, nil
}