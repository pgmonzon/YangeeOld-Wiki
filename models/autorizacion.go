package models

import (
	"time"
	
	"github.com/dgrijalva/jwt-go"
)

type AutorizarToken struct {
	Iss		string				`json:"iss"`
	Pas		string				`json:"pas"`
	iat		time.Time			`json:"iat"`
}

type Autorizar struct {
	Usuario		string				`json:"usuario"`
	Clave 		string				`json:"clave"`
}

type Token struct {
  Token string `json:"token"`
}

type TokenClaims struct {
  *jwt.StandardClaims
	Usuario		string				`json:"usuario"`
	Clave			string				`json:"clave"`
}
