package models

import (
)

type AutorizarToken struct {
	User 		string		`json:"user"`
	Pass		string		`json:"pass"`
}

type Token struct {
  Token string `json:"token"`
}
