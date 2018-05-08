package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Transportista struct {
	ID							bson.ObjectId	 	`bson:"_id" json:"id,omitempty"`
	Empresa_id			bson.ObjectId		`bson:"empresa_id" json:"empresa_id,omitempty"`
	Transportista		string        	`bson:"transportista" json:"transportista"`
  Mail            string          `json:"mail"`
	Contactos				[]ContactoTran	`json:"contactos,omitempty"`
	Tarifarios			[]TarifarioTran	`json:"tarifarios,omitempty"`
	Activo					bool   					`json:"activo"`
	Borrado   			bool          	`bson:"borrado" json:"borrado"`
	Timestamp				time.Time				`bson:"timestamp" json:"timestamp"`
}

type ContactoTran struct {
  Nombre          string          `json:"nombre"`
  Cargo           string          `json:"cargo"`
  Telefono        string          `json:"telefono"`
}

type TarifarioTran struct {
	Tarifario				string					`json:"tarifario"`
	Tipo						string					`json:"tipo"` // Recorrido - Kilometraje - Rango Kilometraje
	TipoUnidad_id		bson.ObjectId		`bson:"tipoUnidad_id" json:"tipoUnidad_id,omitempty"`
	Recorrido				[]ParadasTran		`json:"recorrido,omitempty"` // si es tipo recorrido
	KmDesde					int							`json:"kmDesde,omitempty"` // si es tipo rango kilometraje
	KmHasta					int							`json:"kmHasta,omitempty"` // si es tipo rango kilometraje
	Importe					float64					`json:"importe"`
	VigenteDesde		time.Time				`json:"vigenteDesde"`
	VigenteHasta		time.Time				`json:"vigenteHasta"`
	Activo					bool   					`json:"activo"`
}

type ParadasTran struct {
	Locacion_id			bson.ObjectId		`bson:"locacion_id" json:"locacion_id,omitempty"`
}
