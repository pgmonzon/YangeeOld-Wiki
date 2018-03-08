package models

import (
	"gopkg.in/mgo.v2/bson"
)

type Filosofo struct {
	ID				bson.ObjectId	 	`bson:"_id" json:"id,omitempty"`
	Filosofo	string        	`json:"filosofo"`
	Doctrina	string					`json:"doctrina"`
	Biografia	string					`json:"biografia"`
	Activo		bool   					`json:"activo"`
  Borrado   bool          	`json:"borrado"`
}

type Filosofos struct {
	Filosofos	[]Filosofo			`json:"filosofos"`
}
