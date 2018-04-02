package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Filosofo struct {
	ID				bson.ObjectId	 	`bson:"_id" json:"id,omitempty"`
	Filosofo	string        	`json:"filosofo"`
	Doctrina	string					`json:"doctrina"`
	Biografia	string					`json:"biografia"`
	Activo		bool   					`json:"activo"`
  Borrado   bool          	`json:"borrado,omitempty"`
	Timestamp	time.Time				`json:"timestamp, omitempty"`
}
