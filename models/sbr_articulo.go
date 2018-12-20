package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type SbrArticulo struct {
	ID							bson.ObjectId	 	`bson:"_id" json:"id,omitempty"`
	Empresa_id			bson.ObjectId		`bson:"empresa_id" json:"empresa_id,omitempty"`
	Activo					bool   					`json:"activo"`
	Borrado   			bool          	`bson:"borrado" json:"borrado"`
	Timestamp				time.Time				`bson:"timestamp" json:"timestamp"`
  Rubro_id      	bson.ObjectId	 	`bson:"rubro_id" json:"rubro_id"`
  Rubro           string          `bson:"rubro" json:"rubro"`
  SbrArticulo     string          `bson:"articulo" json:"articulo"`
  CodigoBarras    string          `bson:"codigoBarras" json:"codigoBarras"`
  Precio          float64         `bson:"precio" json:"precio"`
  EsPromo         bool            `bson:"esPromo" json:"esPromo"`
  Promos          []SbrPromos     `bson:"promos" json:"promos"`
}

type SbrPromos struct {
  SbrArticulo_id  bson.ObjectId   `bson:"articulo_id" json:"articulo_id,omitempty"`
  SbrArticulo     string          `bson:"articulo" json:"articulo"`
  Cantidad        int32           `bson:"cantidad" json:"cantidad"`
}

type SbrStock struct {
	ID							bson.ObjectId	 	`bson:"_id" json:"id,omitempty"`
	Empresa_id			bson.ObjectId		`bson:"empresa_id" json:"empresa_id,omitempty"`
  SbrSucursal_id  bson.ObjectId   `bson:"sucursal_id" json:"sucursal_id,omitempty"`
  SbrSucursal     string          `bson:"sucursal" json:"sucursal"`
	SbrArticulo_id  bson.ObjectId   `bson:"articulo_id" json:"articulo_id,omitempty"`
  SbrArticulo     string          `bson:"articulo" json:"articulo"`
  Cantidad        int32           `bson:"cantidad" json:"cantidad"`
}

type SbrStockSucursal struct {
	SbrArticulo_id  bson.ObjectId   `bson:"articulo_id" json:"articulo_id,omitempty"`
  SbrArticulo     string          `bson:"articulo" json:"articulo"`
  Cantidad        int32           `bson:"cantidad" json:"cantidad"`
}

type SbrArticuloStock struct {
	SbrArticulo_id  bson.ObjectId   `bson:"articulo_id" json:"articulo_id,omitempty"`
  SbrArticulo     string          `bson:"articulo" json:"articulo"`
  Total		        int32           `bson:"total" json:"total"`
	Stock						[]SbrDetStock		`bson:"stock" json:"stock"`
}

type SbrDetStock struct {
	SbrSucursal_id  bson.ObjectId   `bson:"sucursal_id" json:"sucursal_id,omitempty"`
  SbrSucursal     string          `bson:"sucursal" json:"sucursal"`
  Cantidad        int32           `bson:"cantidad" json:"cantidad"`
}

type SbrArticuloImportar struct {
	Rubro						string					`bson:"rubro" json:"rubro"`
	Articulo				string					`bson:"articulo" json:"articulo"`
	CodigoBarras		string					`bson:"codigoBarras" json:"codigoBarras"`
	Precio					float64					`bson:"precio" json:"precio"`
}
