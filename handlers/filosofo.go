package handlers

import (
  "encoding/json"
  "net/http"
  "fmt"
  "strings"
  "strconv"
  "time"

  "github.com/pgmonzon/Yangee/models"
  "github.com/pgmonzon/Yangee/core"
  "github.com/pgmonzon/Yangee/config"

  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "github.com/gorilla/context"
  "github.com/gorilla/mux"
)

func FilosofoCrear(w http.ResponseWriter, req *http.Request) {
	var filosofo models.Filosofo

  // Decode del JSON
  // ***************
  decoder := json.NewDecoder(req.Body)
  err := decoder.Decode(&filosofo)
  if err != nil {
    core.RspMsgJSON(w, req, "ERROR", "JSON", "INVALID_PARAMS: JSON decode erróneo", http.StatusBadRequest)
    return
  }

  // Verifico los campos obligatorios
  // ********************************
  if filosofo.Filosofo == "" {
    core.RspMsgJSON(w, req, "ERROR", "Filósofo", "INVALID_PARAMS: Filósofo no puede estar vacío", http.StatusBadRequest)
    return
  }

  // Me fijo si ya Existe
  // ********************
  err, httpStat := FilosofoExiste(filosofo.Filosofo, "")
  if err != nil {
    core.RspMsgJSON(w, req, "ERROR", filosofo.Filosofo, err.Error(), httpStat)
    return
  }

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    core.RspMsgJSON(w, req, "ERROR", "MongoSession", err.Error(), httpStat)
    return
  }
  defer session.Close()

  // Intento el alta
  // ***************
  objID := bson.NewObjectId()
  filosofo.ID = objID
  filosofo.Timestamp = time.Now()
  filosofo.Borrado = false
  collection := session.DB(config.DB_Name).C(config.DB_Filosofo)
  err = collection.Insert(filosofo)
  if err != nil {
    s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
    core.RspMsgJSON(w, req, "ERROR", filosofo.Filosofo, strings.Join(s, ""), http.StatusInternalServerError)
    return
  }

  // Establezco las variables
  // ************************
  context.Set(req, "TipoOper", "#Novedad#")
  context.Set(req, "Coleccion", config.DB_Filosofo)
  context.Set(req, "Objeto_id", filosofo.ID)
  context.Set(req, "Audit", filosofo)

  // Está todo Ok
  // ************
  s := []string{"Agregó el filósofo ", filosofo.Filosofo}
  core.RspMsgJSON(w, req, "OK", filosofo.Filosofo, strings.Join(s, ""), http.StatusCreated)
  return
}

func FilosofoModificar(w http.ResponseWriter, req *http.Request) {
	var filosofo models.Filosofo

  // Verifico el formato del campo ID
  // ********************************
  vars := mux.Vars(req)
  if bson.IsObjectIdHex(vars["filosofoID"]) != true {
    core.RspMsgJSON(w, req, "ERROR", "filosofoID", "INVALID_PARAMS: Formato ID incorrecto", http.StatusBadRequest)
    return
  }
  filosofoID := bson.ObjectIdHex(vars["filosofoID"])

  // Decode del JSON
  // ***************
  decoder := json.NewDecoder(req.Body)
  err := decoder.Decode(&filosofo)
  if err != nil {
    core.RspMsgJSON(w, req, "ERROR", "JSON", "INVALID_PARAMS: JSON decode erróneo", http.StatusBadRequest)
    return
  }

  // Verifico los campos obligatorios
  // ********************************
  if filosofo.Filosofo == "" {
    core.RspMsgJSON(w, req, "ERROR", "Filósofo", "INVALID_PARAMS: Filósofo no puede estar vacío", http.StatusBadRequest)
    return
  }

  // Me fijo si ya Existe
  // ********************
  err, httpStat := FilosofoExiste(filosofo.Filosofo, vars["filosofoID"])
  if err != nil {
    core.RspMsgJSON(w, req, "ERROR", filosofo.Filosofo, err.Error(), httpStat)
    return
  }

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    core.RspMsgJSON(w, req, "ERROR", "MongoSession", err.Error(), httpStat)
    return
  }
  defer session.Close()

  // Intento la modificación
  // ***********************
  filosofo.ID = filosofoID
  collection := session.DB(config.DB_Name).C(config.DB_Filosofo)
  selector := bson.M{"_id": filosofo.ID}
  updator := bson.M{
    "$set": bson.M{
      "filosofo": filosofo.Filosofo,
      "doctrina": filosofo.Doctrina,
      "biografia": filosofo.Biografia,
      "activo": filosofo.Activo,
      "timestamp": time.Now(),
    },
  }
  err = collection.Update(selector, updator)
  if err != nil {
    s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
    core.RspMsgJSON(w, req, "ERROR", filosofo.Filosofo, strings.Join(s, ""), http.StatusInternalServerError)
    return
  }

  // Establezco las variables
  // ************************
  context.Set(req, "TipoOper", "#Novedad#")
  context.Set(req, "Coleccion", config.DB_Filosofo)
  context.Set(req, "Objeto_id", filosofo.ID)
  context.Set(req, "Audit", filosofo)

  // Está todo Ok
  // ************
  s := []string{"Modificó el filósofo ", filosofo.Filosofo}
  core.RspMsgJSON(w, req, "OK", filosofo.Filosofo, strings.Join(s, ""), http.StatusAccepted)
  return
}

