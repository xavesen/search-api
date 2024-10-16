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

type UserStorage interface {
	CheckUserIndexRights(ctx context.Context, userId string, indexId string) (bool, error)
	AddIndexToUser(ctx context.Context, userId string, indexName string) error
	GetUserInfoByLogin (ctx context.Context, login string) (*models.User, error)
	SetRefreshToken(ctx context.Context, userId string, refreshToken string) error
	GetUserInfoById(ctx context.Context, userId string) (*models.User, error)
}