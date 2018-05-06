package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Permiso struct {
	ID        	bson.ObjectId 	`bson:"_id" json:"id,omitempty"`
	Permiso   	string        	`json:"permiso"`
	Modulo_id		bson.ObjectId	 	`bson:"modulo_id" json:"modulo_id"`
  Activo			bool   					`json:"activo"`
  Borrado   	bool          	`json:"borrado"`
	Timestamp		time.Time				`json:"timestamp, omitempty"`
}

type Rol struct {
	ID					bson.ObjectId	 	`bson:"_id" json:"id,omitempty"`
	Rol					string        	`json:"rol"`
	Empresa_id	bson.ObjectId		`bson:"empresa_id" json:"empresa_id"`
	Permisos		[]IdPermiso			`json:"permisos"`
  Activo			bool   					`json:"activo"`
  Borrado   	bool          	`json:"borrado"`
	Timestamp		time.Time				`json:"timestamp, omitempty"`
}

type IdRol struct {
	ID					bson.ObjectId	 	`bson:"_id" json:"id,omitempty"`
}

type IdPermiso struct {
	ID					bson.ObjectId	 	`bson:"_id" json:"id,omitempty"`
}
