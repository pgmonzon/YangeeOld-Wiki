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

  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "github.com/gorilla/context"
  "github.com/gorilla/mux"
)

func SbrVentasCrear(w http.ResponseWriter, req *http.Request) {
  var documentoAlta models.SbrVentasCrear
  audit := "Crear"

  // Decode del JSON
  // ***************
  decoder := json.NewDecoder(req.Body)
  err := decoder.Decode(&documentoAlta)
  if err != nil {
    core.RspMsgJSON(w, req, "ERROR", "JSON", "INVALID_PARAMS: JSON decode erróneo", http.StatusBadRequest)
    return
  }

  // Doy de alta
  // ***********
  estado, valor, mensaje, httpStat, documento := SbrVentasAlta(documentoAlta, req, audit)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }

  // Está todo Ok
  // ************
  s := []string{"Agregaste una caja a la sucursal ", documento.Sucursal}
  core.RspMsgJSON(w, req, "OK", documento.Sucursal, strings.Join(s, ""), http.StatusCreated)
  return
}

func SbrVentasAlta(documentoAlta models.SbrVentasCrear, req *http.Request, audit string) (string, string, string, int, models.SbrVentas) {
  var documento models.SbrVentas
  coll := config.DB_SbrVentas
  empresaID := context.Get(req, "Empresa_id").(bson.ObjectId)

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    return "ERROR", "GetMongoSession", err.Error(), httpStat, documento
  }
  defer session.Close()

  // Me aseguro el índice
  // ********************
  indice := []string{"fecha", "sucursal_id"}
  collection := session.DB(config.DB_Name).C(coll)
  index := mgo.Index{
    Key:        indice,
    Unique:     true,
    DropDups:   false,
    Background: true,
    Sparse:     true,
  }
  err = collection.EnsureIndex(index)
  if err != nil {
    s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
    return "ERROR", "EnsureIndex", strings.Join(s, ""), http.StatusInternalServerError, documento
  }

  // Intento el alta
  // ***************
  t := time.Now()
  objID := bson.NewObjectId()
  documento.ID = objID
  documento.Empresa_id = empresaID
  documento.Timestamp = t
  documento.Fecha = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
  documento.Estado = "Abierto"
  documento.Sucursal_id = documentoAlta.Sucursal_id
  documento.Sucursal = documentoAlta.Sucursal
  documento.Tarjeta = 0
  documento.Efectivo = 0
  documento.Total = 0
  documento.RendidoA_id = config.FakeID
  documento.RendidoA = ""

  collection = session.DB(config.DB_Name).C(coll)
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

