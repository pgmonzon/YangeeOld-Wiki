package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Filosofo struct {
	ID							bson.ObjectId	 	`bson:"_id" json:"id,omitempty"`
	Empresa_id			bson.ObjectId		`bson:"empresa_id" json:"empresa_id,omitempty"`
	Filosofo				string        	`json:"filosofo"`
	Doctrina				string					`json:"doctrina"`
	Biografia				string					`json:"biografia"`
	Activo					bool   					`json:"activo"`
	Borrado   			bool          	`bson:"borrado" json:"borrado"`
	Timestamp				time.Time				`bson:"timestamp" json:"timestamp"`
}
