package models

import (
	"github.com/dgrijalva/jwt-go"
)

type AutorizarToken struct {
	User 		string		`json:"user"`
	Pass		string		`json:"pass"`
}

type Token struct {
  Token string `json:"token"`
}

type Autorizar struct {
	Usuario		string				`json:"usuario"`
	Clave 		string				`json:"clave"`
}

type TokenClaims struct {
  *jwt.StandardClaims
	Usuario		string				`json:"usuario"`
	Clave			string				`json:"clave"`
}
