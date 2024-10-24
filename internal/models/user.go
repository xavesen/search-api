package models

type User struct {
	Id         		string 		`json:"id,omitempty" bson:"_id,omitempty"`
	Login      		string 		`json:"login"`
	Password   		string 		`json:"password"`
	IndexLimit 		int    		`json:"index_limit" bson:"indexlimit"`
	Indexes			[]string	`json:"indexes,omitempty"`
	RefreshToken	string		`json:"refresh_token" bson:"refreshToken"`
}