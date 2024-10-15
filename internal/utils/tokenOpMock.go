package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenOperatorMock struct {

}

func (tom *TokenOperatorMock) GenerateToken(userId string, currentTime time.Time, ttl int, key []byte) (string, error) {
	return "aaa", nil
}

func (tom *TokenOperatorMock) ValidateToken(tokenStr string, key []byte) (bool, *jwt.Token, error) {
	return true, nil, nil
}