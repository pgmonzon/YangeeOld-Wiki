package handlers

import (
  "encoding/json"
  "net/http"
  "strings"
  "time"

  "github.com/pgmonzon/Yangee/models"
  "github.com/pgmonzon/Yangee/core"
  "github.com/pgmonzon/Yangee/config"

  "gopkg.in/mgo.v2/bson"
  "gopkg.in/mgo.v2/txn"
  "github.com/gorilla/context"
  "github.com/gorilla/mux"
)

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
  documento.Estado = "Pendiente"

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

func AutorizacionesTraer(w http.ResponseWriter, req *http.Request) {
  var documento models.Autorizaciones
  var documentos []models.Autorizaciones

  // Decode del JSON
  // ***************
  decoder := json.NewDecoder(req.Body)
  err := decoder.Decode(&documento)
  if err != nil {
    core.RspMsgJSON(w, req, "ERROR", "JSON", "INVALID_PARAMS: JSON decode erróneo", http.StatusBadRequest)
    return
  }

  // Busco
  // *****
  estado, valor, mensaje, httpStat, documentos := AutorizacionesBuscar(documento, "Buscar", req)
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

func AutorizacionesBuscar(documento models.Autorizaciones, audit string, req *http.Request) (string, string, string, int, []models.Autorizaciones) {
  var documentos []models.Autorizaciones
  coll := config.DB_Autorizacion
  empresaID := context.Get(req, "Empresa_id").(bson.ObjectId)

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    return "ERROR", "GetMongoSession", err.Error(), httpStat, documentos
  }
  defer session.Close()

  // Trato de traerlos
  // *****************
  selector := bson.M{
    "empresa_id": empresaID,
    "estado": documento.Estado,
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

func AutorizacionTraer(w http.ResponseWriter, req *http.Request) {
  vars := mux.Vars(req)
  ID := vars["docID"]

  // Verifico el formato del campo ID
  // ********************************
  if bson.IsObjectIdHex(ID) != true {
    core.RspMsgJSON(w, req, "ERROR", ID, "INVALID_PARAMS: Formato ID incorrecto", http.StatusBadRequest)
    return
  }
  documentoID := bson.ObjectIdHex(ID)

  // Busco
  // *****
  estado, valor, mensaje, httpStat, documento := Autorizacion_X_ID(documentoID, "Buscar ID", req)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }

  // Está todo Ok
  // ************
  respuesta, error := json.Marshal(documento)
  core.FatalErr(error)
  core.RspJSON(w, req, respuesta, http.StatusOK)
  return
}

func Autorizacion_X_ID(documentoID bson.ObjectId, audit string, req *http.Request) (string, string, string, int, models.Autorizaciones) {
  var documento models.Autorizaciones
  coll := config.DB_Autorizacion
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

func AutorizacionRechazar(w http.ResponseWriter, req *http.Request) {
  vars := mux.Vars(req)
  ID := vars["docID"]
  audit := "Rechazar"

  // Verifico el formato del campo ID
  // ********************************
  if bson.IsObjectIdHex(ID) != true {
    core.RspMsgJSON(w, req, "ERROR", ID, "INVALID_PARAMS: Formato ID incorrecto", http.StatusBadRequest)
    return
  }
  documentoID := bson.ObjectIdHex(ID)

  // Busco para obtener los campos faltantes
  // ***************************************
  estado, valor, mensaje, httpStat, autorizacionExistente := Autorizacion_X_ID(documentoID, "Buscar ID", req)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }
  estado, valor, mensaje, httpStat, viajeExistente := Viaje_X_ID(autorizacionExistente.Viaje_id, "Buscar ID", req)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }
  if viajeExistente.Editable == false {
    core.RspMsgJSON(w, req, "ERROR", "Viaje", "INVALID_PARAMS: El viaje está bloqueado", http.StatusBadRequest)
    return
  }

  // Modifico
  // ********
  estado, valor, mensaje, httpStat = AutorizacionRechazo(documentoID, autorizacionExistente, viajeExistente, req, audit)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }

  // Está todo Ok
  // ************
  //------------------------------------Modificar ######
  s := []string{"Rechazaste el pedido"}
  //--------------------------------------Modificar ######
  core.RspMsgJSON(w, req, "OK", "Rechazo", strings.Join(s, ""), http.StatusAccepted)
  return
}

