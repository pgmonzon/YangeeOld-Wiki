package handlers

import (
  "time"
  "encoding/json"
  "net/http"
  "fmt"
  "strings"

  "github.com/pgmonzon/Yangee/models"
  "github.com/pgmonzon/Yangee/core"
  "github.com/pgmonzon/Yangee/config"

  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

func ClienteAPIAlta(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	var clienteAPIAlta models.ClienteAPIAlta

  decoder := json.NewDecoder(req.Body)
  err := decoder.Decode(&clienteAPIAlta)
  if err != nil {
    core.RespErrorJSON(w, req, start, fmt.Errorf("INVALID_PARAMS: JSON decode erróneo"), http.StatusBadRequest)
  } else {
    // campos obligatorios
    if clienteAPIAlta.ClienteAPI == "" || clienteAPIAlta.Firma == "" {
      core.RespErrorJSON(w, req, start, fmt.Errorf("INVALID_PARAMS: Cliente API y firma no pueden estar vacíos"), http.StatusBadRequest)
    } else {
      // me fijo si no existe el clienteAPI
      err, httpStat := ClienteAPIExiste(clienteAPIAlta.ClienteAPI)
      if err != nil {
        core.RespErrorJSON(w, req, start, err, httpStat)
      } else {
        // establezco los campos
        var clienteAPI models.ClienteAPI
      	objID := bson.NewObjectId()
      	clienteAPI.ID = objID
        clienteAPI.ClienteAPI = clienteAPIAlta.ClienteAPI
        clienteAPI.Firma = clienteAPIAlta.Firma
        clienteAPI.Aes = clienteAPIAlta.Aes

        // Genero una nueva sesión Mongo
        session, err, httpStat := core.GetMongoSession()
        if err != nil {
          core.RespErrorJSON(w, req, start, err, httpStat)
        } else {
          defer session.Close()

          // Intento el alta
          collection := session.DB(config.DB_Name).C(config.DB_ClienteAPI)
          err = collection.Insert(clienteAPI)
          if err != nil {
            s := []string{"INTERNAL_SERVER_ERROR:", err.Error()}
            core.RespErrorJSON(w, req, start, fmt.Errorf(strings.Join(s, " ")), http.StatusInternalServerError)
          } else {
            core.RespOkJSON(w, req, start, "Ok", http.StatusCreated)
          }
        }
      }
    }
  }
  return
}

func ClienteAPIExiste(clienteAPIExiste string) (error, int) {
  var clienteAPI models.ClienteAPI
  // Genero una nueva sesión Mongo
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    return err, httpStat
  } else {
    defer session.Close()

    // Me aseguro el índice
    collection := session.DB(config.DB_Name).C(config.DB_ClienteAPI)
    index := mgo.Index{
      Key:        []string{"clienteapi"},
      Unique:     true,
      DropDups:   true,
      Background: true,
      Sparse:     true,
    }
    err := collection.EnsureIndex(index)
    if err != nil {
      s := []string{"INTERNAL_SERVER_ERROR:", err.Error()}
      return fmt.Errorf(strings.Join(s, " ")), http.StatusInternalServerError
    }

    collection.Find(bson.M{"clienteapi": clienteAPIExiste}).One(&clienteAPI)
    if clienteAPI.ID == "" {
      return nil, http.StatusOK
    } else {
      return fmt.Errorf("INVALID_PARAMS: El cliente API ya existe"), http.StatusBadRequest
    }
  }
}

func ClienteAPITraerFirma(clienteAPIHeader string) (string, error, int) {
  var clienteAPI models.ClienteAPI
  // Genero una nueva sesión Mongo
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    return "", err, httpStat
  } else {
    defer session.Close()

    collection := session.DB(config.DB_Name).C(config.DB_ClienteAPI)
    collection.Find(bson.M{"clienteapi": clienteAPIHeader}).One(&clienteAPI)
    if clienteAPI.Firma == "" {
      return "", fmt.Errorf("INVALID_PARAMS: El cliente API no tiene firma"), http.StatusBadRequest
    } else {
      return clienteAPI.Firma, nil, http.StatusOK
    }
  }
}

func ClienteAPITraer(cteAPI string) (models.ClienteAPI, error, int) {
  var clienteAPI models.ClienteAPI
  // Genero una nueva sesión Mongo
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    s := []string{"INTERNAL_SERVER_ERROR:", err.Error()}
    return clienteAPI, fmt.Errorf(strings.Join(s, " ")), httpStat
  } else {
    defer session.Close()

    collection := session.DB(config.DB_Name).C(config.DB_ClienteAPI)
    collection.Find(bson.M{"clienteapi": cteAPI}).One(&clienteAPI)
    return clienteAPI, nil, http.StatusOK
  }
}
