package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type SbrIngresoSucursal struct {
	ID							bson.ObjectId	 	              `bson:"_id" json:"id,omitempty"`
	Empresa_id			bson.ObjectId		              `bson:"empresa_id" json:"empresa_id,omitempty"`
	Timestamp				time.Time				              `bson:"timestamp" json:"timestamp"`
  Fecha           time.Time                     `bson:"fecha" json:"fecha"`
  Ingresante_id   bson.ObjectId                 `bson:"ingresante_id" json:"ingresante_id"`
  Ingresante      string                        `bson:"ingresante" json:"ingresante"`
  Sucursal_id     bson.ObjectId                 `bson:"sucursal_id" json:"sucursal_id"`
  Sucursal        string                        `bson:"sucursal" json:"sucursal"`
  Recibio_id      bson.ObjectId                 `bson:"recibio_id" json:"recibio_id"`
  Recibio         string                        `bson:"recibio" json:"recibio"`
  Detalle         []SbrIngresoSucursalDetalle   `bson:"detalle" json:"detalle"`
}

type SbrIngresoSucursalDetalle struct {
  SbrArticulo_id  bson.ObjectId   `bson:"articulo_id" json:"articulo_id,omitempty"`
  SbrArticulo     string          `bson:"articulo" json:"articulo"`
  Cantidad        int32           `bson:"cantidad" json:"cantidad"`
}
