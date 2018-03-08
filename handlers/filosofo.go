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

func CrearFilosofo(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	var Filosofos models.Filosofos
  var resp models.Resp
  var mensaje models.Mensaje

  decoder := json.NewDecoder(req.Body)
  err := decoder.Decode(&Filosofos)
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
      for _, item := range Filosofos.Filosofos {
        if item.Filosofo == "" {
          if resp.EstadoGral != "PARCIAL" {
            if resp.EstadoGral == "OK" {
              resp.EstadoGral = "PARCIAL"
            } else {
              resp.EstadoGral = "ERROR"
            }
          }
          mensaje.Valor = item.Filosofo
          mensaje.Estado = "ERROR"
          s := []string{"INTERNAL_SERVER_ERROR: ", "Filósofo no puede estar vacío"}
          mensaje.Mensaje = strings.Join(s, "")
          resp.Mensajes = append(resp.Mensajes, mensaje)
        } else {
          // Me fijo si ya existe
          err := ExisteFilosofo(item.Filosofo)
          if err != nil {
            if resp.EstadoGral != "PARCIAL" {
              if resp.EstadoGral == "OK" {
                resp.EstadoGral = "PARCIAL"
              } else {
                resp.EstadoGral = "ERROR"
              }
            }
            mensaje.Valor = item.Filosofo
            mensaje.Estado = "ERROR"
            mensaje.Mensaje = err.Error()
            resp.Mensajes = append(resp.Mensajes, mensaje)
          } else {
            objID := bson.NewObjectId()
          	item.ID = objID

            // Intento el alta
            collection := session.DB(config.DB_Name).C(config.DB_Filosofo)
            err = collection.Insert(item)
            if err != nil {
              if resp.EstadoGral != "PARCIAL" {
                if resp.EstadoGral == "OK" {
                  resp.EstadoGral = "PARCIAL"
                } else {
                  resp.EstadoGral = "ERROR"
                }
              }
              mensaje.Valor = item.Filosofo
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
              mensaje.Valor = item.Filosofo
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

func ExisteFilosofo(filosofoExiste string) (error) {
  var filosofo models.Filosofo
  // Genero una nueva sesión Mongo
  session, err, _ := core.GetMongoSession()
  if err != nil {
    s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
    return fmt.Errorf(strings.Join(s, ""))
  } else {
    defer session.Close()
    collection := session.DB(config.DB_Name).C(config.DB_Filosofo)

    // Me aseguro el índice
    index := mgo.Index{
      Key:        []string{"filosofo"},
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

    collection.Find(bson.M{"filosofo": filosofoExiste}).One(&filosofo)
    if filosofo.ID == "" {
      return nil
    } else {
      if filosofo.Borrado == true {
        s := []string{"INVALID_PARAMS: El filósofo [", filosofoExiste,"] ya existe borrado"}
        return fmt.Errorf(strings.Join(s, ""))
      }
      if filosofo.Activo == false {
        s := []string{"INVALID_PARAMS: El filósofo [", filosofoExiste,"] ya existe inactivo"}
        return fmt.Errorf(strings.Join(s, ""))
      }
      s := []string{"INVALID_PARAMS: El filósofo [", filosofoExiste,"] ya existe"}
      return fmt.Errorf(strings.Join(s, ""))
    }
  }
}
