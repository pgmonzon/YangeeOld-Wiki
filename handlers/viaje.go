package handlers

import (
  "encoding/json"
  "net/http"
  "strings"
  "time"
  "strconv"

  "github.com/pgmonzon/Yangee/models"
  "github.com/pgmonzon/Yangee/core"
  "github.com/pgmonzon/Yangee/config"

  "gopkg.in/mgo.v2/bson"
  "github.com/gorilla/context"
  "github.com/gorilla/mux"
)

func ViajeCrear(w http.ResponseWriter, req *http.Request) {
  //-------------------Modificar ###### estas 2 variables
	var documento models.Viaje
  audit := "Crear"

  // Decode del JSON
  // ***************
  decoder := json.NewDecoder(req.Body)
  err := decoder.Decode(&documento)
  if err != nil {
    core.RspMsgJSON(w, req, "ERROR", "JSON", "INVALID_PARAMS: JSON decode erróneo", http.StatusBadRequest)
    return
  }

  // Doy de alta
  // ***********
  //------------------------------------------------------Modificar ######
  estado, valor, mensaje, httpStat, documento := ViajeAlta(documento, req, audit)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }

  // Está todo Ok
  // ************
  s := []string{"Agregaste un viaje para ", documento.Cliente}
  core.RspMsgJSON(w, req, "OK", documento.Cliente, strings.Join(s, ""), http.StatusCreated)
  return
}

// Devuelve Estado, Valor, Mensaje, HttpStat, Collection, Existía
func ViajeAlta(documentoAlta models.Viaje, req *http.Request, audit string) (string, string, string, int, models.Viaje) {
  //-------------------Modificar ###### las 3 variables
	var documento models.Viaje
  //camposVacios := "No podés dejar vacío el campo Hora"
  coll := config.DB_Viaje
  empresaID := context.Get(req, "Empresa_id").(bson.ObjectId)

  // Obtengo la tarifa del cliente
  // *****************************
  tarifarioCliente, tarifaValor := TarifaCliente(documentoAlta, req)

  // Obtengo la tarifa del transportista
  // ***********************************
  tarifarioTransportista, tarifaCosto := TarifaTransportista(documentoAlta, req)

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    return "ERROR", "GetMongoSession", err.Error(), httpStat, documento
  }
  defer session.Close()

  // Intento el alta
  // ***************
  documento = documentoAlta
  objID := bson.NewObjectId()
  documento.ID = objID
  documento.Empresa_id = empresaID
  documento.Timestamp = time.Now()
  documento.Recorrido = RecorridoPuntas(documento)
  documento.TarifarioCliente = tarifarioCliente
  documento.TarifaValor = tarifaValor
  documento.ValorViaje = tarifaValor // en el alta va el mismo que el tarifado
  documento.AutValorViaje_id = config.FakeID // por defecto lo pongo en vacío el usuario
  documento.AutValor = "" // por defecto en vacío el usuario
  documento.AutValorViajeFecha = time.Time{} // por defecto en vacío
  documento.TarifarioTransportista = tarifarioTransportista
  documento.TarifaCosto = tarifaCosto
  documento.CostoViaje = tarifaCosto // en el alta va el mismo que el tarifado
  documento.AutCostoViaje_id = config.FakeID // por defecto lo pongo en vacío el usuario
  documento.AutCosto = "" // por defecto en vacío el usuario
  documento.AutCostoViajeFecha = time.Time{} // por defecto en vacío
  documento.Estado = "Ok" // Ok - Cancelado - Cerrado
  documento.Cancelado_id = config.FakeID // por defecto lo pongo en vacío el usuario
  documento.CanceladoUsuario = "" // por defecto en vacío
  documento.CanceladoFecha = time.Time{} // por defecto en vacío
  documento.CanceladoObser = "" // por defecto en vacío
  documento.Remitos = false // por defecto en false cuando se reciben los remitos cambia
  documento.Remitos_id = config.FakeID // por defecto lo pongo en vacío el usuario
  documento.RemitosUsuario = ""
  documento.RemitosFecha = time.Time{}
  documento.Factura_id = config.FakeID
  documento.Factura = ""
  documento.FechaFacturacion = time.Time{}
  documento.UsuarioFacturacion_id = config.FakeID
  documento.UsuarioFacturacion = ""
  documento.Liquidacion_id = config.FakeID
  documento.Liquidacion = ""
  documento.FechaLiquidacion = time.Time{}
  documento.UsuarioLiquidacion_id = config.FakeID
  documento.UsuarioLiquidacion = ""
  documento.Editable = true

  collection := session.DB(config.DB_Name).C(coll)
  err = collection.Insert(documento)
  if err != nil {
    s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
    return "ERROR", audit, strings.Join(s, ""), http.StatusInternalServerError, documento
  }

  // Está todo Ok
  // ************
  core.Audit(req, coll, documento.ID, audit, documento)
  return "OK", audit, "Ok", http.StatusOK, documento
}

