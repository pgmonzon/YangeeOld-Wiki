package handlers

import (
  "encoding/json"
  "net/http"
  "fmt"
  "strings"

  "github.com/pgmonzon/Yangee/models"
  "github.com/pgmonzon/Yangee/core"
  "github.com/pgmonzon/Yangee/config"

  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "github.com/gorilla/context"
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
  err, httpStat := FilosofoExiste(filosofo.Filosofo)
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
  s := []string{"Agregó el filósofo ", filosofo.Filosofo}
  context.Set(req, "Novedad", strings.Join(s, ""))
  context.Set(req, "Objeto_id", filosofo.ID)
  context.Set(req, "Audit", filosofo)

  // Está todo Ok
  // ************
  core.RspMsgJSON(w, req, "OK", filosofo.Filosofo, "Ok", http.StatusCreated)
  return
}

func FilosofoExiste(filosofoExiste string) (error, int) {
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
