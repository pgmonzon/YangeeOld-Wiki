package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Haberes struct {
	ID							          bson.ObjectId	 	`bson:"_id" json:"id,omitempty"`
	Empresa_id			          bson.ObjectId		`bson:"empresa_id" json:"empresa_id,omitempty"`
  Año                       string          `bson:"año" json:"año"`
	Mes    										string 					`bson:"mes" json:"mes"`
  ComisionesDesde	          time.Time   		`bson:"comisionesDesde" json:"comisionesDesde"`
  ComisionesHasta	          time.Time   		`bson:"comisionesHasta" json:"comisionesHasta"`
  BasicosSindicato          []Basicos       `bson:"basicosSindicato" json:"basicosSindicato"`
  Editable                  bool            `bson:"editable" json:"editable"`
	Timestamp	     			      time.Time				`bson:"timestamp" json:"timestamp"`
}

type Basicos struct {
  Basico_id			    				bson.ObjectId	 	`bson:"basico_id" json:"basico_id,omitempty"`
	BasicoSindicato           string        	`bson:"basico_sindicato" json:"basico_sindicato"`
  Importe                   float64         `bson:"importe" json:"importe"`
}

type Novedades struct {
	ID							          bson.ObjectId	 	`bson:"_id" json:"id,omitempty"`
	Empresa_id			          bson.ObjectId		`bson:"empresa_id" json:"empresa_id,omitempty"`
	Haberes_id								bson.ObjectId		`bson:"haberes_id" json:"haberes_id"`
	Personal_id	   	          bson.ObjectId		`bson:"personal_id" json:"personal_id"`
	Personal									string					`bson:"personal" json:"personal"`
	BasicoSindicato						float64					`bson:"basicoSindicato" json:"basicoSindicato"`
	ViajesMes									float64					`bson:"viajesMes" json:"viajesMes"`
	Comision									float32					`bson:"comision" json:"comision"`
	ComisionEstimada					float64					`bson:"comisionEstimada" json:"comisionEstimada"`
	Diferencia								float64					`bson:"diferencia" json:"diferencia"`
	LiquidacionFinal					float64					`bson:"liquidacionFinal" json:"liquidacionFinal"`
	ComisionReal							float64					`bson:"comisionReal" json:"comisionReal"`
	AnticiposPendientes				float64					`bson:"anticiposPendientes" json:"anticiposPendientes"`
	AnticiposAplicados				float64					`bson:"anticiposAplicados" json: "anticiposAplicados"`
	NetoPagar									float64					`bson:"netoPagar" json:"netoPagar"`
	PagoBanco									float64					`bson:"pagoBanco" json:"pagoBanco"`
	PagoEfectivo							float64					`bson:"pagoEfectivo" json:"pagoEfectivo"`
	Pendiente									float64					`bson:"pendiente" json:"pendiente"`
}