func RecorridoPuntas(documento models.Viaje) (string) {
  origen := ""
  destino := ""
  cant := 0
  s := []string{""}

  for _, item := range documento.Paradas {
    cant = cant + 1
    if origen == "" {
      origen = item.Locacion
    } else {
      destino = item.Locacion
    }
  }

  if cant > 2 {
    s = []string{origen, " <> ", destino}
  } else {
    s = []string{origen, " -> ", destino}
  }

  return strings.Join(s, "")
}

func TarifaCliente(documento models.Viaje, req *http.Request) (string, float64) {
  var tarifarioCliente string
  var tarifaValor float64
  tarifarioCliente = ""
  tarifaValor = 0
  ahora := time.Now()

  // Busco el cliente por los tarifarios
  // ***********************************
  _, _, _, httpStat, cliente := Cliente_X_ID(documento.Cliente_id, "Buscar ID", req)
  if httpStat != http.StatusOK {
    return tarifarioCliente, tarifaValor
  }

  // Recorro los tarifarios para ver cual aplica en orden
  // ****************************************************
  for _, item := range cliente.Tarifarios {
    if item.Activo == true && item.Importe > 0 && (ahora.Equal(item.VigenteDesde) || ahora.After(item.VigenteDesde)) && (ahora.Equal(item.VigenteHasta) || ahora.Before(item.VigenteHasta)) {
      if strings.ToUpper(item.Tipo) == "KILOMETRAJE" {
        if documento.Kilometraje > 0 {
          if item.TipoUnidad_id == config.FakeID || item.TipoUnidad_id == documento.TipoUnidad_id {
            tarifarioCliente = item.Tarifario
            tarifaValor = item.Importe * float64(documento.Kilometraje)
            return tarifarioCliente, tarifaValor
          }
        }
      }
      if strings.ToUpper(item.Tipo) == "RANGO KILOMETRAJE" {
        if documento.Kilometraje > 0 && documento.Kilometraje >= item.KmDesde && documento.Kilometraje <= item.KmHasta {
          if item.TipoUnidad_id == config.FakeID || item.TipoUnidad_id == documento.TipoUnidad_id {
            tarifarioCliente = item.Tarifario
            tarifaValor = item.Importe
            return tarifarioCliente, tarifaValor
          }
        }
      }
      if strings.ToUpper(item.Tipo) == "RECORRIDO" {
        if item.TipoUnidad_id == config.FakeID || item.TipoUnidad_id == documento.TipoUnidad_id {
          encontrados := 0
          for _, itemTar := range item.Recorrido {
            encontrados = 0
            for _, itemVia := range documento.Paradas {
              if itemTar.Locacion_id == itemVia.Locacion_id {
                encontrados = 1
                break
              }
            }
            if encontrados == 0 {
              break
            }
          }
          if encontrados == 1 {
            for _, itemVia := range documento.Paradas {
              encontrados = 0
              for _, itemTar := range item.Recorrido {
                if itemTar.Locacion_id == itemVia.Locacion_id {
                  encontrados = 1
                  break
                }
              }
              if encontrados == 0 {
                break
              }
            }
          }
          if encontrados == 1 {
            tarifarioCliente = item.Tarifario
            tarifaValor = item.Importe
            return tarifarioCliente, tarifaValor
          }
        }
      }
    }
  }
  return tarifarioCliente, tarifaValor
}

