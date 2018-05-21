package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Viaje struct {
	ID							          bson.ObjectId	 	`bson:"_id" json:"id,omitempty"`
	Empresa_id			          bson.ObjectId		`bson:"empresa_id" json:"empresa_id"`
  FechaHora                 time.Time       `bson:"fechaHora" json:"fechaHora"`
  Cliente_id	   	          bson.ObjectId		`bson:"cliente_id" json:"cliente_id"`
  Cliente                   string          `bson:"cliente" json:"cliente"`
  TipoUnidad_id	 	          bson.ObjectId		`bson:"tipoUnidad_id" json:"tipoUnidad_id"`
  TipoUnidad                string          `bson:"tipoUnidad" json:"tipoUnidad"`
  Transportista_id	 	      bson.ObjectId		`bson:"transportista_id" json:"transportista_id"`
  Transportista             string          `bson:"transportista" json:"transportista"`
  Unidad_id	        	      bson.ObjectId		`bson:"unidad_id" json:"unidad_id"`
  Unidad                    string          `bson:"unidad" json:"unidad"`
  Personal_id	   	          bson.ObjectId		`bson:"personal_id" json:"personal_id"`
  Personal                  string          `bson:"personal" json:"personal"`
  Paradas	    		  	      []ParadasServ		`bson:"paradas" json:"paradas"`
  Recorrido                 string          `bson:"recorrido" json:"recorrido"` // La primera y la última parada indicando si hay intermedias
  Kilometraje               int             `bson:"kilometraje" json:"kilometraje"`
  TarifarioCliente          string          `bson:"tarifarioCliente" json:"tarifarioCliente"`
  TarifaValor               float64         `bson:"tarifaValor" json:"tarifaValor"`
  ValorViaje                float64         `bson:"valorViaje" json:"valorViaje"`
  AutValorViaje_id          bson.ObjectId   `bson:"autValorViaje_id" json:"autValorViaje_id"` // el usuario que autorizó la tarifa
  AutValor                  string          `bson:"autValor" json:"autValor"` // usuario que autorizó
  AutValorViajeFecha        time.Time       `bson:"autValorViajeFecha" json:"autValorViajeFecha"`
  TarifarioTransportista    string          `bson:"tarifarioTransportista" json:"tarifarioTransportista"`
  TarifaCosto               float64         `bson:"tarifaCosto" json:"tarifaCosto"`
  CostoViaje                float64         `bson:"costoViaje" json:"costoViaje"`
  AutCostoViaje_id          bson.ObjectId   `bson:"autCostoViaje_id" json:"autCostoViaje_id"` // el usuario que autorizó la tarifa
  AutCosto                  string          `bson:"autCosto" json:"autCosto"` // usuario que autorizó
  AutCostoViajeFecha        time.Time       `bson:"autCostoViajeFecha" json:"autCostoViajeFecha"`
  Peajes                    float64         `bson:"peajes" json:"peajes"`
  Observaciones             string          `bson:"observaciones" json:"observaciones"`
  Estado                    string          `bson:"estado" json:"estado"` // Ok - Sin Tarifa - Cancelado - Cerrado
  Cancelado_id              bson.ObjectId   `bson:"cancelado_id" json:"cancelado_id"` // el usuario que canceló
  CanceladoUsuario          string          `bson:"canceladoUsuario" json:"canceladoUsuario"`
  CanceladoFecha            time.Time       `bson:"canceladoFecha" json:"canceladoFecha"`
  CanceladoObser            string          `bson:"canceladoObser" json:"canceladoObser"`
  Remitos                   bool            `bson:"remitos" json:"remitos"`
  Remitos_id                bson.ObjectId   `bson:"remitos_id" json:"remitos_id"` // usuario que recibió los remitos
  RemitosUsuario            string          `bson:"remitosUsuario" json:"remitosUsuario"`
  RemitosFecha              time.Time       `bson:"remitosFecha" json:"remitosFecha"`
  Factura_id                bson.ObjectId   `bson:"factura_id" json:"factura_id"` // id de la factura
  Factura                   string          `bson:"factura" json:"factura"` // a-006-154648
  FechaFacturacion          time.Time       `bson:"fechaFacturacion" json:"fechaFacturacion"`
  UsuarioFacturacion_id     bson.ObjectId   `bson:"usuarioFacturacion_id" json:"usuarioFacturacion_id"`
  UsuarioFacturacion        string          `bson:"usuarioFacturacion" json:"usuarioFacturacion"`
  Liquidacion_id            bson.ObjectId   `bson:"liquidacion_id" json:"liquidacion_id"` // id de la factura
  Liquidacion               string          `bson:"liquidacion" json:"liquidacion"` // 478
  FechaLiquidacion          time.Time       `bson:"fechaLiquidacion" json:"fechaLiquidacion"`
  UsuarioLiquidacion_id     bson.ObjectId   `bson:"usuarioLiquidacion_id" json:"usuarioLiquidacion_id"`
  UsuarioLiquidacion        string          `bson:"usuarioLiquidacion" json:"usuarioLiquidacion"`
  Editable                  bool            `bson:"editable" json:"editable"`
	Timestamp	     			      time.Time				`bson:"timestamp" json:"timestamp"`
}

