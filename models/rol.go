package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Rol struct {
	ID							bson.ObjectId	 	`bson:"_id" json:"id,omitempty"`
	Empresa_id			bson.ObjectId		`bson:"empresa_id" json:"empresa_id,omitempty"`
	Rol      				string        	`bson:"rol" json:"rol"`
  Menu			      []Opcion       	`json:"menu, omitempty"`
	Activo					bool   					`json:"activo"`
	Borrado   			bool          	`bson:"borrado" json:"borrado"`
	Timestamp				time.Time				`bson:"timestamp" json:"timestamp"`
}

type RolEstado struct {
	ID							bson.ObjectId	 	`bson:"_id" json:"id,omitempty"`
	Empresa_id			bson.ObjectId		`bson:"empresa_id" json:"empresa_id,omitempty"`
	Rol      				string        	`bson:"rol" json:"rol"`
  Menu			      []OpcionEstado 	`json:"menu, omitempty"`
	Activo					bool   					`json:"activo"`
	Borrado   			bool          	`bson:"borrado" json:"borrado"`
	Timestamp				time.Time				`bson:"timestamp" json:"timestamp"`
}
