package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Empresa struct {
	ID				bson.ObjectId	 	`bson:"_id" json:"id,omitempty"`
	Empresa  	string        	`json:"empresa"`
	Logo     	string					`json:"logo"`
	Modulos		[]IdModulo			`json:"modulos"`
	Activo		bool   					`json:"activo"`
  Borrado   bool          	`json:"borrado,omitempty"`
	Timestamp	time.Time				`json:"timestamp, omitempty"`
}

type IdModulo struct {
	ID				bson.ObjectId	 	`bson:"_id" json:"id,omitempty"`
}
