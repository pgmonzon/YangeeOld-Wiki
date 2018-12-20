package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Usuario struct {
	ID					bson.ObjectId	 	`bson:"_id" json:"id,omitempty"`
	Usuario			string        	`json:"usuario"`
	Clave				string					`json:"clave"`
	Mail				string					`json:"mail"`
	Apellido		string					`json:"apellido"`
	Nombre			string					`json:"nombre"`
	Empresa_id	bson.ObjectId	 	`bson:"empresa_id" json:"empresa_id"`
	Activo			bool   					`json:"activo"`
  Borrado   	bool          	`json:"borrado"`
	Roles				[]IdRol					`json:"roles"`
	Menu				[]Opcion				`json:"menu, omitempty"`
	Timestamp		time.Time				`json:"timestamp, omitempty"`
	Rol_id			bson.ObjectId	 	`bson:"rol_id" json:"rol_id"`
}

type UsuarioValidar struct {
	Usuario			string        	`json:"usuario"`
	Clave				string					`json:"clave"`
}

type UsuariosEmpresa struct {
	ID					bson.ObjectId	 	`bson:"_id" json:"id,omitempty"`
	Usuario			string        	`json:"usuario"`
	Apellido		string					`json:"apellido"`
	Nombre			string					`json:"nombre"`
}