func SbrVentasTraer(w http.ResponseWriter, req *http.Request) {
  var documento models.SbrVentas
  var documentos []models.SbrVentas
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
  estado, valor, mensaje, httpStat, documentos := SbrVentasBuscar(documento, orden, limiteInt, sucursalID, false, "Buscar", req)
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

func SbrVentasBuscar(documento models.SbrVentas, orden string, limiteInt int, sucursalID bson.ObjectId, borrados bool, audit string, req *http.Request) (string, string, string, int, []models.SbrVentas) {
  var documentos []models.SbrVentas
  coll := config.DB_SbrVentas
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

func SbrVentasCerrarCaja(w http.ResponseWriter, req *http.Request) {
  var documento models.SbrVentasCerrar
  vars := mux.Vars(req)
  ID := vars["docID"]
  audit := "Cerrar Caja"

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
  estado, valor, mensaje, httpStat := SbrVentasModificarCaja(documentoID, documento, req, audit)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }

  // Está todo Ok
  // ************
  s := []string{"Rendido a ", documento.RendidoA}
  core.RspMsgJSON(w, req, "OK", documento.RendidoA, strings.Join(s, ""), http.StatusAccepted)
  return
}

func SbrVentasModificarCaja(documentoID bson.ObjectId, documentoModi models.SbrVentasCerrar, req *http.Request, audit string) (string, string, string, int) {
  coll := config.DB_SbrVentas
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
      "estado": "Cerrado",
      "rendidoA_id": documentoModi.RendidoA_id,
      "rendidoA": documentoModi.RendidoA,
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

func SbrVentasDetalleCrear(w http.ResponseWriter, req *http.Request) {
  var documento models.SbrVentasDetalle
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
  estado, valor, mensaje, httpStat, documento := SbrVentasDetalleAlta(documento, req, audit)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }

  // Está todo Ok
  // ************
  s := []string{"Agregaste un movimiento a la caja ", documento.SbrArticulo}
  core.RspMsgJSON(w, req, "OK", documento.SbrArticulo, strings.Join(s, ""), http.StatusCreated)
  return
}

func SbrVentasDetalleAlta(documentoAlta models.SbrVentasDetalle, req *http.Request, audit string) (string, string, string, int, models.SbrVentasDetalle) {
  var docStock models.SbrStock
  var documento models.SbrVentasDetalle
  coll := config.DB_SbrVentasDetalle
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

  collection := session.DB(config.DB_Name).C(coll)
  err = collection.Insert(documento)
  if err != nil {
    s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
    return "ERROR", audit, strings.Join(s, ""), http.StatusInternalServerError, documento
  }

  // Descuento SbrStock
  // ******************
  docStock.ID = ""
  collectionStock := session.DB(config.DB_Name).C(config.DB_SbrStock)
  collectionStock.Find(bson.M{"empresa_id": empresaID, "sucursal_id": documento.Sucursal_id, "articulo_id": documento.SbrArticulo_id}).One(&docStock)

  if docStock.ID == "" {
    var docActualizar models.SbrStock
    objID2 := bson.NewObjectId()
    docActualizar.ID = objID2
    docActualizar.Empresa_id = empresaID
    docActualizar.SbrSucursal_id = documento.Sucursal_id
    docActualizar.SbrSucursal = documento.Sucursal
    docActualizar.SbrArticulo_id = documento.SbrArticulo_id
    docActualizar.SbrArticulo = documento.SbrArticulo
    docActualizar.Cantidad = -1

    collectionActualizar := session.DB(config.DB_Name).C(config.DB_SbrStock)
    err = collectionActualizar.Insert(docActualizar)
    if err != nil {
      s := []string{"INTERNAL_SERVER_ERRORdd: ", err.Error()}
      return "ERROR", audit, strings.Join(s, ""), http.StatusInternalServerError, documento
    }
  } else {
    suma := docStock.Cantidad -1
    selector := bson.M{"empresa_id": empresaID, "_id": docStock.ID}
    updator := bson.M{"$set": bson.M{"cantidad": suma}}

    collectionActualizar := session.DB(config.DB_Name).C(config.DB_SbrStock)
    err = collectionActualizar.Update(selector, updator)
    if err != nil {
      s := []string{"INTERNAL_SERVER_ERRORhh: ", err.Error()}
      return "ERROR", audit, strings.Join(s, ""), http.StatusInternalServerError, documento
    }
  }

  // Totalizo
  // ********
  SbrVentasTotalizar(documentoAlta.SbrVentas_id, req)

  // Está todo Ok
  // ************
  core.Audit(req, coll, documento.ID, audit, documento)
  return "OK", audit, "Ok", http.StatusOK, documento
}

func SbrVentasDetalleBorrar(w http.ResponseWriter, req *http.Request) {
  var documento models.SbrVentasDetalle
  vars := mux.Vars(req)
  ID := vars["docID"]
  audit := "Borrar"

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

  // Doy de baja
  // ***********
  estado, valor, mensaje, httpStat, _ := SbrVentasDetalleBaja(documentoID, documento, req, audit)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }

  // Está todo Ok
  // ************
  s := []string{"Borraste un movimiento a la caja"}
  core.RspMsgJSON(w, req, "OK", "Movimiento borrado", strings.Join(s, ""), http.StatusAccepted)
  return
}

