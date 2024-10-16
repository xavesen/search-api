package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

	user, err := s.userStorage.GetUserInfoByLogin(context.TODO(), loginRequest.Login)
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

	hashedRefreshToken := utils.Hash512WithSalt(refreshToken, s.config.JwtSalt)

	err = s.userStorage.SetRefreshToken(context.TODO(), user.Id, hashedRefreshToken)
	if err != nil {
		utils.WriteJSON(w, r, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}

	utils.WriteJSON(w, r, http.StatusOK, true, "", models.TokenResponse{AccessToken: accessToken, RefreshToken: refreshToken})
}

func (s *Server) refresh(w http.ResponseWriter, r *http.Request) {
	var refreshRequest *models.RefreshRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&refreshRequest); err != nil {
		utils.WriteJSON(w, r, http.StatusBadRequest, false, "Invalid request payload", nil)
		return
	}

	// TODO: validate payload

	valid, token, err := s.tokenOp.ValidateToken(refreshRequest.RefreshToken, s.config.JwtKey)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			utils.WriteJSON(w, r, http.StatusUnauthorized, false, "Refresh token has expired", nil)
		} else {
			utils.WriteJSON(w, r, http.StatusUnauthorized, false, "Unauthorized", nil)
		}
		return
	}
	if !valid {
		utils.WriteJSON(w, r, http.StatusUnauthorized, false, "Unauthorized", nil)
		return
	}

	userId, _ := token.Claims.GetSubject()
	hashedRefreshToken := utils.Hash512WithSalt(refreshRequest.RefreshToken, s.config.JwtSalt)

	blacklisted, err := s.userStorage.CheckIfTokenBlacklisted(context.TODO(), hashedRefreshToken)
	if err != nil {
		utils.WriteJSON(w, r, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}

	if blacklisted {
		utils.WriteJSON(w, r, http.StatusUnauthorized, false, "Token is blacklisted", nil)
		return
	}

	user, err := s.userStorage.GetUserInfoById(context.TODO(), userId)
	if err != nil {
		utils.WriteJSON(w, r, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}

	if user.RefreshToken != hashedRefreshToken {
		utils.WriteJSON(w, r, http.StatusUnauthorized, false, "Unauthorized", nil)
		return
	}

	accessToken, err := s.tokenOp.GenerateToken(userId, time.Now(), s.config.JwtAccessTTL, s.config.JwtKey)
	if err != nil {
		utils.WriteJSON(w, r, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}

	refreshToken, err := s.tokenOp.GenerateToken(userId, time.Now(), s.config.JwtRefreshTTL, s.config.JwtKey)
	if err != nil {
		utils.WriteJSON(w, r, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}

	hashedNewRefreshToken := utils.Hash512WithSalt(refreshToken, s.config.JwtSalt)

	err = s.userStorage.SetRefreshToken(context.TODO(), userId, hashedNewRefreshToken)
	if err != nil {
		utils.WriteJSON(w, r, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}

	utils.WriteJSON(w, r, http.StatusOK, true, "", models.TokenResponse{AccessToken: accessToken, RefreshToken: refreshToken})
}
