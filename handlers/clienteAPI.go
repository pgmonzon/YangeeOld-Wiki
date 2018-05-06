package handlers

import (
  "encoding/json"
  "net/http"
  "fmt"
  "strings"
  "time"

  "github.com/pgmonzon/Yangee/models"
  "github.com/pgmonzon/Yangee/core"
  "github.com/pgmonzon/Yangee/config"

  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "github.com/gorilla/context"
)


func ClienteAPICrear(w http.ResponseWriter, req *http.Request) {
	var clienteAPI models.ClienteAPI

  // Decode del JSON
  // ***************
  decoder := json.NewDecoder(req.Body)
  err := decoder.Decode(&clienteAPI)
  if err != nil {
    core.RspMsgJSON(w, req, "ERROR", "JSON", "INVALID_PARAMS: JSON decode erróneo", http.StatusBadRequest)
    return
  }

  // Verifico los campos obligatorios
  // ********************************
  if clienteAPI.ClienteAPI == "" || clienteAPI.Firma == "" || clienteAPI.Aes == "" {
    core.RspMsgJSON(w, req, "ERROR", "ClienteAPI", "INVALID_PARAMS: ClienteAPI, firma y aes no pueden estar vacíos", http.StatusBadRequest)
    return
  }

  // Me fijo si ya Existe
  // ********************
  err, httpStat := ClienteAPIExiste(clienteAPI.ClienteAPI)
  if err != nil {
    core.RspMsgJSON(w, req, "ERROR", clienteAPI.ClienteAPI, err.Error(), httpStat)
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
  clienteAPI.ID = objID
  clienteAPI.Timestamp = time.Now()
  clienteAPI.Borrado = false
  context.Set(req, "Coleccion", config.DB_ClienteAPI)
  collection := session.DB(config.DB_Name).C(config.DB_ClienteAPI)
  err = collection.Insert(clienteAPI)
  if err != nil {
    s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
    core.RspMsgJSON(w, req, "ERROR", clienteAPI.ClienteAPI, strings.Join(s, ""), http.StatusInternalServerError)
    return
  }

  // Establezco las variables
  // ************************
  context.Set(req, "TipoOper", "#Novedad#")
  context.Set(req, "Coleccion", config.DB_ClienteAPI)
  context.Set(req, "Objeto_id", clienteAPI.ID)
  context.Set(req, "Audit", clienteAPI)

  // Está todo Ok
  // ************
  core.RspMsgJSON(w, req, "OK", clienteAPI.ClienteAPI, "Ok", http.StatusCreated)
  return
}

func ClienteAPIExiste(clienteAPIExiste string) (error, int) {
  var clienteAPI models.ClienteAPI

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
  collection := session.DB(config.DB_Name).C(config.DB_ClienteAPI)
  index := mgo.Index{
    Key:        []string{"clienteapi"},
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
  collection.Find(bson.M{"clienteapi": clienteAPIExiste}).One(&clienteAPI)
  // No existe
  if clienteAPI.ID == "" {
    return nil, http.StatusOK
  }
  // Existe borrado
  if clienteAPI.Borrado == true {
    s := []string{"INVALID_PARAMS: El clienteAPI ", clienteAPIExiste," ya existe borrado"}
    return fmt.Errorf(strings.Join(s, "")), http.StatusBadRequest
  }
  // Existe inactivo
  if clienteAPI.Activo == false {
    s := []string{"INVALID_PARAMS: El clienteAPI ", clienteAPIExiste," ya existe inactivo"}
    return fmt.Errorf(strings.Join(s, "")), http.StatusBadRequest
  }
  // Existe
  s := []string{"INVALID_PARAMS: El clienteAPI ", clienteAPIExiste," ya existe"}
  return fmt.Errorf(strings.Join(s, "")), http.StatusBadRequest
}

func ClienteAPI_X_clienteAPI(cteAPI string) (models.ClienteAPI, error, int) {
  var clienteAPI models.ClienteAPI

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    s := []string{err.Error()}
    return clienteAPI, fmt.Errorf(strings.Join(s, "")), httpStat
  }
  defer session.Close()

  // Trato de obtener el ClienteAPI
  // ******************************
  collection := session.DB(config.DB_Name).C(config.DB_ClienteAPI)
  collection.Find(bson.M{"clienteapi": cteAPI}).One(&clienteAPI)
  // Si no existe devuelvo error
  if clienteAPI.ID == "" {
    s := []string{"INVALID_PARAMS: El clienteAPI no existe"}
    return clienteAPI, fmt.Errorf(strings.Join(s, "")), http.StatusBadRequest
  }

  // Existe
  // ******
  return clienteAPI, nil, http.StatusOK
}
