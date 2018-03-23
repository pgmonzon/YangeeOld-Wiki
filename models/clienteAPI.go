package models

import (
  "gopkg.in/mgo.v2/bson"
)

type ClienteAPI struct {
  ID          bson.ObjectId `bson:"_id" json:"id,omitempty"`
	ClienteAPI 	string		    `json:"clienteapi"`
	Firma 		  string		    `json:"firma"`
  Aes         string        `json:"aes"`
  Activo		  bool   				`json:"activo"`
  Borrado     bool         	`json:"borrado"`
}