func AutorizacionRechazo(documentoID bson.ObjectId, autorizacionExistente models.Autorizaciones, viajeExistente models.Viaje, req *http.Request, audit string) (string, string, string, int) {
  empresaID := context.Get(req, "Empresa_id").(bson.ObjectId)
  usuarioID := context.Get(req, "Usuario_id").(bson.ObjectId)
  usuario := context.Get(req, "Usuario").(string)

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    return "ERROR", "GetMongoSession", err.Error(), httpStat
  }
  defer session.Close()

  // Intento la modificación en transaccion
  // **************************************
  transaccionID := bson.NewObjectId()
  runner := txn.NewRunner(session.DB(config.DB_Name).C(config.DB_Transaction))
  if autorizacionExistente.TipoSolicitud == "Tarifario Cliente" {
    ops := []txn.Op{{
      C: config.DB_Autorizacion,
      Id: documentoID,
      Assert: bson.M{"empresa_id": empresaID},
      Update: bson.M{
        "$set": bson.M{
          "estado": "Rechazado",
          "autorizante_id": usuarioID,
          "autorizante": usuario,
          "autorizanteFecha": time.Now(),
          "timestamp": time.Now(),
        },
      },
    }, {
      C: config.DB_Viaje,
      Id: autorizacionExistente.Viaje_id,
      Assert: bson.M{"empresa_id": empresaID, "editable": true},
      Update: bson.M{
        "$set": bson.M{
          "valorViaje": viajeExistente.TarifaValor,
          "autValorViaje_id": config.FakeID,
          "autValor": "",
          "autValorViajeFecha": time.Time{},
          "timestamp": time.Now(),
        },
      },
    }}
    err = runner.Run(ops, transaccionID, nil)
    if err != nil {
      s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
      return "ERROR", audit, strings.Join(s, ""), http.StatusInternalServerError
    }
  } else {
    ops := []txn.Op{{
      C: config.DB_Autorizacion,
      Id: documentoID,
      Assert: bson.M{"empresa_id": empresaID},
      Update: bson.M{
        "$set": bson.M{
          "estado": "Rechazado",
          "autorizante_id": usuarioID,
          "autorizante": usuario,
          "autorizanteFecha": time.Now(),
          "timestamp": time.Now(),
        },
      },
    }, {
      C: config.DB_Viaje,
      Id: autorizacionExistente.Viaje_id,
      Assert: bson.M{"empresa_id": empresaID, "editable": true},
      Update: bson.M{
        "$set": bson.M{
          "costoViaje": viajeExistente.TarifaCosto,
          "autCostoViaje_id": config.FakeID,
          "autCosto": "",
          "autCostoViajeFecha": time.Time{},
          "timestamp": time.Now(),
        },
      },
    }}
    err = runner.Run(ops, transaccionID, nil)
    if err != nil {
      s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
      return "ERROR", audit, strings.Join(s, ""), http.StatusInternalServerError
    }
  }

  // Está todo Ok
  // ************
  core.Audit(req, config.DB_Autorizacion, documentoID, audit, "")
  return "OK", audit, "Ok", http.StatusOK
}

