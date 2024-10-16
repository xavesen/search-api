package storage

import (
	"context"
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/xavesen/search-api/internal/models"
)

type UserStorageMock struct {
	IndexRightsError		error
	AddIndexError 			error
	IndexAccess				bool
	User 					*models.User
	GetUserErr				error
	SetRefreshTokenErr		error
	TokenBlacklisted		bool
	TokenBlacklistedErr		error
	Testing 				*testing.T
	ExpectedToken 			string
}

func (us *UserStorageMock) CheckUserIndexRights(ctx context.Context, userId string, indexId string) (bool, error) {
	return us.IndexAccess, us.IndexRightsError
}

func (us *UserStorageMock) AddIndexToUser(ctx context.Context, userId string, indexName string) error {
	return us.AddIndexError
}

func (us *UserStorageMock) GetUserInfoByLogin(ctx context.Context, login string) (*models.User, error) {
	return us.User, us.GetUserErr
}

func (us *UserStorageMock) SetRefreshToken(ctx context.Context, userId string, refreshToken string) error {
	assert.Equal(us.Testing, refreshToken, us.ExpectedToken, "wrong token hash")
	return us.SetRefreshTokenErr
}

func (us *UserStorageMock) GetUserInfoById(ctx context.Context, userId string) (*models.User, error) {
	return us.User, us.GetUserErr
}

func (us *UserStorageMock) CheckIfTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	return us.TokenBlacklisted, us.TokenBlacklistedErr
}