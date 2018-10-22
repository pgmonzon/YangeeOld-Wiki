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

func SbrRemitoSucursalCrear(w http.ResponseWriter, req *http.Request) {
	var documento models.SbrRemitoSucursal
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
  estado, valor, mensaje, httpStat, documento := SbrRemitoSucursalAlta(documento, req, audit)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }

  // Está todo Ok
  // ************
  s := []string{"Agregaste un remito a la sucursal ", documento.ASucursal}
  core.RspMsgJSON(w, req, "OK", documento.DeSucursal, strings.Join(s, ""), http.StatusCreated)
  return
}

func SbrRemitoSucursalAlta(documentoAlta models.SbrRemitoSucursal, req *http.Request, audit string) (string, string, string, int, models.SbrRemitoSucursal) {
	var documento models.SbrRemitoSucursal
  coll := config.DB_SbrRemitoSucursal
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
  documento.Estado = "Enviado"
  documento.Recibio_id = config.FakeID

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

func SbrRemitosDeSucursalTraer(w http.ResponseWriter, req *http.Request) {
  var documento models.SbrRemitoSucursal
  var documentos []models.SbrRemitoSucursal
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
  estado, valor, mensaje, httpStat, documentos := SbrRemitosDeSucursalBuscar(documento, orden, limiteInt, sucursalID, false, "Buscar", req)
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

func SbrRemitosDeSucursalBuscar(documento models.SbrRemitoSucursal, orden string, limiteInt int, sucursalID bson.ObjectId, borrados bool, audit string, req *http.Request) (string, string, string, int, []models.SbrRemitoSucursal) {
  var documentos []models.SbrRemitoSucursal
  coll := config.DB_SbrRemitoSucursal
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
    "deSucursal_id": sucursalID,
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

func SbrRemitosASucursalTraer(w http.ResponseWriter, req *http.Request) {
  var documento models.SbrRemitoSucursal
  var documentos []models.SbrRemitoSucursal
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
  estado, valor, mensaje, httpStat, documentos := SbrRemitosASucursalBuscar(documento, orden, limiteInt, sucursalID, false, "Buscar", req)
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

func SbrRemitosASucursalBuscar(documento models.SbrRemitoSucursal, orden string, limiteInt int, sucursalID bson.ObjectId, borrados bool, audit string, req *http.Request) (string, string, string, int, []models.SbrRemitoSucursal) {
  var documentos []models.SbrRemitoSucursal
  coll := config.DB_SbrRemitoSucursal
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
    "aSucursal_id": sucursalID,
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

func SbrRemitoSucursalTraer(w http.ResponseWriter, req *http.Request) {
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
  estado, valor, mensaje, httpStat, documento := SbrRemitoSucursal_X_ID(documentoID, "Buscar ID", req)
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

func SbrRemitoSucursal_X_ID(documentoID bson.ObjectId, audit string, req *http.Request) (string, string, string, int, models.SbrRemitoSucursal) {
  var documento models.SbrRemitoSucursal
  coll := config.DB_SbrRemitoSucursal
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

func SbrRemitoSucursalCancelar(w http.ResponseWriter, req *http.Request) {
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

  // Busco para obtener los campos faltantes
  // ***************************************
  estado, valor, mensaje, httpStat, documentoExistente := SbrRemitoSucursal_X_ID(documentoID, "Buscar ID", req)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }
  if documentoExistente.Estado != "Enviado" {
    core.RspMsgJSON(w, req, "ERROR", audit, "Debe estar en estado Enviado para cancelarlo", http.StatusBadRequest)
    return
  }
  documentoExistente.Estado = "Cancelado"

  // Modifico
  // ********
  estado, valor, mensaje, httpStat = SbrRemitoSucursalModificarEstado(documentoID, documentoExistente, req, audit)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }

  // Está todo Ok
  // ************
  s := []string{"Cancelaste envío a ", documentoExistente.ASucursal}
  core.RspMsgJSON(w, req, "OK", documentoExistente.ASucursal, strings.Join(s, ""), http.StatusAccepted)
  return
}

func SbrRemitoSucursalAceptar(w http.ResponseWriter, req *http.Request) {
  var documento models.SbrRemitoSucursalAceptarRechazar
  vars := mux.Vars(req)
  ID := vars["docID"]
  audit := "Aceptar"

  // Decode del JSON
  // ***************
  decoder := json.NewDecoder(req.Body)
  err := decoder.Decode(&documento)
  if err != nil {
    core.RspMsgJSON(w, req, "ERROR", "JSON", "INVALID_PARAMS: JSON decode erróneo", http.StatusBadRequest)
    return
  }

  // Verifico el formato del campo ID
  // ********************************
  if bson.IsObjectIdHex(ID) != true {
    core.RspMsgJSON(w, req, "ERROR", ID, "INVALID_PARAMS: Formato ID incorrecto", http.StatusBadRequest)
    return
  }
  documentoID := bson.ObjectIdHex(ID)

  // Busco para obtener los campos faltantes
  // ***************************************
  estado, valor, mensaje, httpStat, documentoExistente := SbrRemitoSucursal_X_ID(documentoID, "Buscar ID", req)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }
  if documentoExistente.Estado != "Enviado" {
    core.RspMsgJSON(w, req, "ERROR", audit, "Debe estar en estado Enviado para aceptarlo", http.StatusBadRequest)
    return
  }
  documentoExistente.Estado = "Aceptado"

  // Modifico
  // ********
  estado, valor, mensaje, httpStat = SbrRemitoSucursalAceptarRechazar(documentoID, documentoExistente, documento, req, audit)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }

  // Está todo Ok
  // ************
  s := []string{"Aceptaste envío de ", documentoExistente.DeSucursal}
  core.RspMsgJSON(w, req, "OK", documentoExistente.DeSucursal, strings.Join(s, ""), http.StatusAccepted)
  return
}

func SbrRemitoSucursalRechazar(w http.ResponseWriter, req *http.Request) {
  var documento models.SbrRemitoSucursalAceptarRechazar
  vars := mux.Vars(req)
  ID := vars["docID"]
  audit := "Rechazar"

  // Decode del JSON
  // ***************
  decoder := json.NewDecoder(req.Body)
  err := decoder.Decode(&documento)
  if err != nil {
    core.RspMsgJSON(w, req, "ERROR", "JSON", "INVALID_PARAMS: JSON decode erróneo", http.StatusBadRequest)
    return
  }

  // Verifico el formato del campo ID
  // ********************************
  if bson.IsObjectIdHex(ID) != true {
    core.RspMsgJSON(w, req, "ERROR", ID, "INVALID_PARAMS: Formato ID incorrecto", http.StatusBadRequest)
    return
  }
  documentoID := bson.ObjectIdHex(ID)

  // Busco para obtener los campos faltantes
  // ***************************************
  estado, valor, mensaje, httpStat, documentoExistente := SbrRemitoSucursal_X_ID(documentoID, "Buscar ID", req)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }
  if documentoExistente.Estado != "Enviado" {
    core.RspMsgJSON(w, req, "ERROR", audit, "Debe estar en estado Enviado para rechazarlo", http.StatusBadRequest)
    return
  }
  documentoExistente.Estado = "Rechazado"

  // Modifico
  // ********
  estado, valor, mensaje, httpStat = SbrRemitoSucursalAceptarRechazar(documentoID, documentoExistente, documento, req, audit)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }

  // Está todo Ok
  // ************
  s := []string{"Rechazaste envío de ", documentoExistente.DeSucursal}
  core.RspMsgJSON(w, req, "OK", documentoExistente.DeSucursal, strings.Join(s, ""), http.StatusAccepted)
  return
}

func SbrRemitoSucursalAceptarRechazar(documentoID bson.ObjectId, documentoModi models.SbrRemitoSucursal, documento models.SbrRemitoSucursalAceptarRechazar, req *http.Request, audit string) (string, string, string, int) {
  var docStock models.SbrStock
  coll := config.DB_SbrRemitoSucursal
  empresaID := context.Get(req, "Empresa_id").(bson.ObjectId)

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    return "ERROR", "GetMongoSession", err.Error(), httpStat
  }
  defer session.Close()

  // Intento la modificación
  // ***********************
  collection := session.DB(config.DB_Name).C(coll)
  selector := bson.M{"_id": documentoID, "empresa_id": empresaID}
  updator := bson.M{
    "$set": bson.M{
      "estado": documentoModi.Estado,
      "recibio_id": documento.Recibio_id,
      "recibio": documento.Recibio,
      "fechaRecepcion": time.Now(),
      "timestamp": time.Now(),
    },
  }
  err = collection.Update(selector, updator)
  if err != nil {
    s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
    return "ERROR", audit, strings.Join(s, ""), http.StatusInternalServerError
  }

  // Si es Aceptado hago los movimientos de stock
  // ********************************************
  if documentoModi.Estado == "Aceptado" {
    for _, item := range documentoModi.Detalle {
      // Ingreso primero en la sucursal destino
      docStock.ID = ""
      collectionStock := session.DB(config.DB_Name).C(config.DB_SbrStock)
      collectionStock.Find(bson.M{"empresa_id": empresaID, "sucursal_id": documentoModi.ASucursal_id, "articulo_id": item.SbrArticulo_id}).One(&docStock)

      if docStock.ID == "" {
        var docActualizar models.SbrStock
        objID2 := bson.NewObjectId()
        docActualizar.ID = objID2
        docActualizar.Empresa_id = empresaID
        docActualizar.SbrSucursal_id = documentoModi.ASucursal_id
        docActualizar.SbrSucursal = documentoModi.ASucursal
        docActualizar.SbrArticulo_id = item.SbrArticulo_id
        docActualizar.SbrArticulo = item.SbrArticulo
        docActualizar.Cantidad = item.Cantidad

        collectionActualizar := session.DB(config.DB_Name).C(config.DB_SbrStock)
        err = collectionActualizar.Insert(docActualizar)
        if err != nil {
          s := []string{"INTERNAL_SERVER_ERRORdd: ", err.Error()}
          return "ERROR", audit, strings.Join(s, ""), http.StatusInternalServerError
        }
      } else {
        suma := docStock.Cantidad + item.Cantidad
        selector := bson.M{"empresa_id": empresaID, "_id": docStock.ID}
        updator := bson.M{"$set": bson.M{"cantidad": suma}}

        collectionActualizar := session.DB(config.DB_Name).C(config.DB_SbrStock)
        err = collectionActualizar.Update(selector, updator)
        if err != nil {
          s := []string{"INTERNAL_SERVER_ERRORhh: ", err.Error()}
          return "ERROR", audit, strings.Join(s, ""), http.StatusInternalServerError
        }
      }

      // Saco de la sucursal origen
      docStock.ID = ""
      collectionStock = session.DB(config.DB_Name).C(config.DB_SbrStock)
      collectionStock.Find(bson.M{"empresa_id": empresaID, "sucursal_id": documentoModi.DeSucursal_id, "articulo_id": item.SbrArticulo_id}).One(&docStock)

      if docStock.ID == "" {
        var docActualizar models.SbrStock
        objID2 := bson.NewObjectId()
        docActualizar.ID = objID2
        docActualizar.Empresa_id = empresaID
        docActualizar.SbrSucursal_id = documentoModi.DeSucursal_id
        docActualizar.SbrSucursal = documentoModi.DeSucursal
        docActualizar.SbrArticulo_id = item.SbrArticulo_id
        docActualizar.SbrArticulo = item.SbrArticulo
        docActualizar.Cantidad = item.Cantidad * -1

        collectionActualizar := session.DB(config.DB_Name).C(config.DB_SbrStock)
        err = collectionActualizar.Insert(docActualizar)
        if err != nil {
          s := []string{"INTERNAL_SERVER_ERRORdd: ", err.Error()}
          return "ERROR", audit, strings.Join(s, ""), http.StatusInternalServerError
        }
      } else {
        suma := docStock.Cantidad - item.Cantidad
        selector := bson.M{"empresa_id": empresaID, "_id": docStock.ID}
        updator := bson.M{"$set": bson.M{"cantidad": suma}}

        collectionActualizar := session.DB(config.DB_Name).C(config.DB_SbrStock)
        err = collectionActualizar.Update(selector, updator)
        if err != nil {
          s := []string{"INTERNAL_SERVER_ERRORhh: ", err.Error()}
          return "ERROR", audit, strings.Join(s, ""), http.StatusInternalServerError
        }
      }
    }
  }

  // Está todo Ok
  // ************
  core.Audit(req, coll, documentoID, audit, documentoModi)
  return "OK", audit, "Ok", http.StatusOK
}

