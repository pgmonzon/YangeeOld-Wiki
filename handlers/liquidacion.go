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

func LiquidacionCrear(w http.ResponseWriter, req *http.Request) {
	var documento models.Liquidacion
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
  estado, valor, mensaje, httpStat, documento := LiquidacionAlta(documento, req, audit)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }

  // Está todo Ok
  // ************
  s := []string{"Agregaste una liquidación para ", documento.Transportista}
  core.RspMsgJSON(w, req, "OK", documento.Transportista, strings.Join(s, ""), http.StatusCreated)
  return
}

func LiquidacionAlta(documentoAlta models.Liquidacion, req *http.Request, audit string) (string, string, string, int, models.Liquidacion) {
	var documento models.Liquidacion
  coll := config.DB_Liquidacion
  empresaID := context.Get(req, "Empresa_id").(bson.ObjectId)
  usuarioID := context.Get(req, "Usuario_id").(bson.ObjectId)
  usuario := context.Get(req, "Usuario").(string)

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
  documento.FechaLiquidacion = time.Now()
  documento.UsuarioLiquidacion_id = usuarioID
  documento.UsuarioLiquidacion = usuario

  collection := session.DB(config.DB_Name).C(coll)
  err = collection.Insert(documento)
  if err != nil {
    s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
    return "ERROR", audit, strings.Join(s, ""), http.StatusInternalServerError, documento
  }

  // Actualizo los viajes
  // ********************
  collectionViajes := session.DB(config.DB_Name).C(config.DB_Viaje)
  for _, item := range documentoAlta.Viajes {
    selector := bson.M{"_id": item.Viaje_id, "empresa_id": empresaID}
    updator := bson.M{
      "$set": bson.M{
        "liquidacion_id": documento.ID,
        "liquidacion": documento.Liquidacion,
        "fechaLiquidacion": documento.FechaLiquidacion,
        "usuarioLiquidacion_id": documento.UsuarioLiquidacion_id,
        "usuarioLiquidacion": documento.UsuarioLiquidacion,
        "timestamp": time.Now(),
      },
    }
    err = collectionViajes.Update(selector, updator)
    if err != nil {
      s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
      return "ERROR", audit, strings.Join(s, ""), http.StatusInternalServerError, documento
    }
  }

  // Está todo Ok
  // ************
  core.Audit(req, coll, documento.ID, audit, documento)
  return "OK", audit, "Ok", http.StatusOK, documento
}

func LiquidacionTraer(w http.ResponseWriter, req *http.Request) {
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
  estado, valor, mensaje, httpStat, documento := Liquidacion_X_ID(documentoID, "Buscar ID", req)
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

func Liquidacion_X_ID(documentoID bson.ObjectId, audit string, req *http.Request) (string, string, string, int, models.Liquidacion) {
  var documento models.Liquidacion
  coll := config.DB_Liquidacion
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

func LiquidacionesTraer(w http.ResponseWriter, req *http.Request) {
  var documentos []models.Liquidacion

  // Busco
  // *****
  estado, valor, mensaje, httpStat, documentos := LiquidacionesBuscar("Buscar", req)
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
func LiquidacionesBuscar(audit string, req *http.Request) (string, string, string, int, []models.Liquidacion) {
  var documentos []models.Liquidacion
  coll := config.DB_Liquidacion
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
  }
  collection := session.DB(config.DB_Name).C(coll)
  collection.Find(selector).Select(bson.M{"empresa_id":0}).Sort("-fecha").All(&documentos)

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