func SbrVentasDetalleBaja(documentoID bson.ObjectId, documento models.SbrVentasDetalle,req *http.Request, audit string) (string, string, string, int, models.SbrVentasDetalle) {
  var docStock models.SbrStock
  var documentoBaja models.SbrVentasDetalle
  coll := config.DB_SbrVentasDetalle
  empresaID := context.Get(req, "Empresa_id").(bson.ObjectId)

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    return "ERROR", "GetMongoSession", err.Error(), httpStat, documentoBaja
  }
  defer session.Close()

  // Intento la baja
  // ***************
  collection := session.DB(config.DB_Name).C(coll)
  err = collection.Remove(bson.M{"_id": documentoID})
  if err != nil {
    s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
    return "ERROR", audit, strings.Join(s, ""), http.StatusInternalServerError, documentoBaja
  }

  // Descuento SbrStock
  // ******************
  docStock.ID = ""
  collectionStock := session.DB(config.DB_Name).C(config.DB_SbrStock)
  collectionStock.Find(bson.M{"empresa_id": empresaID, "sucursal_id": documento.Sucursal_id, "articulo_id": documento.SbrArticulo_id}).One(&docStock)

  if docStock.ID == "" {
    var docActualizar models.SbrStock
    objID2 := bson.NewObjectId()
    docActualizar.ID = objID2
    docActualizar.Empresa_id = empresaID
    docActualizar.SbrSucursal_id = documento.Sucursal_id
    docActualizar.SbrSucursal = documento.Sucursal
    docActualizar.SbrArticulo_id = documento.SbrArticulo_id
    docActualizar.SbrArticulo = documento.SbrArticulo
    docActualizar.Cantidad = +1

    collectionActualizar := session.DB(config.DB_Name).C(config.DB_SbrStock)
    err = collectionActualizar.Insert(docActualizar)
    if err != nil {
      s := []string{"INTERNAL_SERVER_ERRORdd: ", err.Error()}
      return "ERROR", audit, strings.Join(s, ""), http.StatusInternalServerError, documento
    }
  } else {
    suma := docStock.Cantidad +1
    selector := bson.M{"empresa_id": empresaID, "_id": docStock.ID}
    updator := bson.M{"$set": bson.M{"cantidad": suma}}

    collectionActualizar := session.DB(config.DB_Name).C(config.DB_SbrStock)
    err = collectionActualizar.Update(selector, updator)
    if err != nil {
      s := []string{"INTERNAL_SERVER_ERRORhh: ", err.Error()}
      return "ERROR", audit, strings.Join(s, ""), http.StatusInternalServerError, documento
    }
  }

  // Totalizo
  // ********
  SbrVentasTotalizar(documento.SbrVentas_id, req)

  // Está todo Ok
  // ************
  core.Audit(req, coll, documentoBaja.ID, audit, documentoBaja)
  return "OK", audit, "Ok", http.StatusOK, documentoBaja
}

func SbrVentasDetalleTraer(w http.ResponseWriter, req *http.Request) {
  var documentos []models.SbrVentasDetalle
  vars := mux.Vars(req)
  venta := vars["docID"]

  // Verifico el formato del campo ID
  // ********************************
  if bson.IsObjectIdHex(venta) != true {
    core.RspMsgJSON(w, req, "ERROR", venta, "INVALID_PARAMS: Formato IDVenta incorrecto", http.StatusBadRequest)
    return
  }
  ventaID := bson.ObjectIdHex(venta)

  // Busco
  // *****
  estado, valor, mensaje, httpStat, documentos := SbrVentasDetalleBuscar(ventaID, false, "Buscar", req)
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

func SbrVentasDetalleBuscar(ventaID bson.ObjectId, borrados bool, audit string, req *http.Request) (string, string, string, int, []models.SbrVentasDetalle) {
  var documentos []models.SbrVentasDetalle
  coll := config.DB_SbrVentasDetalle

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
    "sbrVentas_id": ventaID,
  }
  collection := session.DB(config.DB_Name).C(coll)
  collection.Find(selector).Select(bson.M{"empresa_id":0}).Sort("-timestamp").All(&documentos)

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

func SbrVentasGastoCrear(w http.ResponseWriter, req *http.Request) {
  var documento models.SbrVentasGastos
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
  estado, valor, mensaje, httpStat, documento := SbrVentasGastoAlta(documento, req, audit)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }

  // Está todo Ok
  // ************
  s := []string{"Agregaste un gasto a la caja ", documento.CuentaGasto}
  core.RspMsgJSON(w, req, "OK", documento.CuentaGasto, strings.Join(s, ""), http.StatusCreated)
  return
}

func SbrVentasGastoAlta(documentoAlta models.SbrVentasGastos, req *http.Request, audit string) (string, string, string, int, models.SbrVentasGastos) {
  var documento models.SbrVentasGastos
  coll := config.DB_SbrVentasGastos
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

  collection := session.DB(config.DB_Name).C(coll)
  err = collection.Insert(documento)
  if err != nil {
    s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
    return "ERROR", audit, strings.Join(s, ""), http.StatusInternalServerError, documento
  }

  // Totalizo
  // ********
  SbrVentasTotalizar(documentoAlta.SbrVentas_id, req)

  // Está todo Ok
  // ************
  core.Audit(req, coll, documento.ID, audit, documento)
  return "OK", audit, "Ok", http.StatusOK, documento
}

