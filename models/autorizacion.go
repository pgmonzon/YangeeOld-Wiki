package models

import (
)

type AutorizarToken struct {
	User 		string		`json:"user"`
	Pass		string		`json:"pass"`
	Iat			int64	  	`json:"iat"`
	Aud			string		`json:"aud"`
}

type Token struct {
  Token string `json:"token"`
}
