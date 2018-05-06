package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type CuentaGasto struct {
	ID							bson.ObjectId	 	`bson:"_id" json:"id,omitempty"`
	Empresa_id			bson.ObjectId		`bson:"empresa_id" json:"empresa_id,omitempty"`
	CuentaGasto			string        	`bson:"cuenta_gasto" json:"cuenta_gasto"`
	Activo					bool   					`json:"activo"`
	Borrado   			bool          	`bson:"borrado" json:"borrado"`
	Timestamp				time.Time				`bson:"timestamp" json:"timestamp"`
}
