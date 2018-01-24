package models

import (
	"gopkg.in/mgo.v2/bson"
)

type Permiso struct {
	ID        bson.ObjectId 	`bson:"_id" json:"id,omitempty"`
	Permiso   string        	`json:"permiso"`
  Activo		bool   					`json:"activo"`
  Borrado   bool          	`json:"borrado"`
}

type IdPermiso struct {
	ID				bson.ObjectId	 	`bson:"_id" json:"id,omitempty"`
}

type Permisos struct {
	Permisos	[]Permiso				`json:"permisos"`
}

type Rol struct {
	ID				bson.ObjectId	 	`bson:"_id" json:"id,omitempty"`
	Rol				string        	`json:"rol"`
  Activo		bool   					`json:"activo"`
  Borrado   bool          	`json:"borrado"`
	Permisos	[]IdPermiso			`json:"permisos"`
}

type IdRol struct {
	ID				bson.ObjectId	 	`bson:"_id" json:"id,omitempty"`
}

type Roles struct {
	Roles			[]Rol						`json:"roles"`
}
