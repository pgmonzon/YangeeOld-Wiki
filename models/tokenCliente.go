package models

import (
)

type ReqCliente struct {
	Usuario 	string		`json:"usuario"`
	Clave			string		`json:"clave"`
	Audience	string 		`json:"audience"`
}

type TokenCliente struct {
  Token 	string	`json:"token"`
}
