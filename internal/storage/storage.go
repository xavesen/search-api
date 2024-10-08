package storage

import (
	"context"
	"github.com/xavesen/search-api/internal/models"
)

type DocumentStorage interface {
	SearchQuery(ctx context.Context, searchRequest *models.DocumentSearchRequest) ([]models.Document, error)
	IndexExists(ctx context.Context, indexName string) (bool, error)
	NewIndex(ctx context.Context, indexName string) error
}