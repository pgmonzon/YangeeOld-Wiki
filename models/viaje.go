package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Viaje struct {
	ID							          bson.ObjectId	 	`bson:"_id" json:"id,omitempty"`
	Empresa_id			          bson.ObjectId		`bson:"empresa_id" json:"empresa_id,omitempty"`
  FechaHora                 time.Time       `bson:"fechaHora" json:"fechaHora"`
	Vuelta										string					`bson:"vuelta" json:"vuelta"` //  [vacío] - 2da - 3ra - 4ta
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
	Celular										string					`bson:"celular" json:"celular"`
	Tipo											string					`bson:"tipo" json:"tipo"` // chofer - asistente
  Paradas	    		  	      []ParadasServ		`bson:"paradas" json:"paradas"`
  Recorrido                 string          `bson:"recorrido" json:"recorrido"` // La primera y la última parada indicando si hay intermedias
  Kilometraje               int             `bson:"kilometraje" json:"kilometraje"`
	Regreso										int							`bson:"regreso" json:"regreso"` // kilometraje para el regreso de la unidad
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
	Adicional									string					`bson:"adicional" json:"adicional"`
	Importe										float64					`bson:"importe" json:"importe"`
  Observaciones             string          `bson:"observaciones" json:"observaciones"`
  Estado                    string          `bson:"estado" json:"estado"` // Ok - Sin Tarifa - Cancelado - Cerrado
  Cancelado_id              bson.ObjectId   `bson:"cancelado_id" json:"cancelado_id"` // el usuario que canceló
  CanceladoUsuario          string          `bson:"canceladoUsuario" json:"canceladoUsuario"`
  CanceladoFecha            time.Time       `bson:"canceladoFecha" json:"canceladoFecha"`
  CanceladoObser            string          `bson:"canceladoObser" json:"canceladoObser"`
  Remitos                   bool            `bson:"remitos" json:"remitos"`
	RemitosDetalle						string					`bson:"remitosDetalle" json:"remitosDetalle"`
  Remitos_id                bson.ObjectId   `bson:"remitos_id" json:"remitos_id"` // usuario que recibió los remitos
  RemitosUsuario            string          `bson:"remitosUsuario" json:"remitosUsuario"`
  RemitosFecha              time.Time       `bson:"remitosFecha" json:"remitosFecha"`
  Factura_id                bson.ObjectId   `bson:"factura_id" json:"factura_id"` // id de la factura
  Factura                   string          `bson:"factura" json:"factura"` // a-006-154648
  FechaFacturacion          time.Time       `bson:"fechaFacturacion" json:"fechaFacturacion"`
  UsuarioFacturacion_id     bson.ObjectId   `bson:"usuarioFacturacion_id" json:"usuarioFacturacion_id"`
  UsuarioFacturacion        string          `bson:"usuarioFacturacion" json:"usuarioFacturacion"`
  Liquidacion_id            bson.ObjectId   `bson:"liquidacion_id" json:"liquidacion_id"` // id de la factura
  Liquidacion               int		          `bson:"liquidacion" json:"liquidacion"` // 478
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

type ViajesFechas struct {
	FechaDesde			time.Time       `bson:"fechaDesde" json:"fechaDesde"`
	FechaHasta			time.Time       `bson:"fechaHasta" json:"fechaHasta"`
}
