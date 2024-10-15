package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/xavesen/search-api/internal/models"
	"github.com/xavesen/search-api/internal/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	var loginRequest *models.LoginRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&loginRequest); err != nil {
		utils.WriteJSON(w, r, http.StatusBadRequest, false, "Invalid request payload", nil)
		return
	}

	// TODO: validate payload

	user, err := s.userStorage.GetUserInfo(context.TODO(), loginRequest.Login)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.WriteJSON(w, r, http.StatusUnauthorized, false, "Unauthorized", nil)
		} else {
			utils.WriteJSON(w, r, http.StatusInternalServerError, false, "Internal server error", nil)
		}
		return
	}

	if user.Password != loginRequest.Password {
		utils.WriteJSON(w, r, http.StatusUnauthorized, false, "Unauthorized", nil)
		return
	}

	accessToken, err := s.tokenOp.GenerateToken(user.Id, time.Now(), s.config.JwtAccessTTL, s.config.JwtKey)
	if err != nil {
		utils.WriteJSON(w, r, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}

	refreshToken, err := s.tokenOp.GenerateToken(user.Id, time.Now(), s.config.JwtRefreshTTL, s.config.JwtKey)
	if err != nil {
		utils.WriteJSON(w, r, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}

	hashedRefreshToken := utils.Hash512WithSalt(refreshToken, s.config.JwtRefreshSalt)

	err = s.userStorage.SetRefreshToken(context.TODO(), user.Id, hashedRefreshToken)
	if err != nil {
		utils.WriteJSON(w, r, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}

	utils.WriteJSON(w, r, http.StatusOK, true, "", models.TokenResponse{AccessToken: accessToken, RefreshToken: refreshToken})
}