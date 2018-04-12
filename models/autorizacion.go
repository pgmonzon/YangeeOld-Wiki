package models

import (
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2/bson"
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
	Uid	bson.ObjectId `bson:"u_id" json:"uid,omitempty"`
	Usr string				`json:"usr"`
	Rbc	string				`json:"rbc"`
	Cid	bson.ObjectId `bson:"c_id" json:"cid,omitempty"`
	Clt	string				`json:"clt"`
	Eid	bson.ObjectId `bson:"e_id" json:"eid,omitempty"`
	*jwt.StandardClaims
}

type Autorizacion struct {
	Token			string		`json:"token"`
	Logo			string		`json:"logo"`
	Usuario		string		`json:"usuario"`
	Menu			[]Opcion	`json:"menu"`
}
