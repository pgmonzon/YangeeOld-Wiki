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

func PermisoAgregar(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	var Permisos models.Permisos
  var resp models.Resp
  var mensaje models.Mensaje

  decoder := json.NewDecoder(req.Body)
  err := decoder.Decode(&Permisos)
  if err != nil {
    resp.EstadoGral = "ERROR"
    mensaje.Valor = "JSON"
    mensaje.Estado = "ERROR"
    mensaje.Mensaje = "INVALID_PARAMS: JSON decode erróneo"
    resp.Mensajes = append(resp.Mensajes, mensaje)
    respuesta, error := json.Marshal(resp)
    core.FatalErr(error)
    core.RespuestaJSON(w, req, start, respuesta, http.StatusBadRequest)
    return
  } else {
    // Genero una nueva sesión Mongo
    session, err, httpStat := core.GetMongoSession()
    if err != nil {
      resp.EstadoGral = "ERROR"
      mensaje.Valor = "MongoSession"
      mensaje.Estado = "ERROR"
      mensaje.Mensaje = err.Error()
      resp.Mensajes = append(resp.Mensajes, mensaje)
      respuesta, error := json.Marshal(resp)
      core.FatalErr(error)
      core.RespuestaJSON(w, req, start, respuesta, httpStat)
      return
    } else {
      defer session.Close()

      // Recorro el JSON
      for _, item := range Permisos.Permisos {
        if item.Permiso == "" {
          if resp.EstadoGral != "PARCIAL" {
            if resp.EstadoGral == "OK" {
              resp.EstadoGral = "PARCIAL"
            } else {
              resp.EstadoGral = "ERROR"
            }
          }
          mensaje.Valor = item.Permiso
          mensaje.Estado = "ERROR"
          s := []string{"INTERNAL_SERVER_ERROR: ", "El campo Permiso no puede estar vacío"}
          mensaje.Mensaje = strings.Join(s, "")
          resp.Mensajes = append(resp.Mensajes, mensaje)
        } else {
          // Me fijo si ya existe
          err := PermisoExiste(item.Permiso)
          if err != nil {
            if resp.EstadoGral != "PARCIAL" {
              if resp.EstadoGral == "OK" {
                resp.EstadoGral = "PARCIAL"
              } else {
                resp.EstadoGral = "ERROR"
              }
            }
            mensaje.Valor = item.Permiso
            mensaje.Estado = "ERROR"
            mensaje.Mensaje = err.Error()
            resp.Mensajes = append(resp.Mensajes, mensaje)
          } else {
            objID := bson.NewObjectId()
          	item.ID = objID

            // Intento el alta
            collection := session.DB(config.DB_Name).C(config.DB_Permiso)
            err = collection.Insert(item)
            if err != nil {
              if resp.EstadoGral != "PARCIAL" {
                if resp.EstadoGral == "OK" {
                  resp.EstadoGral = "PARCIAL"
                } else {
                  resp.EstadoGral = "ERROR"
                }
              }
              mensaje.Valor = item.Permiso
              mensaje.Estado = "ERROR"
              s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
              mensaje.Mensaje = strings.Join(s, "")
              resp.Mensajes = append(resp.Mensajes, mensaje)
            } else {
              if resp.EstadoGral != "PARCIAL" {
                if resp.EstadoGral == "ERROR" {
                  resp.EstadoGral = "PARCIAL"
                } else {
                  resp.EstadoGral = "OK"
                }
              }
              mensaje.Valor = item.Permiso
              mensaje.Estado = "OK"
              mensaje.Mensaje = "OK"
              resp.Mensajes = append(resp.Mensajes, mensaje)
            }
          }
        }
      }
      var httpStat int
      if resp.EstadoGral == "OK" {
        httpStat = http.StatusCreated
      } else {
        httpStat = http.StatusOK
      }
      respuesta, error := json.Marshal(resp)
      core.FatalErr(error)
      core.RespuestaJSON(w, req, start, respuesta, httpStat)
      return
    }
  }
}

func PermisoExiste(permisoExiste string) (error) {
  var permiso models.Permiso
  // Genero una nueva sesión Mongo
  session, err, _ := core.GetMongoSession()
  if err != nil {
    return err
  } else {
    defer session.Close()
    collection := session.DB(config.DB_Name).C(config.DB_Permiso)

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
      s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
      return fmt.Errorf(strings.Join(s, ""))
    }

    collection.Find(bson.M{"permiso": permisoExiste}).One(&permiso)
    if permiso.ID == "" {
      return nil
    } else {
      s := []string{"INVALID_PARAMS: El permiso [", permisoExiste,"] ya existe"}
      return fmt.Errorf(strings.Join(s, ""))
    }
  }
}

