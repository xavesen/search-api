package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/xavesen/search-api/internal/config"
	"github.com/xavesen/search-api/internal/storage"
	"github.com/xavesen/search-api/internal/utils"
)

type AuthMiddleware struct {
	TokenOp 			utils.TokenOperator
	UserStorage			storage.UserStorage
	Config				*config.Config
}

func (amw *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get(amw.Config.TokenHeaderName)

		valid, token, err := amw.TokenOp.ValidateToken(tokenStr, amw.Config.JwtKey)
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

		hashedToken := utils.Hash512WithSalt(tokenStr, amw.Config.JwtSalt)

		blacklisted, err := amw.UserStorage.CheckIfTokenBlacklisted(context.TODO(), hashedToken)
		if err != nil {
			utils.WriteJSON(w, r, http.StatusInternalServerError, false, "Internal server error", nil)
			return
		}

		if blacklisted {
			utils.WriteJSON(w, r, http.StatusUnauthorized, false, "Token is blacklisted", nil)
			return
		}

		userId, err := token.Claims.GetSubject()
		if err != nil {
			utils.WriteJSON(w, r, http.StatusUnauthorized, false, "Unauthorized", nil)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), utils.ContextKeyUserId, userId)))
	})
}