func TarifaTransportista(documento models.Viaje, req *http.Request) (string, float64) {
  var tarifarioTransportista string
  var tarifaCosto float64
  tarifarioTransportista = ""
  tarifaCosto = 0
  ahora := time.Now()

  // Busco el transportista por los tarifarios
  // *****************************************
  _, _, _, httpStat, transportista := Transportista_X_ID(documento.Transportista_id, "Buscar ID", req)
  if httpStat != http.StatusOK {
    return tarifarioTransportista, tarifaCosto
  }

  // Recorro los tarifarios para ver cual aplica en orden
  // ****************************************************
  for _, item := range transportista.Tarifarios {
    if item.Activo == true && item.Importe > 0 && (ahora.Equal(item.VigenteDesde) || ahora.After(item.VigenteDesde)) && (ahora.Equal(item.VigenteHasta) || ahora.Before(item.VigenteHasta)) {
      if strings.ToUpper(item.Tipo) == "KILOMETRAJE" {
        if documento.Kilometraje > 0 {
          if item.TipoUnidad_id == config.FakeID || item.TipoUnidad_id == documento.TipoUnidad_id {
            tarifarioTransportista = item.Tarifario
            tarifaCosto = item.Importe * float64(documento.Kilometraje)
            return tarifarioTransportista, tarifaCosto
          }
        }
      }
      if strings.ToUpper(item.Tipo) == "RANGO KILOMETRAJE" {
        if documento.Kilometraje > 0 && documento.Kilometraje >= item.KmDesde && documento.Kilometraje <= item.KmHasta {
          if item.TipoUnidad_id == config.FakeID || item.TipoUnidad_id == documento.TipoUnidad_id {
            tarifarioTransportista = item.Tarifario
            tarifaCosto = item.Importe
            return tarifarioTransportista, tarifaCosto
          }
        }
      }
      if strings.ToUpper(item.Tipo) == "RECORRIDO" {
        if item.TipoUnidad_id == config.FakeID || item.TipoUnidad_id == documento.TipoUnidad_id {
          encontrados := 0
          for _, itemTar := range item.Recorrido {
            encontrados = 0
            for _, itemVia := range documento.Paradas {
              if itemTar.Locacion_id == itemVia.Locacion_id {
                encontrados = 1
                break
              }
            }
            if encontrados == 0 {
              break
            }
          }
          if encontrados == 1 {
            for _, itemVia := range documento.Paradas {
              encontrados = 0
              for _, itemTar := range item.Recorrido {
                if itemTar.Locacion_id == itemVia.Locacion_id {
                  encontrados = 1
                  break
                }
              }
              if encontrados == 0 {
                break
              }
            }
          }
          if encontrados == 1 {
            tarifarioTransportista = item.Tarifario
            tarifaCosto = item.Importe
            return tarifarioTransportista, tarifaCosto
          }
        }
      }
    }
  }
  return tarifarioTransportista, tarifaCosto
}

func ViajeGuardar(w http.ResponseWriter, req *http.Request) {
  //-------------------Modificar ######
  var documento models.Viaje
  vars := mux.Vars(req)
  ID := vars["docID"]
  audit := "Guardar"

  // Verifico el formato del campo ID
  // ********************************
  if bson.IsObjectIdHex(ID) != true {
    core.RspMsgJSON(w, req, "ERROR", ID, "INVALID_PARAMS: Formato ID incorrecto", http.StatusBadRequest)
    return
  }
  documentoID := bson.ObjectIdHex(ID)

  // Decode del JSON
  // ***************
  decoder := json.NewDecoder(req.Body)
  err := decoder.Decode(&documento)
  if err != nil {
    core.RspMsgJSON(w, req, "ERROR", "JSON", "INVALID_PARAMS: JSON decode erróneo", http.StatusBadRequest)
    return
  }

  // Modifico
  // ********
  //----------------------------------Modificar ######
  estado, valor, mensaje, httpStat := ViajeModificar(documentoID, documento, req, audit)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }

  // Está todo Ok
  // ************
  s := []string{"Guardaste el viaje de ", documento.Cliente}
  core.RspMsgJSON(w, req, "OK", documento.Cliente, strings.Join(s, ""), http.StatusAccepted)
  return
}

