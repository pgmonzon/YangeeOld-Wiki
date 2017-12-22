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
      // me fijo si no existe el usuario
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
        session, err, httpStat := core.GetMongoSession()
        if err != nil {
          core.RespErrorJSON(w, req, start, err, httpStat)
        } else {
          defer session.Close()

          // Intento el alta
          collection := session.DB(config.DB_Name).C(config.DB_Usuario)
          err = collection.Insert(usuario)
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

func UsuarioExiste(usuarioExiste string) (error, int) {
  var usuario models.Usuario
  // Genero una nueva sesión Mongo
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    return err, httpStat
  } else {
    defer session.Close()
    collection := session.DB(config.DB_Name).C(config.DB_Usuario)

    // Me aseguro el índice
    index := mgo.Index{
      Key:        []string{"usuario"},
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

    collection.Find(bson.M{"usuario": usuarioExiste}).One(&usuario)
    if usuario.ID == "" {
      return nil, http.StatusOK
    } else {
      return fmt.Errorf("INVALID_PARAMS: El usuario ya existe"), http.StatusBadRequest
    }
  }
}

func UsuarioLogin(usuarioLogin string, claveLogin string) (error, int) {
  var usuario models.Usuario
  // Genero una nueva sesión Mongo
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    return err, httpStat
  } else {
    defer session.Close()
    collection := session.DB(config.DB_Name).C(config.DB_Usuario)

    collection.Find(bson.M{"usuario": usuarioLogin, "clave": core.HashSha512(claveLogin)}).One(&usuario)
    if usuario.ID == "" {
      return fmt.Errorf("FORBIDDEN: usuario y clave incorrectos"), http.StatusBadRequest
    } else {
      return nil, http.StatusOK
    }
  }
}
