package storage

import (
	"context"
)

type UserStorageMock struct {
	IndexRightsError	error
	IndexAccess			bool
}

func (us *UserStorageMock) CheckUserIndexRights(ctx context.Context, userId string, indexId string) (bool, error) {
	return us.IndexAccess, us.IndexRightsError
}