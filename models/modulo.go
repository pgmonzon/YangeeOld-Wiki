package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Modulo struct {
	ID				bson.ObjectId	 	`bson:"_id" json:"id,omitempty"`
	Modulo   	string        	`json:"modulo"`
	Activo		bool   					`json:"activo"`
  Borrado   bool          	`json:"borrado,omitempty"`
	Timestamp	time.Time				`json:"timestamp, omitempty"`
}

type Modulos struct {
	Modulos		[]Modulo				`json:"modulos"`
}

type IdModulo struct {
	ID				bson.ObjectId	 	`bson:"_id" json:"id,omitempty"`
}
