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

func SbrIngresoSucursalCrear(w http.ResponseWriter, req *http.Request) {
	var documento models.SbrIngresoSucursal
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
  estado, valor, mensaje, httpStat, documento := SbrIngresoSucursalAlta(documento, req, audit)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }

  // Está todo Ok
  // ************
  s := []string{"Agregaste un ingreso a la sucursal ", documento.Sucursal}
  core.RspMsgJSON(w, req, "OK", documento.Sucursal, strings.Join(s, ""), http.StatusCreated)
  return
}

func SbrIngresoSucursalAlta(documentoAlta models.SbrIngresoSucursal, req *http.Request, audit string) (string, string, string, int, models.SbrIngresoSucursal) {
	var documento models.SbrIngresoSucursal
  var docStock models.SbrStock
  coll := config.DB_SbrIngresoSucursal
  empresaID := context.Get(req, "Empresa_id").(bson.ObjectId)

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
  documento.Fecha = time.Now()

  collection := session.DB(config.DB_Name).C(coll)
  err = collection.Insert(documento)
  if err != nil {
    s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
    return "ERROR", audit, strings.Join(s, ""), http.StatusInternalServerError, documento
  }

  // Actualizo el stock
  // ******************
  for _, item := range documentoAlta.Detalle {
    // me fijo si existe
    docStock.ID = ""
    collectionStock := session.DB(config.DB_Name).C(config.DB_SbrStock)
    collectionStock.Find(bson.M{"empresa_id": empresaID, "sucursal_id": documentoAlta.Sucursal_id, "articulo_id": item.SbrArticulo_id}).One(&docStock)

    if docStock.ID == "" {
      var docActualizar models.SbrStock
      objID2 := bson.NewObjectId()
      docActualizar.ID = objID2
      docActualizar.Empresa_id = empresaID
      docActualizar.SbrSucursal_id = documentoAlta.Sucursal_id
      docActualizar.SbrSucursal = documentoAlta.Sucursal
      docActualizar.SbrArticulo_id = item.SbrArticulo_id
      docActualizar.SbrArticulo = item.SbrArticulo
      docActualizar.Cantidad = item.Cantidad

      collectionActualizar := session.DB(config.DB_Name).C(config.DB_SbrStock)
      err = collectionActualizar.Insert(docActualizar)
      if err != nil {
        s := []string{"INTERNAL_SERVER_ERRORdd: ", err.Error()}
        return "ERROR", audit, strings.Join(s, ""), http.StatusInternalServerError, documento
      }
    } else {
      suma := docStock.Cantidad + item.Cantidad
      selector := bson.M{"empresa_id": empresaID, "_id": docStock.ID}
      updator := bson.M{"$set": bson.M{"cantidad": suma}}

      collectionActualizar := session.DB(config.DB_Name).C(config.DB_SbrStock)
      err = collectionActualizar.Update(selector, updator)
      if err != nil {
        s := []string{"INTERNAL_SERVER_ERRORhh: ", err.Error()}
        return "ERROR", audit, strings.Join(s, ""), http.StatusInternalServerError, documento
      }
    }
  }

  // Está todo Ok
  // ************
  core.Audit(req, coll, documento.ID, audit, documento)
  return "OK", audit, "Ok", http.StatusOK, documento
}

func SbrIngresosSucursalesTraer(w http.ResponseWriter, req *http.Request) {
  var documento models.SbrIngresoSucursal
  var documentos []models.SbrIngresoSucursal
  vars := mux.Vars(req)
  orden := vars["orden"]
  limite := vars["limite"]
  sucrusal := vars["sucursal"]

  // Verifico el formato del campo ID
  // ********************************
  if bson.IsObjectIdHex(sucrusal) != true {
    core.RspMsgJSON(w, req, "ERROR", sucrusal, "INVALID_PARAMS: Formato IDSucursal incorrecto", http.StatusBadRequest)
    return
  }
  sucursalID := bson.ObjectIdHex(sucrusal)

  // Verifico el formato del campo limite
  // ************************************
  limiteInt, err := strconv.Atoi(limite)
  if err != nil {
    core.RspMsgJSON(w, req, "ERROR", "Límite debe ser numérico", err.Error(), http.StatusBadRequest)
    return
  }

  // Busco
  // *****
  estado, valor, mensaje, httpStat, documentos := SbrIngresosSucursalesBuscar(documento, orden, limiteInt, sucursalID, false, "Buscar", req)
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

func SbrIngresosSucursalesBuscar(documento models.SbrIngresoSucursal, orden string, limiteInt int, sucursalID bson.ObjectId, borrados bool, audit string, req *http.Request) (string, string, string, int, []models.SbrIngresoSucursal) {
  var documentos []models.SbrIngresoSucursal
  coll := config.DB_SbrIngresoSucursal
  empresaID := context.Get(req, "Empresa_id").(bson.ObjectId)

  // Verifico que el campo orden sea Unique
  // **************************************
  if orden != "fecha" && orden != "-fecha" {
    s := []string{"No puedo ordenar por ", orden}
    return "ERROR", "Buscar", strings.Join(s, ""), http.StatusNonAuthoritativeInfo, documentos
  }

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
    "sucursal_id": sucursalID,
  }
  collection := session.DB(config.DB_Name).C(coll)
  collection.Find(selector).Select(bson.M{"empresa_id":0}).Sort(orden).Limit(limiteInt).All(&documentos)

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

func SbrIngresoSucursalTraer(w http.ResponseWriter, req *http.Request) {
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
  //----------------------------------------------Modificar ######
  estado, valor, mensaje, httpStat, documento := SbrIngresoSucursal_X_ID(documentoID, "Buscar ID", req)
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

func SbrIngresoSucursal_X_ID(documentoID bson.ObjectId, audit string, req *http.Request) (string, string, string, int, models.SbrIngresoSucursal) {
  var documento models.SbrIngresoSucursal
  coll := config.DB_SbrIngresoSucursal
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