func FilosofoBorrar(w http.ResponseWriter, req *http.Request) {

  // Verifico el formato del campo ID
  // ********************************
  vars := mux.Vars(req)
  if bson.IsObjectIdHex(vars["filosofoID"]) != true {
    core.RspMsgJSON(w, req, "ERROR", "filosofoID", "INVALID_PARAMS: Formato ID incorrecto", http.StatusBadRequest)
    return
  }
  filosofoID := bson.ObjectIdHex(vars["filosofoID"])

  // Traigo los datos
  // ****************
  filosofo, err, httpStat := Filosofo_X_ID(filosofoID)
  if err != nil {
    core.RspMsgJSON(w, req, "ERROR", vars["filosofoID"], err.Error(), httpStat)
    return
  }

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    core.RspMsgJSON(w, req, "ERROR", "MongoSession", err.Error(), httpStat)
    return
  }
  defer session.Close()

  // Intento borrarlo
  // ****************
  collection := session.DB(config.DB_Name).C(config.DB_Filosofo)
  selector := bson.M{"_id": filosofo.ID}
  updator := bson.M{
    "$set": bson.M{
      "borrado": true,
      "timestamp": time.Now(),
    },
  }
  err = collection.Update(selector, updator)
  if err != nil {
    s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
    core.RspMsgJSON(w, req, "ERROR", filosofo.Filosofo, strings.Join(s, ""), http.StatusInternalServerError)
    return
  }
  filosofo.Borrado = true

  // Establezco las variables
  // ************************
  context.Set(req, "TipoOper", "#Novedad#")
  context.Set(req, "Coleccion", config.DB_Filosofo)
  context.Set(req, "Objeto_id", filosofo.ID)
  context.Set(req, "Audit", filosofo)

  // Está todo Ok
  // ************
  s := []string{"Borró el filósofo ", filosofo.Filosofo}
  core.RspMsgJSON(w, req, "OK", filosofo.Filosofo, strings.Join(s, ""), http.StatusAccepted)
  return
}

func FilosofoTraer(w http.ResponseWriter, req *http.Request) {

  // Verifico el formato del campo ID
  // ********************************
  vars := mux.Vars(req)
  if bson.IsObjectIdHex(vars["filosofoID"]) != true {
    core.RspMsgJSON(w, req, "ERROR", "filosofoID", "INVALID_PARAMS: Formato ID incorrecto", http.StatusBadRequest)
    return
  }
  filosofoID := bson.ObjectIdHex(vars["filosofoID"])

  // Traigo los datos
  // ****************
  filosofo, err, httpStat := Filosofo_X_ID(filosofoID)
  if err != nil {
    core.RspMsgJSON(w, req, "ERROR", vars["filosofoID"], err.Error(), httpStat)
    return
  }

  // Establezco las variables
  // ************************
  context.Set(req, "TipoOper", "#Oper#")
  context.Set(req, "Coleccion", config.DB_Filosofo)
  context.Set(req, "Objeto_id", filosofo.ID)
  context.Set(req, "Audit", "")

  // Está todo Ok
  // ************
  s := []string{"Trajo el filósofo ", filosofo.Filosofo}
  context.Set(req, "Novedad", strings.Join(s, ""))
  respuesta, error := json.Marshal(filosofo)
  core.FatalErr(error)
  core.RspJSON(w, req, respuesta, http.StatusOK)
  return
}

