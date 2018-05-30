package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Autorizaciones struct {
	ID							          bson.ObjectId	 	`bson:"_id" json:"id,omitempty"`
	Empresa_id			          bson.ObjectId		`bson:"empresa_id" json:"empresa_id"`
	Viaje_id			            bson.ObjectId   `bson:"viaje_id" json:"viaje_id"`
	FechaHora                 time.Time       `bson:"fechaHora" json:"fechaHora"`
	Recorrido                 string          `bson:"recorrido" json:"recorrido"` // La primera y la última parada indicando si hay intermedias
	Kilometraje               int             `bson:"kilometraje" json:"kilometraje"`
	Solicitante_id            bson.ObjectId   `bson:"solicitante_id" json:"solicitante_id"`
  Solicitante		            string          `bson:"solicitante" json:"solicitante"`
  SolicitanteFecha          time.Time       `bson:"solicitanteFecha" json:"solicitanteFecha"`
	TipoSolicitud							string					`bson:"tipoSolicitud" json:"tipoSolicitud"` // Tarifa Cliente - Tarifa Transportista
	Titular_id				 	      bson.ObjectId		`bson:"titular_id" json:"titular_id"`
  Titular			         	    string          `bson:"titular" json:"titular"`
	ImporteTarifario					float64         `bson:"importeTarifario" json:"importeTarifario"`
	ImporteSugerido						float64         `bson:"importeSugerido" json:"importeSugerido"`
	Autorizante_id            bson.ObjectId   `bson:"autorizante_id" json:"solicitante_id"`
  Autorizante		            string          `bson:"autorizante" json:"autorizante"`
  AutorizanteFecha          time.Time       `bson:"autorizanteFecha" json:"autorizanteFecha"`
	ImporteAutorizado					float64         `bson:"importeAutorizado" json:"importeAutorizado"`
  Estado                    string          `bson:"estado" json:"estado"` // Pendiente - Autorizado - Rechazado
	Timestamp	     			      time.Time				`bson:"timestamp" json:"timestamp"`
}

type ImporteSugerido struct {
	Importe										float64         `bson:"importe" json:"importe"`
}
