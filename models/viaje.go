package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Viaje struct {
	ID							          bson.ObjectId	 	`bson:"_id" json:"id,omitempty"`
	Empresa_id			          bson.ObjectId		`bson:"empresa_id" json:"empresa_id,omitempty"`
  FechaHora                 time.Time       `json:"fechaHora"`
  Cliente_id	   	          bson.ObjectId		`bson:"cliente_id" json:"cliente_id,omitempty"`
  Cliente                   string          `json:"cliente"`
  TipoUnidad_id	 	          bson.ObjectId		`bson:"tipoUnidad_id" json:"tipoUnidad_id,omitempty"`
  TipoUnidad                string          `json:"tipoUnidad"`
  Transportista_id	 	      bson.ObjectId		`bson:"transportista_id" json:"transportista_id,omitempty"`
  Transportista             string          `json:"transportista"`
  Unidad_id	        	      bson.ObjectId		`bson:"unidad_id" json:"unidad_id,omitempty"`
  Unidad                    string          `json:"unidad"`
  Personal_id	   	          bson.ObjectId		`bson:"personal_id" json:"personal_id,omitempty"`
  Personal                  string          `json:"personal"`
  Paradas	    		  	      []ParadasServ		`json:"paradas,omitempty"`
  Recorrido                 string          `json:"recorrido,omitempty"` // La primera y la última parada indicando si hay intermedias
  Kilometraje               int             `json:"kilometraje"`
  TarifarioCliente          string          `json:"tarifarioCliente,omitempty"`
  TarifaValor               float64         `json:"tarifaValor,omitempty"`
  ValorViaje                float64         `json:"valorViaje,omitempty"`
  AutValorViaje_id          bson.ObjectId   `bson:"autValorViaje_id" json:"autValorViaje_id,omitempty"` // el usuario que autorizó la tarifa
  AutValor                  string          `json:"autValor,omitempty"` // usuario que autorizó
  AutValorViajeFecha        time.Time       `json:"autValorViajeFecha,omitempty"`
  TarifarioTransportista    string          `json:"tarifarioTransportista,omitempty"`
  TarifaCosto               float64         `json:"TarifaCosto,omitempty"`
  CostoViaje                float64         `json:"costoViaje,omitempty"`
  AutCostoViaje_id          bson.ObjectId   `bson:"autCostoViaje_id" json:"autCostoViaje_id,omitempty"` // el usuario que autorizó la tarifa
  AutCosto                  string          `json:"autCosto,omitempty"` // usuario que autorizó
  AutCostoViajeFecha        time.Time       `json:"autCostoViajeFecha,omitempty"`
  Peajes                    float64         `json:"peajes,omitempty"`
  Observaciones             string          `json:"observaciones,omitempty"`
  Estado                    string          `json:"estado,omitempty"` // Ok - Sin Tarifa - Cancelado - Cerrado
  Cancelado_id              bson.ObjectId   `bson:"cancelado_id" json:"cancelado_id,omitempty"` // el usuario que canceló
  CanceladoUsuario          string          `json:"canceladousuario,omitempty"`
  CanceladoFecha            time.Time       `json:"canceladofecha,omitempty"`
  CanceladoObser            string          `json:"canceladoobser,omitempty"`
  Remitos                   bool            `json:"remitos,omitempty"`
  Remitos_id                bson.ObjectId   `bson:"remitos_id" json:"remitos_id,omitempty"` // usuario que recibió los remitos
  RemitosUsuario            string          `json:"remitosusuario,omitempty"`
  RemitosFecha              time.Time       `json:"remitosfecha,omitempty"`
  Factura_id                bson.ObjectId   `bson:"factura_id" json:"factura_id,omitempty"` // id de la factura
  Factura                   string          `json:"factura,omitempty"` // a-006-154648
  FechaFacturacion          time.Time       `json:"fechafacturacion,omitempty"`
  UsuarioFacturacion_id     bson.ObjectId   `bson:"usuariofacturacion_id" json:"usuarioFacturacion_id,omitempty"`
  UsuarioFacturacion        string          `json:"usuariofacturacion,omitempty"`
  Liquidacion_id            bson.ObjectId   `bson:"liquidacion_id" json:"liquidacion_id,omitempty"` // id de la factura
  Liquidacion               string          `json:"liquidacion,omitempty"` // 478
  FechaLiquidacion          time.Time       `json:"fechaliquidacion,omitempty"`
  UsuarioLiquidacion_id     bson.ObjectId   `bson:"usuarioliquidacion_id" json:"usuarioFacturacion_id,omitempty"`
  UsuarioLiquidacion        string          `json:"usuarioliquidacion,omitempty"`
  Editable                  bool            `json:"editable,omitempty"`
	Timestamp	     			      time.Time				`bson:"timestamp" json:"timestamp"`
}

type ParadasServ struct {
	Locacion_id			bson.ObjectId		`bson:"locacion_id" json:"locacion_id,omitempty"`
  Locacion        string          `json:"locacion"`
}

type CanceladoObser struct {
	Observacion			string					`json:"observacion"`
}

type Autorizaciones struct {
	ID							          bson.ObjectId	 	`bson:"_id" json:"id,omitempty"`
	Empresa_id			          bson.ObjectId		`bson:"empresa_id" json:"empresa_id,omitempty"`
	Viaje_id			            bson.ObjectId   `bson:"viaje_id" json:"viaje_id"`
	FechaHora                 time.Time       `json:"fechaHora"`
	Recorrido                 string          `json:"recorrido,omitempty"` // La primera y la última parada indicando si hay intermedias
	Kilometraje               int             `json:"kilometraje"`
	Solicitante_id            bson.ObjectId   `bson:"solicitante_id" json:"solicitante_id,omitempty"`
  Solicitante		            string          `json:"solicitante,omitempty"`
  SolicitanteFecha          time.Time       `json:"solicitantefecha,omitempty"`
	TipoSolicitud							string					`json:"tiposolicitud,omitempty"` // Tarifa Cliente - Tarifa Transportista
	Titular_id				 	      bson.ObjectId		`bson:"titular_id" json:"titular_id,omitempty"`
  Titular			         	    string          `json:"titular"`
	ImporteTarifario					float64         `json:"importetarifario,omitempty"`
	ImporteSugerido						float64         `json:"importesugerido,omitempty"`
	Autorizante_id            bson.ObjectId   `bson:"autorizante_id" json:"solicitante_id,omitempty"`
  Autorizante		            string          `json:"autorizante,omitempty"`
  AutorizanteFecha          time.Time       `json:"autorizantefecha,omitempty"`
	ImporteAutorizado					float64         `json:"importeautorizado,omitempty"`
	Timestamp	     			      time.Time				`bson:"timestamp" json:"timestamp"`
}

type ImporteSugerido struct {
	Importe										float64         `json:"importe"`
}