func ViajeModificar(documentoID bson.ObjectId, documentoModi models.Viaje, req *http.Request, audit string) (string, string, string, int) {
  coll := config.DB_Viaje
  empresaID := context.Get(req, "Empresa_id").(bson.ObjectId)

  // Busco para obtener los campos del viaje original
  // ************************************************
  //------------------------------------------------------Modificar ######
  estado, valor, mensaje, httpStat, documentoExistente := Viaje_X_ID(documentoID, "Buscar ID", req)
  if httpStat != http.StatusOK {
    return estado, valor, mensaje, httpStat
  }

  // Me fijo si está Editable
  // ************************
  if documentoExistente.Editable == false {
    s := []string{"El viaje está bloqueado"}
    return "ERROR", "Modificar", strings.Join(s, ""), http.StatusNonAuthoritativeInfo
  }

  // Si cambia alguno de los valores que modifica las tarifas
  // la recalculo y blanqueo las autorizaciones
  // ********************************************************

  if documentoExistente.Kilometraje != documentoModi.Kilometraje || documentoExistente.TipoUnidad_id != documentoModi.TipoUnidad_id || documentoExistente.Recorrido != documentoModi.Recorrido || documentoExistente.Cliente_id != documentoModi.Cliente_id || documentoExistente.Transportista_id != documentoModi.Transportista_id {
    // Obtengo la tarifa del cliente
    // *****************************
    tarifarioCliente, tarifaValor := TarifaCliente(documentoModi, req)

    // Obtengo la tarifa del transportista
    // ***********************************
    tarifarioTransportista, tarifaCosto := TarifaTransportista(documentoModi, req)

    documentoModi.TarifarioCliente = tarifarioCliente
    documentoModi.TarifaValor = tarifaValor
    documentoModi.ValorViaje = tarifaValor // en el alta va el mismo que el tarifado
    documentoModi.AutValorViaje_id = config.FakeID // por defecto lo pongo en vacío el usuario
    documentoModi.AutValor = "" // por defecto en vacío el usuario
    documentoModi.AutValorViajeFecha = time.Time{} // por defecto en vacío
    documentoModi.TarifarioTransportista = tarifarioTransportista
    documentoModi.TarifaCosto = tarifaCosto
    documentoModi.CostoViaje = tarifaCosto // en el alta va el mismo que el tarifado
    documentoModi.AutCostoViaje_id = config.FakeID // por defecto lo pongo en vacío el usuario
    documentoModi.AutCosto = "" // por defecto en vacío el usuario
    documentoModi.AutCostoViajeFecha = time.Time{} // por defecto en vacío
  } else {
    documentoModi.TarifarioCliente = documentoExistente.TarifarioCliente
    documentoModi.TarifaValor = documentoExistente.TarifaValor
    documentoModi.ValorViaje = documentoExistente.ValorViaje
    documentoModi.AutValorViaje_id = documentoExistente.AutValorViaje_id
    documentoModi.AutValor = documentoExistente.AutValor
    documentoModi.AutValorViajeFecha = documentoExistente.AutValorViajeFecha
    documentoModi.TarifarioTransportista = documentoExistente.TarifarioTransportista
    documentoModi.TarifaCosto = documentoExistente.TarifaCosto
    documentoModi.CostoViaje = documentoExistente.CostoViaje
    documentoModi.AutCostoViaje_id = documentoExistente.AutCostoViaje_id
    documentoModi.AutCosto = documentoExistente.AutCosto
    documentoModi.AutCostoViajeFecha = documentoExistente.AutCostoViajeFecha
  }

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    return "ERROR", "GetMongoSession", err.Error(), httpStat
  }
  defer session.Close()

  // Intento la modificación
  // ***********************
  documentoModi.ID = documentoID
  documentoModi.Empresa_id = empresaID
  collection := session.DB(config.DB_Name).C(coll)
  selector := bson.M{"_id": documentoID, "empresa_id": empresaID}
  updator := bson.M{
    "$set": bson.M{
      "fechaHora": documentoModi.FechaHora,
      "cliente_id": documentoModi.Cliente_id,
      "cliente": documentoModi.Cliente,
      "tipoUnidad_id": documentoModi.TipoUnidad_id,
      "tipoUnidad": documentoModi.TipoUnidad,
      "transportista_id": documentoModi.Transportista_id,
      "transportista": documentoModi.Transportista,
      "unidad_id": documentoModi.Unidad_id,
      "unidad": documentoModi.Unidad,
      "personal_id": documentoModi.Personal_id,
      "personal": documentoModi.Personal,
      "paradas": documentoModi.Paradas,
      "recorrido": RecorridoPuntas(documentoModi),
      "kilometraje": documentoModi.Kilometraje,
      "peajes": documentoModi.Peajes,
      "observaciones": documentoModi.Observaciones,
      "tarifarioCliente": documentoModi.TarifarioCliente,
      "tarifaValor": documentoModi.TarifaValor,
      "valorViaje": documentoModi.ValorViaje,
      "autValorViaje_id": documentoModi.AutValorViaje_id,
      "autValor": documentoModi.AutValor,
      "autValorViajeFecha": documentoModi.AutValorViajeFecha,
      "tarifarioTransportista": documentoModi.TarifarioTransportista,
      "tarifaCosto": documentoModi.TarifaCosto,
      "costoViaje": documentoModi.CostoViaje,
      "autCostoViaje_id": documentoModi.AutCostoViaje_id,
      "autCosto": documentoModi.AutCosto,
      "autCostoViajeFecha": documentoModi.AutCostoViajeFecha,
      "timestamp": time.Now(),
    },
  }
  err = collection.Update(selector, updator)
  if err != nil {
    s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
    return "ERROR", audit, strings.Join(s, ""), http.StatusInternalServerError
  }

  // Está todo Ok
  // ************
  core.Audit(req, coll, documentoID, audit, documentoModi)
  return "OK", audit, "Ok", http.StatusOK
}

