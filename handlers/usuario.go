package handlers

import (
  "time"
  "encoding/json"
  "net/http"
  "fmt"

  "github.com/pgmonzon/Yangee/models"
  "github.com/pgmonzon/Yangee/core"
  "github.com/pgmonzon/Yangee/config"

  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

func UsuarioRegistrar(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	var usuarioRegistrar models.UsuarioRegistrar

  decoder := json.NewDecoder(req.Body)
  err := decoder.Decode(&usuarioRegistrar)
  if err != nil {
    core.RespErrorJSON(w, req, start, fmt.Errorf("INVALID_PARAMS: JSON decode erróneo"), http.StatusBadRequest)
  } else {
    // campos obligatorios
    if usuarioRegistrar.Usuario == "" || usuarioRegistrar.Clave == "" || usuarioRegistrar.Mail == "" {
      core.RespErrorJSON(w, req, start, fmt.Errorf("INVALID_PARAMS: Usuario, clave y mail no pueden estar vacíos"), http.StatusBadRequest)
    } else {
      // me fijo si no existe el usuarioRegistrar
      err, httpStat := UsuarioExiste(usuarioRegistrar.Usuario)
      if err != nil {
        core.RespErrorJSON(w, req, start, err, httpStat)
      } else {
        // establezco los campos
        var usuario models.Usuario
      	objID := bson.NewObjectId()
      	usuario.ID = objID
        usuario.Usuario = usuarioRegistrar.Usuario
        usuario.Clave = core.HashSha512(usuarioRegistrar.Clave)
        usuario.Mail = usuarioRegistrar.Mail

        // Genero una nueva sesión Mongo
        session := core.GetMongoSession()
        defer session.Close()

        // Intento el alta
        collection := session.DB(config.DB_Name).C(config.DB_Usuario)
        err = collection.Insert(usuario)
        if err != nil {
          core.RespErrorJSON(w, req, start, err, http.StatusInternalServerError)
        } else {
          core.RespOkJSON(w, req, start, "Ok", http.StatusCreated)
        }
      }
    }
  }
  return
}

func UsuarioExiste(usuarioExiste string) (error, int) {
  var usuario models.Usuario
  // Genero una nueva sesión Mongo
  session := core.GetMongoSession()
  defer session.Close()

  // Me aseguro el índice
  collection := session.DB(config.DB_Name).C(config.DB_Usuario)
  index := mgo.Index{
    Key:        []string{"usuario"},
    Unique:     true,
    DropDups:   true,
    Background: true,
    Sparse:     true,
  }
  err := collection.EnsureIndex(index)
  if err != nil {
    return err, http.StatusInternalServerError
  }

  collection.Find(bson.M{"usuario": usuarioExiste}).One(&usuario)
  if usuario.ID == "" {
    return nil, http.StatusOK
  } else {
    return fmt.Errorf("INVALID_PARAMS: El usuario ya existe"), http.StatusBadRequest
  }
}