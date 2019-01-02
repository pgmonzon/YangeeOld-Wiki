package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type SbrVentas struct {
	ID							bson.ObjectId	 	              `bson:"_id" json:"id,omitempty"`
	Empresa_id			bson.ObjectId		              `bson:"empresa_id" json:"empresa_id,omitempty"`
	Timestamp				time.Time				              `bson:"timestamp" json:"timestamp"`
  Sucursal_id     bson.ObjectId                 `bson:"sucursal_id" json:"sucursal_id"`
  Sucursal        string                        `bson:"sucursal" json:"sucursal"`
  Fecha           time.Time                     `bson:"fecha" json:"fecha"`
  Tarjeta         float64                       `bson:"tarjeta" json:"tarjeta"`
  Efectivo        float64                       `bson:"efectivo" json:"efectivo"`
  Total           float64                       `bson:"total" json:"total"`
	Gastos					float64												`bson:"gastos" json:"gastos"`
	ImporteRendir		float64												`bson:"importeRendir" json:"importeRendir"`
  Estado          string                        `bson:"estado" json:"estado"`
  RendidoA_id     bson.ObjectId                 `bson:"rendidoA_id" json:"rendidoA_id"`
  RendidoA        string                        `bson:"rendidoA" json:"rendidoA"`
}

type SbrVentasCrear struct {
  Sucursal_id     bson.ObjectId                 `bson:"sucursal_id" json:"sucursal_id"`
  Sucursal        string                        `bson:"sucursal" json:"sucursal"`
}

type SbrVentasCerrar struct {
	RendidoA_id     bson.ObjectId                 `bson:"rendidoA_id" json:"rendidoA_id"`
	RendidoA        string                        `bson:"rendidoA" json:"rendidoA"`
}

type SbrVentasDetalle struct {
	ID							bson.ObjectId	 	              `bson:"_id" json:"id,omitempty"`
	Empresa_id			bson.ObjectId		              `bson:"empresa_id" json:"empresa_id,omitempty"`
	Sucursal_id     bson.ObjectId                 `bson:"sucursal_id" json:"sucursal_id"`
	Sucursal        string                        `bson:"sucursal" json:"sucursal"`
  SbrVentas_id		bson.ObjectId		              `bson:"sbrVentas_id" json:"sbrVentas_id,omitempty"`
	Timestamp				time.Time				              `bson:"timestamp" json:"timestamp"`
  Vendedor_id     bson.ObjectId                 `bson:"vendedor_id" json:"vendedor_id"`
  Vendedor        string                        `bson:"vendedor" json:"vendedor"`
  SbrArticulo_id  bson.ObjectId                 `bson:"sbrArticulo_id" json:"sbrArticulo_id"`
  SbrArticulo     string                        `bson:"sbrArticulo" json:"sbrArticulo"`
  Importe         float64                       `bson:"importe" json:"importe"`
  Descuento       float64                       `bson:"descuento" json:"descuento"`
  Cobrado         float64                       `bson:"cobrado" json:"cobrado"`
  FormaPago       string                        `bson:"formaPago" json:"formaPago"`
}

type SbrVentasGastos struct {
	ID							bson.ObjectId	 	              `bson:"_id" json:"id,omitempty"`
	Empresa_id			bson.ObjectId		              `bson:"empresa_id" json:"empresa_id,omitempty"`
	Sucursal_id     bson.ObjectId                 `bson:"sucursal_id" json:"sucursal_id"`
	Sucursal        string                        `bson:"sucursal" json:"sucursal"`
  SbrVentas_id		bson.ObjectId		              `bson:"sbrVentas_id" json:"sbrVentas_id,omitempty"`
	Timestamp				time.Time				              `bson:"timestamp" json:"timestamp"`
  Vendedor_id     bson.ObjectId                 `bson:"vendedor_id" json:"vendedor_id"`
  Vendedor        string                        `bson:"vendedor" json:"vendedor"`
  CuentaGasto_id 	bson.ObjectId                 `bson:"cuentaGasto_id" json:"cuentaGasto_id"`
  CuentaGasto     string                        `bson:"cuentaGasto" json:"cuentaGasto"`
  Importe         float64                       `bson:"importe" json:"importe"`
}