// Devuelve Estado, Valor, Mensaje, HttpStat, collection
func Viaje_X_ID(documentoID bson.ObjectId, audit string, req *http.Request) (string, string, string, int, models.Viaje) {
  //-------------------Modificar ######
  var documento models.Viaje
  coll := config.DB_Viaje
  empresaID := context.Get(req, "Empresa_id").(bson.ObjectId)

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    return "ERROR", "GetMongoSession", err.Error(), httpStat, documento
  }
  defer session.Close()

  // Trato de traerlo
  // ****************
  collection := session.DB(config.DB_Name).C(coll)
  collection.Find(bson.M{"_id": documentoID, "empresa_id": empresaID}).Select(bson.M{"empresa_id":0}).One(&documento)
  // No existe
  if documento.ID == "" {
    s := []string{"No encuentro el documento"}
    return "ERROR", audit, strings.Join(s, ""), http.StatusNonAuthoritativeInfo, documento
  }
  // Existe
  return "OK", audit, "Ok", http.StatusOK, documento
}

func ViajesTraer(w http.ResponseWriter, req *http.Request) {
  //-------------------Modificar ###### estas 2 variables
  var documento models.Viaje
  var documentos []models.Viaje
  vars := mux.Vars(req)
  ano, _ := strconv.Atoi(vars["ano"])
  mes, _ := strconv.Atoi(vars["mes"])
  dia, _ := strconv.Atoi(vars["dia"])

  // Busco
  // *****
  //----------------------------------------------Modificar ######
  estado, valor, mensaje, httpStat, documentos := ViajesBuscar(documento, ano, mes, dia, "Buscar", req)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }

  // Está todo Ok
  // ************
  respuesta, error := json.Marshal(documentos)
  core.FatalErr(error)
  core.RspJSON(w, req, respuesta, http.StatusOK)
  return
}

// Devuelve Estado, Valor, Mensaje, HttpStat, Collection
func ViajesBuscar(documento models.Viaje, ano int, mes int, dia int, audit string, req *http.Request) (string, string, string, int, []models.Viaje) {
  //----------------------Modificar ###### estas 2 variables
  var documentos []models.Viaje
  coll := config.DB_Viaje
  empresaID := context.Get(req, "Empresa_id").(bson.ObjectId)
  //fechaDesde := time.Date(ano, time.Month(mes), dia, 0, 0, 0, 0, time.UTC)
  //fechaHasta := time.Date(ano, time.Month(mes), dia, 23, 59, 59, 999999999, time.UTC)

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    return "ERROR", "GetMongoSession", err.Error(), httpStat, documentos
  }
  defer session.Close()

  // Trato de traerlos
  // *****************
  //----------Modificar ###### en forma manual
  selector := bson.M{
    "empresa_id": empresaID,
    //"fechaHora": bson.M{"$gte": fechaDesde, "$lte": fechaHasta},
  }
  collection := session.DB(config.DB_Name).C(coll)
  collection.Find(selector).Select(bson.M{"empresa_id":0}).All(&documentos)

  // Si el resultado es vacío devuelvo ERROR
  // ***************************************
  if len(documentos) == 0 {
    s := []string{"No encontré documentos"}
    return "ERROR", audit, strings.Join(s, ""), http.StatusNonAuthoritativeInfo, documentos
  }

  // Está todo Ok
  // ************
  return "OK", audit, "Ok", http.StatusOK, documentos
}

