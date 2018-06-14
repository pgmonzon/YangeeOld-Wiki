package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Liquidacion struct {
	ID							          bson.ObjectId	 	`bson:"_id" json:"id,omitempty"`
	Empresa_id			          bson.ObjectId		`bson:"empresa_id" json:"empresa_id,omitempty"`
  Liquidacion               int		          `bson:"liquidacion" json:"liquidacion"`
  Fecha                     time.Time       `bson:"fecha" json:"fecha"`
  Transportista_id          bson.ObjectId		`bson:"transportista_id" json:"transportista_id"`
  Transportista             string          `bson:"transportista" json:"transportista"`
  Descripcion               string          `bson:"descripcion" json:"descripcion"`
  Neto                      float64         `bson:"neto" json:"neto"`
  Iva105                    float64         `bson:"iva105" json:"iva105"`
  Iva21                     float64         `bson:"iva21" json:"iva21"`
  Total                     float64         `bson:"total" json:"total"`
  CuentaIngreso             string          `bson:"cuentaIngreso" json:"cuentaIngreso"`
  Viajes                    []ViajesLiq     `bson:"viajesLiq" json:"viajesLiq"`
  FechaLiquidacion          time.Time       `bson:"fechaLiquidacion" json:"fechaLiquidacion"`
  UsuarioLiquidacion_id     bson.ObjectId   `bson:"usuarioLiquidacion_id" json:"usuarioLiquidacion_id"`
  UsuarioLiquidacion        string          `bson:"usuarioLiquidacion" json:"usuarioLiquidacion"`
	Timestamp	     			      time.Time				`bson:"timestamp" json:"timestamp"`
}

type ViajesLiq struct {
	Viaje_id		   	bson.ObjectId		`bson:"viaje_id" json:"viaje_id,omitempty"`
  FechaHora       time.Time       `bson:"fechaHora" json:"fechaHora"`
  Recorrido       string          `bson:"recorrido" json:"recorrido"`
  Valor           float64	        `bson:"valor" json:"valor"`
}
