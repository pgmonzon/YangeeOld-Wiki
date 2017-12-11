package models

import (
	"github.com/dgrijalva/jwt-go"
)

type AutorizarToken struct {
	User 		string		`json:"user"`
	Pass		string		`json:"pass"`
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
