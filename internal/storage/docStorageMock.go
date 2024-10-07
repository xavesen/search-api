package storage

import (
	"context"

	"github.com/xavesen/search-api/internal/models"
)

type DocStorageMock struct {
	Error 			error
	Documents 		[]models.Document
	EsIndexExists 	bool
}

func (ds *DocStorageMock) SearchQuery(ctx context.Context, searchRequest *models.DocumentSearchRequest) ([]models.Document, error) {
	if ds.Error != nil {
		return []models.Document{}, ds.Error
	}

	return ds.Documents, nil
}

func (ds *DocStorageMock) IndexExists(ctx context.Context, indexName string) (bool, error) {
	if ds.Error != nil {
		return false, ds.Error
	}

	return ds.EsIndexExists, nil
}