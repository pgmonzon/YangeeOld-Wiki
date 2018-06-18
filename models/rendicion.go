package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Rendicion struct {
	ID							          bson.ObjectId	 	`bson:"_id" json:"id,omitempty"`
	Empresa_id			          bson.ObjectId		`bson:"empresa_id" json:"empresa_id,omitempty"`
  Personal_id	   	          bson.ObjectId		`bson:"personal_id" json:"personal_id"`
  Personal                  string          `bson:"personal" json:"personal"`
  FechaHora                 time.Time       `bson:"fechaHora" json:"fechaHora"`
  CuentaGasto_id 	          bson.ObjectId		`bson:"cuentaGasto_id" json:"cuentaGasto_id"`
  CuentaGasto               string          `bson:"cuentaGasto" json:"cuentaGasto"`
  Observaciones             string          `bson:"observaciones" json:"observaciones"`
  Ingreso                   float64         `bson:"ingreso" json:"ingreso"`
  Egreso                    float64         `bson:"egreso" json:"egreso"`
  Saldo                     float64         `bson:"saldo" json:"saldo"`
  Timestamp	     			      time.Time				`bson:"timestamp" json:"timestamp"`
}