func ViajeCancelar(w http.ResponseWriter, req *http.Request) {
  //-------------------Modificar ######
  var obser models.CanceladoObser
  vars := mux.Vars(req)
  ID := vars["docID"]
  audit := "Cancelar"

  // Verifico el formato del campo ID
  // ********************************
  if bson.IsObjectIdHex(ID) != true {
    core.RspMsgJSON(w, req, "ERROR", ID, "INVALID_PARAMS: Formato ID incorrecto", http.StatusBadRequest)
    return
  }
  documentoID := bson.ObjectIdHex(ID)

  // Decode del JSON
  // ***************
  decoder := json.NewDecoder(req.Body)
  err := decoder.Decode(&obser)
  if err != nil {
    core.RspMsgJSON(w, req, "ERROR", "JSON", "INVALID_PARAMS: JSON decode erróneo", http.StatusBadRequest)
    return
  }

  // Modifico
  // ********
  //----------------------------------Modificar ######
  estado, valor, mensaje, httpStat := ViajeCancel(documentoID, obser, req, audit)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }

  // Está todo Ok
  // ************
  //------------------------------------Modificar ######
  s := []string{"Cancelaste ", obser.Observacion}
  //--------------------------------------Modificar ######
  core.RspMsgJSON(w, req, "OK", obser.Observacion, strings.Join(s, ""), http.StatusAccepted)
  return
}

func ViajeCancel(documentoID bson.ObjectId, obser models.CanceladoObser, req *http.Request, audit string) (string, string, string, int) {
  //-------------------Modificar ###### las 2 variables
  coll := config.DB_Viaje
  empresaID := context.Get(req, "Empresa_id").(bson.ObjectId)
  usuarioID := context.Get(req, "Usuario_id").(bson.ObjectId)
  usuario := context.Get(req, "Usuario")

  // Busco para obtener los campos del viaje original
  // ************************************************
  //------------------------------------------------------Modificar ######
  estado, valor, mensaje, httpStat, documentoExistente := Viaje_X_ID(documentoID, "Buscar ID", req)
  if httpStat != http.StatusOK {
    return estado, valor, mensaje, httpStat
  }

  // Me fijo si está Editable
  // ************************
  if documentoExistente.Editable == false {
    s := []string{"El viaje está bloqueado"}
    return "ERROR", "Modificar", strings.Join(s, ""), http.StatusNonAuthoritativeInfo
  }

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    return "ERROR", "GetMongoSession", err.Error(), httpStat
  }
  defer session.Close()

  // Intento la modificación
  // ***********************
  //-------------------Modificar ###### en forma manual
  collection := session.DB(config.DB_Name).C(coll)
  selector := bson.M{"_id": documentoID, "empresa_id": empresaID}
  updator := bson.M{
    "$set": bson.M{
      "estado": "Cancelado",
      "cancelado_id": usuarioID,
      "canceladousuario": usuario,
      "canceladofecha": time.Now(),
      "canceladoobser": obser.Observacion,
      "editable": false,
      "timestamp": time.Now(),
    },
  }
  err = collection.Update(selector, updator)
  if err != nil {
    s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
    return "ERROR", audit, strings.Join(s, ""), http.StatusInternalServerError
  }

  // Está todo Ok
  // ************
  core.Audit(req, coll, documentoID, audit, obser)
  return "OK", audit, "Ok", http.StatusOK
}

