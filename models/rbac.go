package models

import (
	"gopkg.in/mgo.v2/bson"
)

type Permiso struct {
	ID        bson.ObjectId 	`bson:"_id" json:"id"`
	Permiso   string        	`json:"permiso"`
  Activo		bool   					`json:"activo"`
  Borrado   bool          	`json:"borrado"`
}

type PermisoAlta struct {
	Permiso   string        	`json:"permiso"`
  Activo		bool   					`json:"activo"`
  Borrado   bool          	`json:"borrado"`
}
