package storage

import (
	"context"
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