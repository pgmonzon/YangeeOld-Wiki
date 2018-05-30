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
