package middleware

import (
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/xavesen/search-api/internal/utils"
)

type AuthMiddleware struct {
	TokenHeaderName		string
	TokenOp 			utils.TokenOperator
	JwtKey				[]byte
}

func (amw *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get(amw.TokenHeaderName)

		valid, _, err := amw.TokenOp.ValidateToken(tokenStr, amw.JwtKey)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				utils.WriteJSON(w, r, http.StatusUnauthorized, false, "Token has expired, refresh it or login again", nil)
			} else {
				utils.WriteJSON(w, r, http.StatusUnauthorized, false, "Unauthorized", nil)
			}
			return
		}

		if !valid {
			utils.WriteJSON(w, r, http.StatusUnauthorized, false, "Unauthorized", nil)
			return
		}

		// TODO: check if token is blacklisted

		next.ServeHTTP(w, r)
	})
}