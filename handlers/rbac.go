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

func PermisoAlta(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	var permisoAlta models.PermisoAlta

  decoder := json.NewDecoder(req.Body)
  err := decoder.Decode(&permisoAlta)
  if err != nil {
    core.RespErrorJSON(w, req, start, fmt.Errorf("INVALID_PARAMS: JSON decode erróneo"), http.StatusBadRequest)
  } else {
    // campos obligatorios
    if permisoAlta.Permiso == "" {
      core.RespErrorJSON(w, req, start, fmt.Errorf("INVALID_PARAMS: Permiso no puede estar vacío"), http.StatusBadRequest)
    } else {
      // me fijo si no existe el usuario
      err, httpStat := PermisoExiste(permisoAlta.Permiso)
      if err != nil {
        core.RespErrorJSON(w, req, start, err, httpStat)
      } else {
        // establezco los campos
        var permiso models.Permiso
      	objID := bson.NewObjectId()
      	permiso.ID = objID
        permiso.Permiso = permisoAlta.Permiso
        if permiso.Activo == "" {
          permiso.Activo = true
        } else {
          permiso.Activo = permisoAlta.Activo
        }
        if permiso.Borrado == "" {
          permiso.Borrado = false
        } else {
          permiso.Borrado = permisoAlta.Borrado
        }

        // Genero una nueva sesión Mongo
        session, err, httpStat := core.GetMongoSession()
        if err != nil {
          core.RespErrorJSON(w, req, start, err, httpStat)
        } else {
          defer session.Close()

          // Intento el alta
          collection := session.DB(config.DB_Name).C(config.DB_Permiso)
          err = collection.Insert(permiso)
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

func PermisoExiste(permisoExiste string) (error, int) {
  var permiso models.Permiso
  // Genero una nueva sesión Mongo
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    return err, httpStat
  } else {
    defer session.Close()
    collection := session.DB(config.DB_Name).C(config.DB_Usuario)

    // Me aseguro el índice
    index := mgo.Index{
      Key:        []string{"permiso"},
      Unique:     true,
      DropDups:   false,
      Background: true,
      Sparse:     true,
    }
    err := collection.EnsureIndex(index)
    if err != nil {
      s := []string{"INTERNAL_SERVER_ERROR:", err.Error()}
      return fmt.Errorf(strings.Join(s, " ")), http.StatusInternalServerError
    }

    collection.Find(bson.M{"permiso": permisoExiste}).One(&permiso)
    if permiso.ID == "" {
      return nil, http.StatusOK
    } else {
      return fmt.Errorf("INVALID_PARAMS: El usuario ya existe"), http.StatusBadRequest
    }
  }
}
