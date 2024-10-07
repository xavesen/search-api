package storage

import (
	"context"

	"github.com/xavesen/search-api/internal/models"
)

type DocStorageMock struct {
	IndexError 		error
	SearchError		error
	Documents 		[]models.Document
	EsIndexExists 	bool
}

func (ds *DocStorageMock) SearchQuery(ctx context.Context, searchRequest *models.DocumentSearchRequest) ([]models.Document, error) {
	if ds.SearchError != nil {
		return []models.Document{}, ds.SearchError
	}

	return ds.Documents, nil
}

func (ds *DocStorageMock) IndexExists(ctx context.Context, indexName string) (bool, error) {
	if ds.IndexError != nil {
		return false, ds.IndexError
	}

	return ds.EsIndexExists, nil
}