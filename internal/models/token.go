package models

type LoginRequest struct {
	Login 		string 	`json:"login"`
	Password 	string	`json:"password"`
}

type TokenResponse struct {
	AccessToken 	string	`json:"access_token"`
	RefreshToken 	string	`json:"refresh_token"`
}