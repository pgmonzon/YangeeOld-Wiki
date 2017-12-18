package models

import (
	//"time"
)

type AutorizarToken struct {
	User 		string		`json:"user"`
	Pass		string		`json:"pass"`
	Iat			int64	  	`json:"iat"`
	Iss			string		`json:"iss"`
	Sub			string		`json:"sub"`
	Aud			string		`json:"aud"`

}

type Token struct {
  Token string `json:"token"`
}