func FilosofosTraer(w http.ResponseWriter, req *http.Request) {
  var filosofo models.Filosofo
  var filosofos []models.Filosofo
  vars := mux.Vars(req)
  orden := vars["orden"]
  limite := vars["limite"]

  // Verifico el formato del campo limite
  // ************************************
  limiteInt, err := strconv.Atoi(limite)
  if err != nil {
    core.RspMsgJSON(w, req, "ERROR", "Límite debe ser numérico", err.Error(), http.StatusBadRequest)
    return
  }

  // Verifico que el campo orden sea Unique
  // **************************************
  if orden != "filosofo" && orden != "-filosofo" {
    s := []string{"El campo ", orden, " no está soportado para ordenar"}
    core.RspMsgJSON(w, req, "ERROR", "Orden", strings.Join(s, ""), http.StatusBadRequest)
    return
  }

  // Decode del JSON
  // ***************
  decoder := json.NewDecoder(req.Body)
  err = decoder.Decode(&filosofo)
  if err != nil {
    core.RspMsgJSON(w, req, "ERROR", "JSON", "INVALID_PARAMS: JSON decode erróneo", http.StatusBadRequest)
    return
  }

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    core.RspMsgJSON(w, req, "ERROR", "MongoSession", err.Error(), httpStat)
    return
  }
  defer session.Close()

  // Trato de traerlos
  // *****************
  selector := bson.M{
    "filosofo": bson.M{"$regex": bson.RegEx{filosofo.Filosofo, "i"}},
    "doctrina": bson.M{"$regex": bson.RegEx{filosofo.Doctrina, "i"}},
    "biografia": bson.M{"$regex": bson.RegEx{filosofo.Biografia, "i"}},
    "activo": filosofo.Activo,
    "borrado": false,
  }
  collection := session.DB(config.DB_Name).C(config.DB_Filosofo)
  collection.Find(selector).Sort(orden).Limit(limiteInt).All(&filosofos)

  // Si el resultado es vacío devuelvo ERROR
  // ***************************************
  if filosofos == nil {
    core.RspMsgJSON(w, req, "ERROR", "Filosofos", "INVALID_PARAMS: No hay filosofos para tu búsqueda", http.StatusBadRequest)
    return
  }

  // Establezco las variables
  // ************************
  context.Set(req, "TipoOper", "#Oper#")
  context.Set(req, "Coleccion", config.DB_Filosofo)
  context.Set(req, "Objeto_id", context.Get(req, "CicloDeVida_id").(bson.ObjectId))
  context.Set(req, "Audit", "")

  // Está todo Ok
  // ************
  s := []string{"Trajo filósofos /orden: ", orden, " /limite: ", limite}
  context.Set(req, "Novedad", strings.Join(s, ""))
  respuesta, error := json.Marshal(filosofos)
  core.FatalErr(error)
  core.RspJSON(w, req, respuesta, http.StatusOK)
  return
}