type ParadasServ struct {
	Locacion_id			bson.ObjectId		`bson:"locacion_id" json:"locacion_id,omitempty"`
  Locacion        string          `bson:"locacion" json:"locacion"`
}

type CanceladoObser struct {
	Observacion			string					`bson:"observacion" json:"observacion"`
}

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
	Timestamp	     			      time.Time				`bson:"timestamp" json:"timestamp"`
}

type ImporteSugerido struct {
	Importe										float64         `bson:"importe" json:"importe"`
}

/*
type Viaje struct {
	ID							          bson.ObjectId	 	`bson:"_id" json:"id,omitempty"`
	Empresa_id			          bson.ObjectId		`bson:"empresa_id" json:"empresa_id,omitempty"`
  FechaHora                 time.Time       `bson:"fechaHora" json:"fechaHora"`
  Cliente_id	   	          bson.ObjectId		`bson:"cliente_id,omitempty" json:"cliente_id,omitempty"`
  Cliente                   string          `bson:"cliente,omitempty" json:"cliente"`
  TipoUnidad_id	 	          bson.ObjectId		`bson:"tipoUnidad_id,omitempty" json:"tipoUnidad_id,omitempty"`
  TipoUnidad                string          `bson:"tipoUnidad,omitempty" json:"tipoUnidad"`
  Transportista_id	 	      bson.ObjectId		`bson:"transportista_id,omitempty" json:"transportista_id,omitempty"`
  Transportista             string          `bson:"transportista,omitempty" json:"transportista"`
  Unidad_id	        	      bson.ObjectId		`bson:"unidad_id" json:"unidad_id,omitempty"`
  Unidad                    string          `bson:"unidad" json:"unidad"`
  Personal_id	   	          bson.ObjectId		`bson:"personal_id" json:"personal_id,omitempty"`
  Personal                  string          `bson:"personal" json:"personal"`
  Paradas	    		  	      []ParadasServ		`bson:"paradas" json:"paradas,omitempty"`
  Recorrido                 string          `bson:"recorrido" json:"recorrido,omitempty"` // La primera y la última parada indicando si hay intermedias
  Kilometraje               int             `bson:"kilometraje" json:"kilometraje,omitempty"`
  TarifarioCliente          string          `bson:"tarifarioCliente,omitempty" json:"tarifarioCliente,omitempty"`
  TarifaValor               float64         `bson:"tarifaValor,omitempty" json:"tarifaValor,omitempty"`
  ValorViaje                float64         `bson:"valorViaje,omitempty" json:"valorViaje,omitempty"`
  AutValorViaje_id          bson.ObjectId   `bson:"autValorViaje_id,omitempty" json:"autValorViaje_id,omitempty"` // el usuario que autorizó la tarifa
  AutValor                  string          `bson:"autValor,omitempty" json:"autValor,omitempty"` // usuario que autorizó
  AutValorViajeFecha        time.Time       `bson:"autValorViajeFecha,omitempty" json:"autValorViajeFecha,omitempty"`
  TarifarioTransportista    string          `bson:"tarifarioTransportista,omitempty" json:"tarifarioTransportista,omitempty"`
  TarifaCosto               float64         `bson:"tarifaCosto,omitempty" json:"tarifaCosto,omitempty"`
  CostoViaje                float64         `bson:"costoViaje,omitempty" json:"costoViaje,omitempty"`
  AutCostoViaje_id          bson.ObjectId   `bson:"autCostoViaje_id,omitempty" json:"autCostoViaje_id,omitempty"` // el usuario que autorizó la tarifa
  AutCosto                  string          `bson:"autCosto,omitempty" json:"autCosto,omitempty"` // usuario que autorizó
  AutCostoViajeFecha        time.Time       `bson:"autCostoViajeFecha,omitempty" json:"autCostoViajeFecha,omitempty"`
  Peajes                    float64         `bson:"peajes,omitempty" json:"peajes,omitempty"`
  Observaciones             string          `bson:"observaciones,omitempty" json:"observaciones,omitempty"`
  Estado                    string          `bson:"estado,omitempty" json:"estado,omitempty"` // Ok - Sin Tarifa - Cancelado - Cerrado
  Cancelado_id              bson.ObjectId   `bson:"cancelado_id,omitempty" json:"cancelado_id,omitempty"` // el usuario que canceló
  CanceladoUsuario          string          `bson:"canceladoUsuario,omitempty" json:"canceladoUsuario,omitempty"`
  CanceladoFecha            time.Time       `bson:"canceladoFecha,omitempty" json:"canceladoFecha,omitempty"`
  CanceladoObser            string          `bson:"canceladoObser,omitempty" json:"canceladoObser,omitempty"`
  Remitos                   bool            `bson:"remitos,omitempty" json:"remitos,omitempty"`
  Remitos_id                bson.ObjectId   `bson:"remitos_id,omitempty" json:"remitos_id,omitempty"` // usuario que recibió los remitos
  RemitosUsuario            string          `bson:"remitosUsuario,omitempty" json:"remitosUsuario,omitempty"`
  RemitosFecha              time.Time       `bson:"remitosFecha,omitempty" json:"remitosFecha,omitempty"`
  Factura_id                bson.ObjectId   `bson:"factura_id,omitempty" json:"factura_id,omitempty"` // id de la factura
  Factura                   string          `bson:"factura,omitempty" json:"factura,omitempty"` // a-006-154648
  FechaFacturacion          time.Time       `bson:"fechaFacturacion,omitempty" json:"fechaFacturacion,omitempty"`
  UsuarioFacturacion_id     bson.ObjectId   `bson:"usuarioFacturacion_id,omitempty" json:"usuarioFacturacion_id,omitempty"`
  UsuarioFacturacion        string          `bson:"usuarioFacturacion,omitempty" json:"usuarioFacturacion,omitempty"`
  Liquidacion_id            bson.ObjectId   `bson:"liquidacion_id,omitempty" json:"liquidacion_id,omitempty"` // id de la factura
  Liquidacion               string          `bson:"liquidacion,omitempty" json:"liquidacion,omitempty"` // 478
  FechaLiquidacion          time.Time       `bson:"fechaLiquidacion,omitempty" json:"fechaLiquidacion,omitempty"`
  UsuarioLiquidacion_id     bson.ObjectId   `bson:"usuarioLiquidacion_id,omitempty" json:"usuarioLiquidacion_id,omitempty"`
  UsuarioLiquidacion        string          `bson:"usuarioLiquidacion,omitempty" json:"usuarioLiquidacion,omitempty"`
  Editable                  bool            `bson:"editable,omitempty" json:"editable,omitempty"`
	Timestamp	     			      time.Time				`bson:"timestamp,omitempty" json:"timestamp"`
}

*/
