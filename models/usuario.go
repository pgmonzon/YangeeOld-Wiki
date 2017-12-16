package models

import (
	"gopkg.in/mgo.v2/bson"
)

type Usuario struct {
	ID        bson.ObjectId 	`bson:"_id" json:"id"`
	Usuario   string        	`json:"usuario"`
  Clave			int64 					`json:"clave"`
  Mail      string        	`json:"mail"`
}

type UsuarioRegistrar struct {
	Usuario			  string        	`json:"usuario"`
  Clave					string 					`json:"clave"`
  Mail      		string        	`json:"mail"`
}
