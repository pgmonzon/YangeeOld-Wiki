package models

import (
	"github.com/dgrijalva/jwt-go"
)

type AutorizarToken struct {
	Usuario string		`json:"usuario"`
	Clave		string		`json:"clave"`
	Iat			int64	  	`json:"iat"`
	Aud			string		`json:"aud"`
}

type Token struct {
  Token string `json:"token"`
}

type TokenAutorizar struct {
	Usuario 	string		`json:"usuario"`
	Clave			string		`json:"clave"`
	Audience	string 		`json:"audience"`
}

type TokenAutorizarClaims struct {
	Usuario 	string		`json:"usuario"`
	Clave			string		`json:"clave"`
	jwt.StandardClaims
}