func SbrVentasGastoBorrar(w http.ResponseWriter, req *http.Request) {
  var documento models.SbrVentasGastos
  vars := mux.Vars(req)
  ID := vars["docID"]
  audit := "Borrar"

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

  // Doy de baja
  // ***********
  estado, valor, mensaje, httpStat, _ := SbrVentasGastoBaja(documentoID, documento, req, audit)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }

  // Está todo Ok
  // ************
  s := []string{"Borraste un gasto de la caja"}
  core.RspMsgJSON(w, req, "OK", "Movimiento borrado", strings.Join(s, ""), http.StatusAccepted)
  return
}

func SbrVentasGastoBaja(documentoID bson.ObjectId, documento models.SbrVentasGastos,req *http.Request, audit string) (string, string, string, int, models.SbrVentasGastos) {
  var documentoBaja models.SbrVentasGastos
  coll := config.DB_SbrVentasGastos

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    return "ERROR", "GetMongoSession", err.Error(), httpStat, documentoBaja
  }
  defer session.Close()

  // Intento la baja
  // ***************
  collection := session.DB(config.DB_Name).C(coll)
  err = collection.Remove(bson.M{"_id": documentoID})
  if err != nil {
    s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
    return "ERROR", audit, strings.Join(s, ""), http.StatusInternalServerError, documentoBaja
  }

  // Totalizo
  // ********
  SbrVentasTotalizar(documento.SbrVentas_id, req)

  // Está todo Ok
  // ************
  core.Audit(req, coll, documentoBaja.ID, audit, documentoBaja)
  return "OK", audit, "Ok", http.StatusOK, documentoBaja
}

func SbrVentasGastoTraer(w http.ResponseWriter, req *http.Request) {
  var documentos []models.SbrVentasGastos
  vars := mux.Vars(req)
  venta := vars["docID"]

  // Verifico el formato del campo ID
  // ********************************
  if bson.IsObjectIdHex(venta) != true {
    core.RspMsgJSON(w, req, "ERROR", venta, "INVALID_PARAMS: Formato IDVenta incorrecto", http.StatusBadRequest)
    return
  }
  ventaID := bson.ObjectIdHex(venta)

  // Busco
  // *****
  estado, valor, mensaje, httpStat, documentos := SbrVentasGastoBuscar(ventaID, false, "Buscar", req)
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

func SbrVentasGastoBuscar(ventaID bson.ObjectId, borrados bool, audit string, req *http.Request) (string, string, string, int, []models.SbrVentasGastos) {
  var documentos []models.SbrVentasGastos
  coll := config.DB_SbrVentasGastos

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
    "sbrVentas_id": ventaID,
  }
  collection := session.DB(config.DB_Name).C(coll)
  collection.Find(selector).Select(bson.M{"empresa_id":0}).Sort("-timestamp").All(&documentos)

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

func SbrVentasTotalizar(ventaID bson.ObjectId, req *http.Request) {
  var tarjeta float64
  var efectivo float64
  var total float64
  var gastos float64
  var importeRendir float64
  coll := config.DB_SbrVentas
  empresaID := context.Get(req, "Empresa_id").(bson.ObjectId)

  // Busco las ventas
  // ****************
  _, _, _, httpStat, docVentas := SbrVentasDetalleBuscar(ventaID, false, "Buscar", req)
  if httpStat == http.StatusOK {
    // Recorro y sumo
    // **************
    for _, item := range docVentas {
      if item.FormaPago == "Tarjeta" {
        tarjeta = tarjeta + item.Cobrado
      } else {
        efectivo = efectivo + item.Cobrado
      }
    }
  }

  // Busco los gastos
  // ****************
  _, _, _, httpStat, docGastos := SbrVentasGastoBuscar(ventaID, false, "Buscar", req)
  if httpStat == http.StatusOK {
    // Recorro y sumo
    // **************
    for _, item := range docGastos {
      gastos = gastos + item.Importe
    }
  }

  total = efectivo + tarjeta
  importeRendir = total - gastos

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    return
  }
  defer session.Close()

  // Intento la modificación
  // ***********************
  collection := session.DB(config.DB_Name).C(coll)
  selector := bson.M{"_id": ventaID, "empresa_id": empresaID}
  updator := bson.M{
    "$set": bson.M{
      "tarjeta": tarjeta,
      "efectivo": efectivo,
      "total": total,
      "gastos": gastos,
      "importeRendir": importeRendir,
      "timestamp": time.Now(),
    },
  }
  err = collection.Update(selector, updator)
  if err != nil {
    return
  }

  // Está todo Ok
  // ************
  return
}
