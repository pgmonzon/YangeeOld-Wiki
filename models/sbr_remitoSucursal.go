package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type SbrRemitoSucursal struct {
	ID							bson.ObjectId	 	              `bson:"_id" json:"id,omitempty"`
	Empresa_id			bson.ObjectId		              `bson:"empresa_id" json:"empresa_id,omitempty"`
	Timestamp				time.Time				              `bson:"timestamp" json:"timestamp"`
  Fecha           time.Time                     `bson:"fecha" json:"fecha"`
  Envia_id        bson.ObjectId                 `bson:"envia_id" json:"envia_id"`
  Envia           string                        `bson:"envia" json:"envia"`
  DeSucursal_id   bson.ObjectId                 `bson:"deSucursal_id" json:"deSucursal_id"`
  DeSucursal      string                        `bson:"deSucursal" json:"deSucursal"`
  ASucursal_id    bson.ObjectId                 `bson:"aSucursal_id" json:"aSucursal_id"`
  ASucursal       string                        `bson:"aSucursal" json:"aSucursal"`
  Recibio_id      bson.ObjectId                 `bson:"recibio_id" json:"recibio_id"`
  Recibio         string                        `bson:"recibio" json:"recibio"`
  FechaRecepcion  time.Time                     `bson:"fechaRecepcion" json:"fechaRecepcion"`
  Estado          string                        `bson:"estado" json:"estado"` // Enviado, Recibido, Rechazado, Cancelado
  Detalle         []SbrRemitoSucursalDetalle    `bson:"detalle" json:"detalle"`
}

type SbrRemitoSucursalDetalle struct {
  SbrArticulo_id  bson.ObjectId   `bson:"articulo_id" json:"articulo_id,omitempty"`
  SbrArticulo     string          `bson:"articulo" json:"articulo"`
  Cantidad        int32           `bson:"cantidad" json:"cantidad"`
}
