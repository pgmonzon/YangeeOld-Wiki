package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type TipoUnidad struct {
	ID							bson.ObjectId	 	`bson:"_id" json:"id,omitempty"`
	Empresa_id			bson.ObjectId		`bson:"empresa_id" json:"empresa_id,omitempty"`
	TipoUnidad			string        	`bson:"tipo_unidad" json:"tipo_unidad"`
	Activo					bool   					`json:"activo"`
	Borrado   			bool          	`bson:"borrado" json:"borrado"`
	Timestamp				time.Time				`bson:"timestamp" json:"timestamp"`
}