func FilosofosTraerSiguiente(w http.ResponseWriter, req *http.Request) {
  var filosofo models.Filosofo
  var filosofos []models.Filosofo
  vars := mux.Vars(req)
  orden := vars["orden"]
  limite := vars["limite"]
  ultimoOrden := vars["ultimo_campo_orden"]

  // Verifico el formato del campo limite
  // ************************************
  limiteInt, err := strconv.Atoi(limite)
  if err != nil {
    core.RspMsgJSON(w, req, "ERROR", "Límite debe ser numérico", err.Error(), http.StatusBadRequest)
    return
  }

  // Verifico que el campo orden sea Unique
  // **************************************
  if orden != "filosofo" && orden != "-filosofo" {
    s := []string{"El campo ", orden, " no está soportado para ordenar"}
    core.RspMsgJSON(w, req, "ERROR", "Orden", strings.Join(s, ""), http.StatusBadRequest)
    return
  }

  // Decode del JSON
  // ***************
  decoder := json.NewDecoder(req.Body)
  err = decoder.Decode(&filosofo)
  if err != nil {
    core.RspMsgJSON(w, req, "ERROR", "JSON", "INVALID_PARAMS: JSON decode erróneo", http.StatusBadRequest)
    return
  }

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    core.RspMsgJSON(w, req, "ERROR", "MongoSession", err.Error(), httpStat)
    return
  }
  defer session.Close()

  // Trato de traerlos
  // *****************
  var selector bson.M
  if strings.HasPrefix(orden, "-") {
    selector = bson.M{
      "filosofo": bson.M{"$lt": ultimoOrden, "$regex": bson.RegEx{filosofo.Filosofo, "i"}},
      "doctrina": bson.M{"$regex": bson.RegEx{filosofo.Doctrina, "i"}},
      "biografia": bson.M{"$regex": bson.RegEx{filosofo.Biografia, "i"}},
      "activo": filosofo.Activo,
      "borrado": false,
    }
  } else {
    selector = bson.M{
      "filosofo": bson.M{"$gt": ultimoOrden, "$regex": bson.RegEx{filosofo.Filosofo, "i"}},
      "doctrina": bson.M{"$regex": bson.RegEx{filosofo.Doctrina, "i"}},
      "biografia": bson.M{"$regex": bson.RegEx{filosofo.Biografia, "i"}},
      "activo": filosofo.Activo,
      "borrado": false,
    }
  }
  collection := session.DB(config.DB_Name).C(config.DB_Filosofo)
  collection.Find(selector).Sort(orden).Limit(limiteInt).All(&filosofos)

  // Si el resultado es vacío devuelvo ERROR
  // ***************************************
  if filosofos == nil {
    core.RspMsgJSON(w, req, "ERROR", "Filosofos", "INVALID_PARAMS: No hay filosofos para tu búsqueda", http.StatusBadRequest)
    return
  }

  // Establezco las variables
  // ************************
  context.Set(req, "TipoOper", "#Oper#")
  context.Set(req, "Coleccion", config.DB_Filosofo)
  context.Set(req, "Objeto_id", context.Get(req, "CicloDeVida_id").(bson.ObjectId))
  context.Set(req, "Audit", "")

  // Está todo Ok
  // ************
  s := []string{"Trajo filósofos siguiente /orden: ", orden, " /limite: ", limite}
  context.Set(req, "Novedad", strings.Join(s, ""))
  respuesta, error := json.Marshal(filosofos)
  core.FatalErr(error)
  core.RspJSON(w, req, respuesta, http.StatusOK)
  return
}

func FilosofosTraerAnterior(w http.ResponseWriter, req *http.Request) {
  var filosofo models.Filosofo
  var filosofos []models.Filosofo
  vars := mux.Vars(req)
  orden := vars["orden"]
  limite := vars["limite"]
  primerOrden := vars["primer_campo_orden"]

  // Verifico el formato del campo limite
  // ************************************
  limiteInt, err := strconv.Atoi(limite)
  if err != nil {
    core.RspMsgJSON(w, req, "ERROR", "Límite debe ser numérico", err.Error(), http.StatusBadRequest)
    return
  }

  // Verifico que el campo orden sea Unique
  // **************************************
  if orden != "filosofo" && orden != "-filosofo" {
    s := []string{"El campo ", orden, " no está soportado para ordenar"}
    core.RspMsgJSON(w, req, "ERROR", "Orden", strings.Join(s, ""), http.StatusBadRequest)
    return
  }

  // Decode del JSON
  // ***************
  decoder := json.NewDecoder(req.Body)
  err = decoder.Decode(&filosofo)
  if err != nil {
    core.RspMsgJSON(w, req, "ERROR", "JSON", "INVALID_PARAMS: JSON decode erróneo", http.StatusBadRequest)
    return
  }

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    core.RspMsgJSON(w, req, "ERROR", "MongoSession", err.Error(), httpStat)
    return
  }
  defer session.Close()

  // Trato de traerlos
  // *****************
  var selector bson.M
  if strings.HasPrefix(orden, "-") {
    selector = bson.M{
      "filosofo": bson.M{"$gt": primerOrden, "$regex": bson.RegEx{filosofo.Filosofo, "i"}},
      "doctrina": bson.M{"$regex": bson.RegEx{filosofo.Doctrina, "i"}},
      "biografia": bson.M{"$regex": bson.RegEx{filosofo.Biografia, "i"}},
      "activo": filosofo.Activo,
      "borrado": false,
    }
  } else {
    selector = bson.M{
      "filosofo": bson.M{"$lt": primerOrden, "$regex": bson.RegEx{filosofo.Filosofo, "i"}},
      "doctrina": bson.M{"$regex": bson.RegEx{filosofo.Doctrina, "i"}},
      "biografia": bson.M{"$regex": bson.RegEx{filosofo.Biografia, "i"}},
      "activo": filosofo.Activo,
      "borrado": false,
    }
  }
  collection := session.DB(config.DB_Name).C(config.DB_Filosofo)
  collection.Find(selector).Sort(orden).Limit(limiteInt).All(&filosofos)

  // Si el resultado es vacío devuelvo ERROR
  // ***************************************
  if filosofos == nil {
    core.RspMsgJSON(w, req, "ERROR", "Filosofos", "INVALID_PARAMS: No hay filosofos para tu búsqueda", http.StatusBadRequest)
    return
  }

  // Establezco las variables
  // ************************
  context.Set(req, "TipoOper", "#Oper#")
  context.Set(req, "Coleccion", config.DB_Filosofo)
  context.Set(req, "Objeto_id", context.Get(req, "CicloDeVida_id").(bson.ObjectId))
  context.Set(req, "Audit", "")

  // Está todo Ok
  // ************
  s := []string{"Trajo filósofos anterior /orden: ", orden, " /limite: ", limite}
  context.Set(req, "Novedad", strings.Join(s, ""))
  respuesta, error := json.Marshal(filosofos)
  core.FatalErr(error)
  core.RspJSON(w, req, respuesta, http.StatusOK)
  return
}