func RolAgregar(w http.ResponseWriter, req *http.Request) {
  start := time.Now()
	var Roles models.Roles
  var resp models.Resp
  var mensaje models.Mensaje

  decoder := json.NewDecoder(req.Body)
  err := decoder.Decode(&Roles)
  if err != nil {
    resp.EstadoGral = "ERROR"
    mensaje.Valor = "JSON"
    mensaje.Estado = "ERROR"
    mensaje.Mensaje = "INVALID_PARAMS: JSON decode erróneo"
    resp.Mensajes = append(resp.Mensajes, mensaje)
    respuesta, error := json.Marshal(resp)
    core.FatalErr(error)
    core.RespuestaJSON(w, req, start, respuesta, http.StatusBadRequest)
    return
  } else {
    // Genero una nueva sesión Mongo
    session, err, httpStat := core.GetMongoSession()
    if err != nil {
      resp.EstadoGral = "ERROR"
      mensaje.Valor = "MongoSession"
      mensaje.Estado = "ERROR"
      mensaje.Mensaje = err.Error()
      resp.Mensajes = append(resp.Mensajes, mensaje)
      respuesta, error := json.Marshal(resp)
      core.FatalErr(error)
      core.RespuestaJSON(w, req, start, respuesta, httpStat)
      return
    } else {
      defer session.Close()

      // Recorro el JSON
      for _, item := range Roles.Roles {
        if item.Rol == "" {
          if resp.EstadoGral != "PARCIAL" {
            if resp.EstadoGral == "OK" {
              resp.EstadoGral = "PARCIAL"
            } else {
              resp.EstadoGral = "ERROR"
            }
          }
          mensaje.Valor = item.Rol
          mensaje.Estado = "ERROR"
          s := []string{"INTERNAL_SERVER_ERROR: ", "El campo Rol no puede estar vacío"}
          mensaje.Mensaje = strings.Join(s, "")
          resp.Mensajes = append(resp.Mensajes, mensaje)
        } else {
          // Me fijo si ya existe
          err := RolExiste(item.Rol)
          if err != nil {
            if resp.EstadoGral != "PARCIAL" {
              if resp.EstadoGral == "OK" {
                resp.EstadoGral = "PARCIAL"
              } else {
                resp.EstadoGral = "ERROR"
              }
            }
            mensaje.Valor = item.Rol
            mensaje.Estado = "ERROR"
            mensaje.Mensaje = err.Error()
            resp.Mensajes = append(resp.Mensajes, mensaje)
          } else {
            objID := bson.NewObjectId()
          	item.ID = objID

            // Intento el alta
            collection := session.DB(config.DB_Name).C(config.DB_Rol)
            err = collection.Insert(item)
            if err != nil {
              if resp.EstadoGral != "PARCIAL" {
                if resp.EstadoGral == "OK" {
                  resp.EstadoGral = "PARCIAL"
                } else {
                  resp.EstadoGral = "ERROR"
                }
              }
              mensaje.Valor = item.Rol
              mensaje.Estado = "ERROR"
              s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
              mensaje.Mensaje = strings.Join(s, "")
              resp.Mensajes = append(resp.Mensajes, mensaje)
            } else {
              if resp.EstadoGral != "PARCIAL" {
                if resp.EstadoGral == "ERROR" {
                  resp.EstadoGral = "PARCIAL"
                } else {
                  resp.EstadoGral = "OK"
                }
              }
              mensaje.Valor = item.Rol
              mensaje.Estado = "OK"
              mensaje.Mensaje = "OK"
              resp.Mensajes = append(resp.Mensajes, mensaje)
            }
          }
        }
      }
      var httpStat int
      if resp.EstadoGral == "OK" {
        httpStat = http.StatusCreated
      } else {
        httpStat = http.StatusOK
      }
      respuesta, error := json.Marshal(resp)
      core.FatalErr(error)
      core.RespuestaJSON(w, req, start, respuesta, httpStat)
      return
    }
  }
}

func RolExiste(rolExiste string) (error) {
  var rol models.Rol
  // Genero una nueva sesión Mongo
  session, err, _ := core.GetMongoSession()
  if err != nil {
    return err
  } else {
    defer session.Close()
    collection := session.DB(config.DB_Name).C(config.DB_Rol)

    // Me aseguro el índice
    index := mgo.Index{
      Key:        []string{"rol"},
      Unique:     true,
      DropDups:   false,
      Background: true,
      Sparse:     true,
    }
    err := collection.EnsureIndex(index)
    if err != nil {
      s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
      return fmt.Errorf(strings.Join(s, ""))
    }

    collection.Find(bson.M{"rol": rolExiste}).One(&rol)
    if rol.ID == "" {
      return nil
    } else {
      s := []string{"INVALID_PARAMS: El rol [", rolExiste,"] ya existe"}
      return fmt.Errorf(strings.Join(s, ""))
    }
  }
}
