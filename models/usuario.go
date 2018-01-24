package models

import (
	"gopkg.in/mgo.v2/bson"
)

type Usuario struct {
	ID				bson.ObjectId	 	`bson:"_id" json:"id,omitempty"`
	Usuario		string        	`json:"usuario"`
	Clave			string					`json:"clave"`
  //Encrip  	int64  					`json:"encrip,omitempty"`
	Mail			string					`json:"mail"`
	Activo		bool   					`json:"activo"`
  Borrado   bool          	`json:"borrado"`
	Roles			[]IdRol					`json:"roles"`
}

type Usuarios struct {
	Usuarios	[]Usuario				`json:"usuarios"`
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
