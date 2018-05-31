package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Factura struct {
	ID							          bson.ObjectId	 	`bson:"_id" json:"id,omitempty"`
	Empresa_id			          bson.ObjectId		`bson:"empresa_id" json:"empresa_id,omitempty"`
  Letra                     string          `bson:"letra" json:"letra"`
  Suc                       string          `bson:"suc" json:"suc"`
  Numero                    string          `bson:"numero" json:"numero"`
  Fecha                     time.Time       `bson:"fecha" json:"fecha"`
  Vencimiento               time.Time       `bson:"vencimiento" json:"vencimiento"`
  Cliente_id	   	          bson.ObjectId		`bson:"cliente_id" json:"cliente_id"`
  Cliente                   string          `bson:"cliente" json:"cliente"`
  Descripcion               string          `bson:"descripcion" json:"descripcion"`
  Neto                      float64         `bson:"neto" json:"neto"`
  Iva105                    float64         `bson:"iva105" json:"iva105"`
  Iva21                     float64         `bson:"iva21" json:"iva21"`
  Total                     float64         `bson:"total" json:"total"`
  CuentaIngreso             string          `bson:"cuentaIngreso" json:"cuentaIngresa"`
  Viajes                    []ViajesFact    `bson:"viajesFact" json:"viajesFact"`
  FechaFacturacion          time.Time       `bson:"fechaFacturacion" json:"fechaFacturacion"`
  UsuarioFacturacion_id     bson.ObjectId   `bson:"usuarioFacturacion_id" json:"usuarioFacturacion_id"`
  UsuarioFacturacion        string          `bson:"usuarioFacturacion" json:"usuarioFacturacion"`
	Timestamp	     			      time.Time				`bson:"timestamp" json:"timestamp"`
}

type ViajesFact struct {
	Viaje_id		   	bson.ObjectId		`bson:"viaje_id" json:"viaje_id,omitempty"`
  FechaHora       time.Time       `bson:"fechaHora" json:"fechaHora"`
  Recorrido       string          `bson:"recorrido" json:"recorrido"`
  Valor           float64	        `bson:"valor" json:"valor"`
}
