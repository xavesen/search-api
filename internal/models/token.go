package models

type LoginRequest struct {
	Login 		string 	`json:"login"`
	Password 	string	`json:"password"`
}

type RefreshRequest struct {
	RefreshToken	string	`json:"refresh_token"`
}

type TokenResponse struct {
	AccessToken 	string	`json:"access_token"`
	RefreshToken 	string	`json:"refresh_token"`
}