package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Locacion struct {
	ID							bson.ObjectId	 	`bson:"_id" json:"id,omitempty"`
	Empresa_id			bson.ObjectId		`bson:"empresa_id" json:"empresa_id,omitempty"`
	Locacion				string        	`bson:"locacion" json:"locacion"`
	Direccion				string					`json:"direccion"`
	Latitud  				string					`json:"latitud"`
  Longitud 				string					`json:"longitud"`
	Activo					bool   					`json:"activo"`
	Borrado   			bool          	`bson:"borrado" json:"borrado"`
	Timestamp				time.Time				`bson:"timestamp" json:"timestamp"`
}
