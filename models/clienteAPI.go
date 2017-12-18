package models

import (
  "gopkg.in/mgo.v2/bson"
)

type ClienteAPI struct {
  ID          bson.ObjectId 	`bson:"_id" json:"id"`
	ClienteAPI 	string		      `json:"clienteapi"`
	Firma 		  string		      `json:"firma"`
}

type ClienteAPIAlta struct {
	ClienteAPI 	string		      `json:"clienteapi"`
	Firma 		  string		      `json:"firma"`
}