func AutorizacionAutorizar(w http.ResponseWriter, req *http.Request) {
  var importe models.ImporteSugerido
  vars := mux.Vars(req)
  ID := vars["docID"]
  audit := "Autorizar"

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

  // Busco para obtener los campos faltantes
  // ***************************************
  estado, valor, mensaje, httpStat, autorizacionExistente := Autorizacion_X_ID(documentoID, "Buscar ID", req)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }
  estado, valor, mensaje, httpStat, viajeExistente := Viaje_X_ID(autorizacionExistente.Viaje_id, "Buscar ID", req)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }
  if viajeExistente.Editable == false {
    core.RspMsgJSON(w, req, "ERROR", "Viaje", "INVALID_PARAMS: El viaje está bloqueado", http.StatusBadRequest)
    return
  }

  // Modifico
  // ********
  estado, valor, mensaje, httpStat = AutorizacionAutorizo(documentoID, importe, autorizacionExistente, viajeExistente, req, audit)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }

  // Está todo Ok
  // ************
  //------------------------------------Modificar ######
  s := []string{"Autorizaste el pedido"}
  //--------------------------------------Modificar ######
  core.RspMsgJSON(w, req, "OK", "Autorizo", strings.Join(s, ""), http.StatusAccepted)
  return
}

func AutorizacionAutorizo(documentoID bson.ObjectId, importe models.ImporteSugerido, autorizacionExistente models.Autorizaciones, viajeExistente models.Viaje, req *http.Request, audit string) (string, string, string, int) {
  empresaID := context.Get(req, "Empresa_id").(bson.ObjectId)
  usuarioID := context.Get(req, "Usuario_id").(bson.ObjectId)
  usuario := context.Get(req, "Usuario").(string)

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    return "ERROR", "GetMongoSession", err.Error(), httpStat
  }
  defer session.Close()

  // Intento la modificación en transaccion
  // **************************************
  transaccionID := bson.NewObjectId()
  runner := txn.NewRunner(session.DB(config.DB_Name).C(config.DB_Transaction))
  if autorizacionExistente.TipoSolicitud == "Tarifario Cliente" {
    ops := []txn.Op{{
      C: config.DB_Autorizacion,
      Id: documentoID,
      Assert: bson.M{"empresa_id": empresaID},
      Update: bson.M{
        "$set": bson.M{
          "estado": "Autorizado",
          "autorizante_id": usuarioID,
          "autorizante": usuario,
          "autorizanteFecha": time.Now(),
          "importeAutorizado": importe.Importe,
          "timestamp": time.Now(),
        },
      },
    }, {
      C: config.DB_Viaje,
      Id: autorizacionExistente.Viaje_id,
      Assert: bson.M{"empresa_id": empresaID, "editable": true},
      Update: bson.M{
        "$set": bson.M{
          "valorViaje": importe.Importe,
          "autValorViaje_id": usuarioID,
          "autValor": usuario,
          "autValorViajeFecha": time.Now(),
          "timestamp": time.Now(),
        },
      },
    }}
    err = runner.Run(ops, transaccionID, nil)
    if err != nil {
      s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
      return "ERROR", audit, strings.Join(s, ""), http.StatusInternalServerError
    }
  } else {
    ops := []txn.Op{{
      C: config.DB_Autorizacion,
      Id: documentoID,
      Assert: bson.M{"empresa_id": empresaID},
      Update: bson.M{
        "$set": bson.M{
          "estado": "Autorizado",
          "autorizante_id": usuarioID,
          "autorizante": usuario,
          "autorizanteFecha": time.Now(),
          "importeAutorizado": importe.Importe,
          "timestamp": time.Now(),
        },
      },
    }, {
      C: config.DB_Viaje,
      Id: autorizacionExistente.Viaje_id,
      Assert: bson.M{"empresa_id": empresaID, "editable": true},
      Update: bson.M{
        "$set": bson.M{
          "costoViaje": importe.Importe,
          "autCostoViaje_id": usuarioID,
          "autCosto": usuario,
          "autCostoViajeFecha": time.Now(),
          "timestamp": time.Now(),
        },
      },
    }}
    err = runner.Run(ops, transaccionID, nil)
    if err != nil {
      s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
      return "ERROR", audit, strings.Join(s, ""), http.StatusInternalServerError
    }
  }

  // Está todo Ok
  // ************
  core.Audit(req, config.DB_Autorizacion, documentoID, audit, "")
  return "OK", audit, "Ok", http.StatusOK
}
