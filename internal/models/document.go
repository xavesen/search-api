package models

type Document struct {
	Title	string	`json:"title"`
	Text	string	`json:"text"`
}

type DocumentsForIndexing struct {
	Index		string		`json:"index_name"`
	UserId		string		`json:"user_id,omitempty"`
	Documents 	[]Document	`json:"documents"`
}

type DocumentSearchRequest struct {
	Index 		string	`json:"index_name"`
	Query		string	`json:"query"`
}