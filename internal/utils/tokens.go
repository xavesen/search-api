package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
)

var ErrWrongSigningMethod = errors.New("expected HMAC signing method")

type TokenOperator interface {
	GenerateToken(userId string, currentTime time.Time, ttl int, key []byte) (string, error)
	ValidateToken(tokenStr string, key []byte) (bool, *jwt.Token, error)
}

type JwtTokenOperator struct {
}

func (jto *JwtTokenOperator) GenerateToken(userId string, currentTime time.Time, ttl int, key []byte) (string, error) {
	expirationTime := currentTime.Add(time.Duration(ttl) * time.Second)
	claims := jwt.RegisteredClaims{
		Subject: userId,
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(key)
	if err != nil {
		log.Errorf("Error signing jwt token: %s", err)
		return "", err
	}
	return tokenString, nil
}

func (jto *JwtTokenOperator) ValidateToken(tokenStr string, key []byte) (bool, *jwt.Token, error) {
	if tokenStr == "" {
		log.Warning("Token validation error: no token passed")
		return false, nil, jwt.ErrTokenMalformed
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Errorf("Token validation error: token is signed with wrong signing method, expected HMAC got %s", token.Header["alg"])
			return nil, ErrWrongSigningMethod
		}

		if _, err := token.Claims.GetSubject(); err != nil {
			log.Errorf("Token validation error: token misses sub claim")
			return nil, jwt.ErrTokenRequiredClaimMissing
		}
	
		return key, nil
	}, jwt.WithExpirationRequired())

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			log.Debugf("Token validation error: %s", err)
		} else {
			log.Errorf("Token validation error: %s", err)
		}
		return false, nil, err
	}

	if !token.Valid {
		return false, nil, nil
	}

	return true, token, nil
}