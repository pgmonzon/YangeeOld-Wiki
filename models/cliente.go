package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Cliente struct {
	ID							bson.ObjectId	 	`bson:"_id" json:"id,omitempty"`
	Empresa_id			bson.ObjectId		`bson:"empresa_id" json:"empresa_id,omitempty"`
	Cliente  				string        	`bson:"cliente" json:"cliente"`
	Cuit						string					`json:"cuit"`
	Contactos				[]ContactoClie	`json:"contactos,omitempty"`
	Tarifarios			[]TarifarioClie	`json:"tarifarios,omitempty"`
	Activo					bool   					`json:"activo"`
	Borrado   			bool          	`bson:"borrado" json:"borrado"`
	Timestamp				time.Time				`bson:"timestamp" json:"timestamp"`
}

type ContactoClie struct {
  Nombre          string          `json:"nombre"`
  Cargo           string          `json:"cargo"`
  Telefono        string          `json:"telefono"`
}

type TarifarioClie struct {
	Tarifario				string					`json:"tarifario"`
	Tipo						string					`json:"tipo"` // Recorrido - Kilometraje - Rango Kilometraje
	TipoUnidad_id		bson.ObjectId		`bson:"tipoUnidad_id" json:"tipoUnidad_id,omitempty"`
	Vuelta					string					`bson:"vuelta" json:"vuelta"` //  [vacío] - 2da - 3ra - 4ta
	TipoServicio		string					`bson:"tipoServicio" json:"tipoServicio"` // chofer - asistente
	Recorrido				[]ParadasClie		`json:"recorrido,omitempty"` // si es tipo recorrido
	KmDesde					int							`json:"kmDesde,omitempty"` // si es tipo rango kilometraje
	KmHasta					int							`json:"kmHasta,omitempty"` // si es tipo rango kilometraje
	Importe					float64					`json:"importe"`
	VigenteDesde		time.Time				`json:"vigenteDesde"`
	VigenteHasta		time.Time				`json:"vigenteHasta"`
	Activo					bool   					`json:"activo"`
}

type ParadasClie struct {
	Locacion_id			bson.ObjectId		`bson:"locacion_id" json:"locacion_id,omitempty"`
}
