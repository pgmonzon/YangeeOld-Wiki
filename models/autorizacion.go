package models

import (
	"github.com/dgrijalva/jwt-go"
)

type AutorizarTokenCliente struct {
	Usr string		`json:"usr"`
	Pas	string		`json:"pas"`
	*jwt.StandardClaims
}

type Token struct {
  Token 	string	`json:"token"`
}

type TokenAutorizado struct {
	Usr string		`json:"usr"`
	Rbc	string		`json:"rbc"`
	*jwt.StandardClaims
}
