package handlers

import (
  "time"
  "encoding/json"
  "net/http"
  "fmt"
  "strings"
  "strconv"

  "github.com/pgmonzon/Yangee/models"
  "github.com/pgmonzon/Yangee/core"
  "github.com/pgmonzon/Yangee/config"

  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

func UsuarioRegistrar(w http.ResponseWriter, req *http.Request) {
  start := time.Now()
	var Usuarios models.Usuarios
  var resp models.Resp
  var mensaje models.Mensaje

  decoder := json.NewDecoder(req.Body)
  err := decoder.Decode(&Usuarios)
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
      for _, item := range Usuarios.Usuarios {
        if item.Usuario == "" || item.Clave == "" || item.Mail == "" {
          if resp.EstadoGral != "PARCIAL" {
            if resp.EstadoGral == "OK" {
              resp.EstadoGral = "PARCIAL"
            } else {
              resp.EstadoGral = "ERROR"
            }
          }
          mensaje.Valor = item.Usuario
          mensaje.Estado = "ERROR"
          s := []string{"INTERNAL_SERVER_ERROR: ", "El campo usuario, clave y mail no pueden estar vacíos"}
          mensaje.Mensaje = strings.Join(s, "")
          resp.Mensajes = append(resp.Mensajes, mensaje)
        } else {
          // Me fijo si ya existe
          err := UsuarioExiste(item.Usuario)
          if err != nil {
            if resp.EstadoGral != "PARCIAL" {
              if resp.EstadoGral == "OK" {
                resp.EstadoGral = "PARCIAL"
              } else {
                resp.EstadoGral = "ERROR"
              }
            }
            mensaje.Valor = item.Usuario
            mensaje.Estado = "ERROR"
            mensaje.Mensaje = err.Error()
            resp.Mensajes = append(resp.Mensajes, mensaje)
          } else {
            objID := bson.NewObjectId()
          	item.ID = objID
            //item.Encrip = core.HashSha512(item.Clave)
            item.Clave = strconv.FormatInt(core.HashSha512(item.Clave),16)

            // Intento el alta
            collection := session.DB(config.DB_Name).C(config.DB_Usuario)
            err = collection.Insert(item)
            if err != nil {
              if resp.EstadoGral != "PARCIAL" {
                if resp.EstadoGral == "OK" {
                  resp.EstadoGral = "PARCIAL"
                } else {
                  resp.EstadoGral = "ERROR"
                }
              }
              mensaje.Valor = item.Usuario
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
              mensaje.Valor = item.Usuario
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

func UsuarioExiste(usuarioExiste string) (error) {
  var usuario models.Usuario
  // Genero una nueva sesión Mongo
  session, err, _ := core.GetMongoSession()
  if err != nil {
    return err
  } else {
    defer session.Close()
    collection := session.DB(config.DB_Name).C(config.DB_Usuario)

    // Me aseguro el índice
    index := mgo.Index{
      Key:        []string{"usuario"},
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

    collection.Find(bson.M{"usuario": usuarioExiste}).One(&usuario)
    if usuario.ID == "" {
      return nil
    } else {
      s := []string{"INVALID_PARAMS: El usuario [", usuarioExiste,"] ya existe"}
      return fmt.Errorf(strings.Join(s, ""))
    }
  }
}

////////////////////////////////////////////////////////////////////
/**
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
**/
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
      return fmt.Errorf("FORBIDDEN: usuario y clave incorrectos"), http.StatusForbidden
    } else {
      return nil, http.StatusOK
    }
  }
}
