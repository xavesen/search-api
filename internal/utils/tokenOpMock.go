package utils

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenOperatorMock struct {
	Token			string
	GenerateErr		error
	ValidateErr		error
	TokenValid 		bool
	ReturnedToken	*jwt.Token
}

func (tom *TokenOperatorMock) GenerateToken(userId string, currentTime time.Time, ttl int, key []byte) (string, error) {
	if tom.GenerateErr != nil {
		return "", tom.GenerateErr
	}
	return tom.Token + strconv.Itoa(ttl), nil
}

func (tom *TokenOperatorMock) ValidateToken(tokenStr string, key []byte) (bool, *jwt.Token, error) {
	if tom.ValidateErr != nil {
		return false, nil, tom.ValidateErr
	}
	return tom.TokenValid, tom.ReturnedToken, nil
}