func FilosofoExiste(filosofoExiste string, filosofoString string) (error, int) {
  var filosofo models.Filosofo

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, _ := core.GetMongoSession()
  if err != nil {
    s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
    return fmt.Errorf(strings.Join(s, "")), http.StatusInternalServerError
  }
  defer session.Close()

  // Me aseguro el índice
  // ********************
  collection := session.DB(config.DB_Name).C(config.DB_Filosofo)
  index := mgo.Index{
    Key:        []string{"filosofo"},
    Unique:     true,
    DropDups:   false,
    Background: true,
    Sparse:     true,
  }
  err = collection.EnsureIndex(index)
  if err != nil {
    s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
    return fmt.Errorf(strings.Join(s, "")), http.StatusInternalServerError
  }

  // Verifico si Existe
  // ******************
  collection.Find(bson.M{"filosofo": filosofoExiste}).One(&filosofo)
  // Me fijo si coincide con el recibido
  if bson.IsObjectIdHex(filosofoString) == true {
    if filosofo.ID == bson.ObjectIdHex(filosofoString) {
      return nil, http.StatusOK
    }
  }
  // No existe
  if filosofo.ID == "" {
    return nil, http.StatusOK
  }
  // Existe borrado
  if filosofo.Borrado == true {
    s := []string{"INVALID_PARAMS: El filósofo ", filosofoExiste," ya existe borrado"}
    return fmt.Errorf(strings.Join(s, "")), http.StatusBadRequest
  }
  // Existe inactivo
  if filosofo.Activo == false {
    s := []string{"INVALID_PARAMS: El filósofo ", filosofoExiste," ya existe inactivo"}
    return fmt.Errorf(strings.Join(s, "")), http.StatusBadRequest
  }
  // Existe
  s := []string{"INVALID_PARAMS: El filósofo ", filosofoExiste," ya existe"}
  return fmt.Errorf(strings.Join(s, "")), http.StatusBadRequest
}

func Filosofo_X_ID(filosofoID bson.ObjectId) (models.Filosofo, error, int) {
  var filosofo models.Filosofo

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, _ := core.GetMongoSession()
  if err != nil {
    s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
    return filosofo, fmt.Errorf(strings.Join(s, "")), http.StatusInternalServerError
  }
  defer session.Close()

  // Trato de traerlo
  // ****************
  collection := session.DB(config.DB_Name).C(config.DB_Filosofo)
  collection.Find(bson.M{"_id": filosofoID}).One(&filosofo)
  // Si no existe devuelvo error
  if filosofo.ID == "" {
    s := []string{"INVALID_PARAMS: El filósofo no existe"}
    return filosofo, fmt.Errorf(strings.Join(s, "")), http.StatusBadRequest
  }

  // Existe
  // ******
  return filosofo, nil, http.StatusOK
}
