package handlers

import (
  "time"
  "encoding/json"
  "net/http"
  "fmt"

  "github.com/pgmonzon/Yangee/models"
  "github.com/pgmonzon/Yangee/core"
  "github.com/pgmonzon/Yangee/config"

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
        return
      }
      core.RespOkJSON(w, req, start, "Ok", http.StatusCreated)
    }
  }
  return
}
