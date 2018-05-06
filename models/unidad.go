package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Unidad struct {
	ID							bson.ObjectId	 	`bson:"_id" json:"id,omitempty"`
	Empresa_id			bson.ObjectId		`bson:"empresa_id" json:"empresa_id,omitempty"`
	Unidad   				string        	`bson:"unidad" json:"unidad"`
	Propia   				bool  					`json:"propia"`
	VTV      				time.Time				`json:"vtv"`
  Ruta     				time.Time				`json:"ruta"`
  Poliza   				time.Time				`json:"poliza"`
  Seguro   				time.Time				`json:"seguro"`
	Activo					bool   					`json:"activo"`
	Borrado   			bool          	`bson:"borrado" json:"borrado"`
	Timestamp				time.Time				`bson:"timestamp" json:"timestamp"`
}