func SbrRemitoSucursalModificarEstado(documentoID bson.ObjectId, documentoModi models.SbrRemitoSucursal, req *http.Request, audit string) (string, string, string, int) {
  var docStock models.SbrStock
  coll := config.DB_SbrRemitoSucursal
  empresaID := context.Get(req, "Empresa_id").(bson.ObjectId)

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    return "ERROR", "GetMongoSession", err.Error(), httpStat
  }
  defer session.Close()

  // Intento la modificación
  // ***********************
  collection := session.DB(config.DB_Name).C(coll)
  selector := bson.M{"_id": documentoID, "empresa_id": empresaID}
  updator := bson.M{
    "$set": bson.M{
      "estado": documentoModi.Estado,
      "timestamp": time.Now(),
    },
  }
  err = collection.Update(selector, updator)
  if err != nil {
    s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
    return "ERROR", audit, strings.Join(s, ""), http.StatusInternalServerError
  }

  // Si es Aceptado hago los movimientos de stock
  // ********************************************
  if documentoModi.Estado == "Aceptado" {
    for _, item := range documentoModi.Detalle {
      // Ingreso primero en la sucursal destino
      docStock.ID = ""
      collectionStock := session.DB(config.DB_Name).C(config.DB_SbrStock)
      collectionStock.Find(bson.M{"empresa_id": empresaID, "sucursal_id": documentoModi.ASucursal_id, "articulo_id": item.SbrArticulo_id}).One(&docStock)

      if docStock.ID == "" {
        var docActualizar models.SbrStock
        objID2 := bson.NewObjectId()
        docActualizar.ID = objID2
        docActualizar.Empresa_id = empresaID
        docActualizar.SbrSucursal_id = documentoModi.ASucursal_id
        docActualizar.SbrSucursal = documentoModi.ASucursal
        docActualizar.SbrArticulo_id = item.SbrArticulo_id
        docActualizar.SbrArticulo = item.SbrArticulo
        docActualizar.Cantidad = item.Cantidad

        collectionActualizar := session.DB(config.DB_Name).C(config.DB_SbrStock)
        err = collectionActualizar.Insert(docActualizar)
        if err != nil {
          s := []string{"INTERNAL_SERVER_ERRORdd: ", err.Error()}
          return "ERROR", audit, strings.Join(s, ""), http.StatusInternalServerError
        }
      } else {
        suma := docStock.Cantidad + item.Cantidad
        selector := bson.M{"empresa_id": empresaID, "_id": docStock.ID}
        updator := bson.M{"$set": bson.M{"cantidad": suma}}

        collectionActualizar := session.DB(config.DB_Name).C(config.DB_SbrStock)
        err = collectionActualizar.Update(selector, updator)
        if err != nil {
          s := []string{"INTERNAL_SERVER_ERRORhh: ", err.Error()}
          return "ERROR", audit, strings.Join(s, ""), http.StatusInternalServerError
        }
      }

      // Saco de la sucursal origen
      docStock.ID = ""
      collectionStock = session.DB(config.DB_Name).C(config.DB_SbrStock)
      collectionStock.Find(bson.M{"empresa_id": empresaID, "sucursal_id": documentoModi.DeSucursal_id, "articulo_id": item.SbrArticulo_id}).One(&docStock)

      if docStock.ID == "" {
        var docActualizar models.SbrStock
        objID2 := bson.NewObjectId()
        docActualizar.ID = objID2
        docActualizar.Empresa_id = empresaID
        docActualizar.SbrSucursal_id = documentoModi.DeSucursal_id
        docActualizar.SbrSucursal = documentoModi.DeSucursal
        docActualizar.SbrArticulo_id = item.SbrArticulo_id
        docActualizar.SbrArticulo = item.SbrArticulo
        docActualizar.Cantidad = item.Cantidad * -1

        collectionActualizar := session.DB(config.DB_Name).C(config.DB_SbrStock)
        err = collectionActualizar.Insert(docActualizar)
        if err != nil {
          s := []string{"INTERNAL_SERVER_ERRORdd: ", err.Error()}
          return "ERROR", audit, strings.Join(s, ""), http.StatusInternalServerError
        }
      } else {
        suma := docStock.Cantidad - item.Cantidad
        selector := bson.M{"empresa_id": empresaID, "_id": docStock.ID}
        updator := bson.M{"$set": bson.M{"cantidad": suma}}

        collectionActualizar := session.DB(config.DB_Name).C(config.DB_SbrStock)
        err = collectionActualizar.Update(selector, updator)
        if err != nil {
          s := []string{"INTERNAL_SERVER_ERRORhh: ", err.Error()}
          return "ERROR", audit, strings.Join(s, ""), http.StatusInternalServerError
        }
      }
    }
  }

  // Está todo Ok
  // ************
  core.Audit(req, coll, documentoID, audit, documentoModi)
  return "OK", audit, "Ok", http.StatusOK
}