func ViajeRemitos(w http.ResponseWriter, req *http.Request) {
  //-------------------Modificar ######
  vars := mux.Vars(req)
  ID := vars["docID"]
  audit := "Remitos Recibidos"

  // Verifico el formato del campo ID
  // ********************************
  if bson.IsObjectIdHex(ID) != true {
    core.RspMsgJSON(w, req, "ERROR", ID, "INVALID_PARAMS: Formato ID incorrecto", http.StatusBadRequest)
    return
  }
  documentoID := bson.ObjectIdHex(ID)

  // Modifico
  // ********
  //----------------------------------Modificar ######
  estado, valor, mensaje, httpStat := ViajeRemitosRecibidos(documentoID, req, audit)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }

  // Está todo Ok
  // ************
  //------------------------------------Modificar ######
  s := []string{"Recibiste los remitos"}
  //--------------------------------------Modificar ######
  core.RspMsgJSON(w, req, "OK", "Recibiste los remitos", strings.Join(s, ""), http.StatusAccepted)
  return
}

func ViajeRemitosRecibidos(documentoID bson.ObjectId, req *http.Request, audit string) (string, string, string, int) {
  //-------------------Modificar ###### las 2 variables
  coll := config.DB_Viaje
  empresaID := context.Get(req, "Empresa_id").(bson.ObjectId)
  usuarioID := context.Get(req, "Usuario_id").(bson.ObjectId)
  usuario := context.Get(req, "Usuario")

  // Busco para obtener los campos del viaje original
  // ************************************************
  //------------------------------------------------------Modificar ######
  estado, valor, mensaje, httpStat, documentoExistente := Viaje_X_ID(documentoID, "Buscar ID", req)
  if httpStat != http.StatusOK {
    return estado, valor, mensaje, httpStat
  }

  // Me fijo si está Editable
  // ************************
  if documentoExistente.Editable == false {
    s := []string{"El viaje está bloqueado"}
    return "ERROR", "Modificar", strings.Join(s, ""), http.StatusNonAuthoritativeInfo
  }

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    return "ERROR", "GetMongoSession", err.Error(), httpStat
  }
  defer session.Close()

  // Intento la modificación
  // ***********************
  //-------------------Modificar ###### en forma manual
  collection := session.DB(config.DB_Name).C(coll)
  selector := bson.M{"_id": documentoID, "empresa_id": empresaID}
  updator := bson.M{
    "$set": bson.M{
      "remitos": true,
      "remitos_id": usuarioID,
      "remitosusuario": usuario,
      "remitosfecha": time.Now(),
      "timestamp": time.Now(),
    },
  }
  err = collection.Update(selector, updator)
  if err != nil {
    s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
    return "ERROR", audit, strings.Join(s, ""), http.StatusInternalServerError
  }

  // Está todo Ok
  // ************
  core.Audit(req, coll, documentoID, audit, "Recibiste los remitos")
  return "OK", audit, "Ok", http.StatusOK
}

func ViajeAutValor(w http.ResponseWriter, req *http.Request) {
  //-------------------Modificar ######
  var importe models.ImporteSugerido
  vars := mux.Vars(req)
  ID := vars["docID"]
  audit := "Autorizar Valor"

  // Verifico el formato del campo ID
  // ********************************
  if bson.IsObjectIdHex(ID) != true {
    core.RspMsgJSON(w, req, "ERROR", ID, "INVALID_PARAMS: Formato ID incorrecto", http.StatusBadRequest)
    return
  }
  documentoID := bson.ObjectIdHex(ID)

  // Decode del JSON
  // ***************
  decoder := json.NewDecoder(req.Body)
  err := decoder.Decode(&importe)
  if err != nil {
    core.RspMsgJSON(w, req, "ERROR", "JSON", "INVALID_PARAMS: JSON decode erróneo", http.StatusBadRequest)
    return
  }

  // Creo la autorización
  // ********************
  estado, valor, mensaje, httpStat := AutorizacionAlta(documentoID, "Cliente", importe.Importe, req, audit)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }

  // Está todo Ok
  // ************
  //------------------------------------Modificar ######
  s := []string{"Pediste autorización"}
  //--------------------------------------Modificar ######
  core.RspMsgJSON(w, req, "OK", "Pediste autorización", strings.Join(s, ""), http.StatusAccepted)
  return
}

