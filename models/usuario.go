package models

import (
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
	Menu				[]Opcion				`json:"menu"`
}

type Usuarios struct {
	Usuarios	[]Usuario2				`json:"usuarios"`
}
//**********SACAR
type Usuario2 struct {
	ID					bson.ObjectId	 	`bson:"_id" json:"id,omitempty"`
	Usuario			string        	`json:"usuario"`
	Clave				string					`json:"clave"`
	Mail				string					`json:"mail"`
}


type UsuarioX struct {
	ID        bson.ObjectId 	`bson:"_id" json:"id"`
	Usuario   string        	`json:"usuario"`
  Clave			int64 					`json:"clave"`
  Mail      string        	`json:"mail"`
}

type UsuarioRegistrar struct {
	Usuario			  string        	`json:"usuario"`
  Clave					string 					`json:"clave"`
  Mail      		string        	`json:"mail"`
}
