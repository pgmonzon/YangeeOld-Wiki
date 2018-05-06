package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Personal struct {
	ID							      bson.ObjectId	 	`bson:"_id" json:"id,omitempty"`
	Empresa_id			      bson.ObjectId		`bson:"empresa_id" json:"empresa_id,omitempty"`
	Apellido 				      string        	`bson:"apellido" json:"apellido"`
  Nombre                string          `bson:"nombre" json:"nombre"`
  Categoria_id		      bson.ObjectId		`bson:"categoria_id" json:"categoria_id"`
	Propio   				      bool  					`json:"propio"`
  BasicoSindicato_id		bson.ObjectId		`bson:"basicoSindicato_id" json:"basicoSindicato_id,omitempty"`
  Comision              float32         `json:"comision"`
  Curso    				      time.Time				`json:"curso"`
	LNH      				      time.Time				`json:"lnh"`
  Registro     				  time.Time				`json:"registro"`
	Activo					      bool   					`json:"activo"`
	Borrado   			      bool          	`bson:"borrado" json:"borrado"`
	Timestamp			       	time.Time				`bson:"timestamp" json:"timestamp"`
}
