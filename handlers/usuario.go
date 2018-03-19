package handlers

import (
  "time"
  "encoding/json"
  "net/http"
  "fmt"
  "strings"
  "strconv"
  "log"

  "github.com/pgmonzon/Yangee/models"
  "github.com/pgmonzon/Yangee/core"
  "github.com/pgmonzon/Yangee/config"

  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "github.com/gorilla/context"
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

func UsuarioPermisos(usuarioPermisos string) (string, error, int) {
  // Genero una nueva sesión Mongo
  session, err, _ := core.GetMongoSession()
  if err != nil {
    s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
    return "", fmt.Errorf(strings.Join(s, "")), http.StatusInternalServerError
  } else {
    defer session.Close()
    cUsuario := session.DB(config.DB_Name).C(config.DB_Usuario)

    // Busco el usuario y verifico que esté activo
    var usuario models.Usuario
    cUsuario.Find(bson.M{"usuario": usuarioPermisos, "activo": true, "borrado": false}).One(&usuario)
    if usuario.ID == "" {
      s := []string{"INVALID_PARAMS: El usuario no existe o está inactivo"}
      return "", fmt.Errorf(strings.Join(s, "")), http.StatusBadRequest
    } else {
      // Obtengo los ID roles del usuario
      rolesArr := []bson.ObjectId{}
      for _, item := range usuario.Roles {
        if item.ID != "" {
          rolesArr = append(rolesArr, item.ID)
        }
      }

      roles := make([]models.Rol, 0)
      cRoles := session.DB(config.DB_Name).C(config.DB_Rol)
      cRoles.Find(bson.M{"_id": bson.M{"$in": rolesArr}}).All(&roles)

      // Obtengo los ID permisos de los roles
      permisosArr := []bson.ObjectId{}
      for _, itemRol := range roles {
        for _, itemPermiso := range itemRol.Permisos {
          if itemPermiso.ID != "" {
            permisosArr = append(permisosArr, itemPermiso.ID)
          }
        }
      }

      permisos := make([]models.Permiso, 0)
      cPermisos := session.DB(config.DB_Name).C(config.DB_Permiso)
      cPermisos.Find(bson.M{"_id": bson.M{"$in": permisosArr}}).All(&permisos)

      // Junto los permisos en un string
      permisosStr := []string{}
      permisosStr = append(permisosStr, "#")
      for _, itemItem := range permisos {
        if itemItem.Permiso != "" {
          permisosStr = append(permisosStr, itemItem.Permiso)
        }
      }
      permisosStr = append(permisosStr, "#")
      permisosUsuario := strings.Join(permisosStr, "#")
      return permisosUsuario, nil, http.StatusOK
    }
  }
}

func UsuarioLogin(usuarioLogin string, claveLogin string, req *http.Request) (error, int) {
  var usuario models.Usuario
  // Genero una nueva sesión Mongo
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    return err, httpStat
  } else {
    defer session.Close()
    collection := session.DB(config.DB_Name).C(config.DB_Usuario)
    collection.Find(bson.M{"usuario": usuarioLogin, "clave": strconv.FormatInt(core.HashSha512(claveLogin),16), "activo": true, "borrado": false}).One(&usuario)
    if usuario.ID == "" {
      s := []string{"FORBIDDEN: ", "Usuario y clave incorrectos"}
      return fmt.Errorf(strings.Join(s, "")), http.StatusForbidden
    } else {
      context.Set(req, "Usuario_id", usuario.ID)
      context.Set(req, "Usuario", usuario.Usuario)
      return nil, http.StatusOK
    }
  }
}

func TestPermisos(w http.ResponseWriter, req *http.Request) {
  start := time.Now()
  var resp models.Resp
  var mensaje models.Mensaje
  usuarioPermisos := "patricio"
  // Genero una nueva sesión Mongo
  session, err, _ := core.GetMongoSession()
  if err != nil {
    resp.EstadoGral = "ERROR"
    mensaje.Valor = "MongoSession"
    mensaje.Estado = "ERROR"
    mensaje.Mensaje = err.Error()
    resp.Mensajes = append(resp.Mensajes, mensaje)
    respuesta, error := json.Marshal(resp)
    core.FatalErr(error)
    core.RespuestaJSON(w, req, start, respuesta, http.StatusBadRequest)
    return
  } else {
    defer session.Close()
    cUsuario := session.DB(config.DB_Name).C(config.DB_Usuario)

    // Busco el usuario y verifico que esté activo
    var usuario models.Usuario
    cUsuario.Find(bson.M{"usuario": usuarioPermisos, "activo": true, "borrado": false}).One(&usuario)
    if usuario.ID == "" {
      resp.EstadoGral = "ERROR"
      mensaje.Valor = "Usuario inválido"
      mensaje.Estado = "ERROR"
      mensaje.Mensaje = "Usuario inválido"
      resp.Mensajes = append(resp.Mensajes, mensaje)
      respuesta, error := json.Marshal(resp)
      core.FatalErr(error)
      core.RespuestaJSON(w, req, start, respuesta, http.StatusBadRequest)
      return
    } else {
      // Obtengo los ID roles del usuario
      rolesArr := []bson.ObjectId{}
      for _, item := range usuario.Roles {
        if item.ID != "" {
          rolesArr = append(rolesArr, item.ID)
        }
      }

      roles := make([]models.Rol, 0)
      cRoles := session.DB(config.DB_Name).C(config.DB_Rol)
      cRoles.Find(bson.M{"_id": bson.M{"$in": rolesArr}}).All(&roles)

      // Obtengo los ID permisos de los roles
      permisosArr := []bson.ObjectId{}
      for _, itemRol := range roles {
        for _, itemPermiso := range itemRol.Permisos {
          if itemPermiso.ID != "" {
            permisosArr = append(permisosArr, itemPermiso.ID)
          }
        }
      }

      permisos := make([]models.Permiso, 0)
      cPermisos := session.DB(config.DB_Name).C(config.DB_Permiso)
      cPermisos.Find(bson.M{"_id": bson.M{"$in": permisosArr}}).All(&permisos)

      // Junto los permisos en un string
      permisosStr := []string{}
      permisosStr = append(permisosStr, "#")
      for _, itemItem := range permisos {
        if itemItem.Permiso != "" {
          permisosStr = append(permisosStr, itemItem.Permiso)
        }
      }
      permisosStr = append(permisosStr, "#")
      permisosUsuario := strings.Join(permisosStr, "#")
      //return permisosUsuario, nil, http.StatusOK
      log.Printf(permisosUsuario)

      respuesta, error := json.Marshal(permisos)
      core.FatalErr(error)
      core.RespuestaJSON(w, req, start, respuesta, http.StatusOK)
      return
    }
  }
}