func ViajeAutCosto(w http.ResponseWriter, req *http.Request) {
  //-------------------Modificar ######
  var importe models.ImporteSugerido
  vars := mux.Vars(req)
  ID := vars["docID"]
  audit := "Autorizar Costo"

  // Verifico el formato del campo ID
  // ********************************
  if bson.IsObjectIdHex(ID) != true {
    core.RspMsgJSON(w, req, "ERROR", ID, "INVALID_PARAMS: Formato ID incorrecto", http.StatusBadRequest)
    return
  }
  documentoID := bson.ObjectIdHex(ID)

  // Decode del JSON
  // ***************
  decoder := json.NewDecoder(req.Body)
  err := decoder.Decode(&importe)
  if err != nil {
    core.RspMsgJSON(w, req, "ERROR", "JSON", "INVALID_PARAMS: JSON decode erróneo", http.StatusBadRequest)
    return
  }

  // Creo la autorización
  // ********************
  estado, valor, mensaje, httpStat := AutorizacionAlta(documentoID, "Transportista", importe.Importe, req, audit)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }

  // Está todo Ok
  // ************
  //------------------------------------Modificar ######
  s := []string{"Pediste autorización"}
  //--------------------------------------Modificar ######
  core.RspMsgJSON(w, req, "OK", "Pediste autorización", strings.Join(s, ""), http.StatusAccepted)
  return
}

func AutorizacionAlta(documentoID bson.ObjectId, tipo string, importe float64, req *http.Request, audit string) (string, string, string, int) {
  //-------------------Modificar ###### las 3 variables
	var documento models.Autorizaciones
  coll := config.DB_Autorizacion
  empresaID := context.Get(req, "Empresa_id").(bson.ObjectId)
  usuarioID := context.Get(req, "Usuario_id").(bson.ObjectId)
  usuario := context.Get(req, "Usuario").(string)

  // Busco para obtener los campos del viaje original
  // ************************************************
  //------------------------------------------------------Modificar ######
  estado, valor, mensaje, httpStat, documentoExistente := Viaje_X_ID(documentoID, "Buscar ID", req)
  if httpStat != http.StatusOK {
    return estado, valor, mensaje, httpStat
  }

  // Me fijo si está Editable
  // ************************
  if documentoExistente.Editable == false {
    s := []string{"El viaje está bloqueado"}
    return "ERROR", "Modificar", strings.Join(s, ""), http.StatusNonAuthoritativeInfo
  }

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    return "ERROR", "GetMongoSession", err.Error(), httpStat
  }
  defer session.Close()

  // Intento el alta
  // ***************
  objID := bson.NewObjectId()
  documento.ID = objID
  documento.Empresa_id = empresaID
  documento.Timestamp = time.Now()
	documento.Viaje_id = documentoExistente.ID
	documento.FechaHora = documentoExistente.FechaHora
	documento.Recorrido = documentoExistente.Recorrido
	documento.Kilometraje = documentoExistente.Kilometraje
	documento.Solicitante_id = usuarioID
  documento.Solicitante = usuario
  documento.SolicitanteFecha = time.Now()
  if tipo == "Cliente" {
    documento.TipoSolicitud = "Tarifario Cliente"
  	documento.Titular_id = documentoExistente.Cliente_id
    documento.Titular = documentoExistente.Cliente
  	documento.ImporteTarifario = documentoExistente.TarifaValor
  } else {
    documento.TipoSolicitud = "Tarifario Transportista"
  	documento.Titular_id = documentoExistente.Transportista_id
    documento.Titular = documentoExistente.Transportista
  	documento.ImporteTarifario = documentoExistente.TarifaCosto
  }
	documento.ImporteSugerido = importe
	documento.Autorizante_id = config.FakeID
  documento.Autorizante = ""
  documento.AutorizanteFecha = time.Time{}
	documento.ImporteAutorizado = 0

  collection := session.DB(config.DB_Name).C(coll)
  err = collection.Insert(documento)
  if err != nil {
    s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
    return "ERROR", audit, strings.Join(s, ""), http.StatusInternalServerError
  }

  // Está todo Ok
  // ************
  core.Audit(req, coll, documento.ID, audit, documento)
  return "OK", audit, "Ok", http.StatusOK
}