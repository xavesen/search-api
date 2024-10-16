package storage

import (
	"context"

	"github.com/xavesen/search-api/internal/models"
)

type UserStorageMock struct {
	IndexRightsError	error
	AddIndexError 		error
	IndexAccess			bool
}

func (us *UserStorageMock) CheckUserIndexRights(ctx context.Context, userId string, indexId string) (bool, error) {
	return us.IndexAccess, us.IndexRightsError
}

func (us *UserStorageMock) AddIndexToUser(ctx context.Context, userId string, indexName string) error {
	return us.AddIndexError
}

func (us *UserStorageMock) GetUserInfoByLogin(ctx context.Context, login string) (*models.User, error) {
	return nil, nil
}

func (us *UserStorageMock) SetRefreshToken(ctx context.Context, userId string, refreshToken string) error {
	return nil
}

func (us *UserStorageMock) GetUserInfoById(ctx context.Context, userId string) (*models.User, error) {
	return nil, nil
}