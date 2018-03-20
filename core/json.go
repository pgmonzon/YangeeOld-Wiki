package core

import (
  "log"
  "net/http"
  "time"
  "encoding/json"

  "github.com/pgmonzon/Yangee/models"
  "github.com/pgmonzon/Yangee/config"

  "github.com/gorilla/context"
  "gopkg.in/mgo.v2/bson"
)

func RespuestaJSON(w http.ResponseWriter, req *http.Request, start time.Time, respuesta []byte, code int) {
  w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
  if string(respuesta) != "" {
		w.Write(respuesta)
	}

  RegistrarCicloDeVida(req, code, start)

  log.Printf("%s\t%s\t%s\t%s\t%d\t%d\t%s",
		req.RemoteAddr,
		req.Method,
		req.RequestURI,
		req.Proto,
		code,
		len(respuesta),
		time.Since(start),
	)
}

func RegistrarCicloDeVida(req *http.Request, code int, start time.Time) {
  session, err, _ := GetMongoSession()
  if err == nil {
    defer session.Close()
    var cicloVida models.CicloDeVida
    cicloVida.ID = context.Get(req, "CicloDeVida_id").(bson.ObjectId)
    cicloVida.Fecha = context.Get(req, "Start").(time.Time)
    cicloVida.RemoteAddr = req.RemoteAddr
    cicloVida.Metodo = req.Method
    cicloVida.RequestURI = req.RequestURI
    cicloVida.Protocolo = req.Proto
    cicloVida.Codigo = code
    cicloVida.Duracion = time.Since(start)
    //bodyBytes, _ := ioutil.ReadAll(req.Body)
    //bodyString := string(bodyBytes)
    //cicloVida.Body = bodyString
    if context.Get(req, "ClienteAPI_id") == nil {
      cicloVida.ClienteAPI_id = bson.ObjectIdHex(config.CliAPI_id)
      cicloVida.ClienteAPI = config.CliAPI
    } else {
      cicloVida.ClienteAPI_id = context.Get(req, "ClienteAPI_id").(bson.ObjectId)
      cicloVida.ClienteAPI = context.Get(req, "ClienteAPI").(string)
    }
    if context.Get(req, "Usuario_id") == nil {
      cicloVida.Usuario_id = bson.ObjectIdHex(config.Usr_id)
      cicloVida.Usuario = config.Usr
    } else {
      cicloVida.Usuario_id = context.Get(req, "Usuario_id").(bson.ObjectId)
      cicloVida.Usuario = context.Get(req, "Usuario").(string)
    }
    cicloVida.Proceso = context.Get(req, "Proceso").(string)
    cicloVida.Coleccion = context.Get(req, "Coleccion").(string)
    if context.Get(req, "Objeto_id") == nil {
      cicloVida.Objeto_id = bson.ObjectIdHex(config.Obj_id)
    } else {
      cicloVida.ClienteAPI_id = context.Get(req, "Objeto_id").(bson.ObjectId)
    }
    cicloVida.Novedad = context.Get(req, "Novedad").(string)

    collection := session.DB(config.DB_Name).C(config.DB_CicloVida)
    err :=collection.Insert(cicloVida)
    if err != nil {
      log.Printf("Falló Insert ciclo de vida")
    }
  } else {
    log.Printf("Falló GetMongoSession ciclo de vida")
  }

}

func RespErrorJSON(w http.ResponseWriter, req *http.Request, start time.Time, err error, httpStat int) {
  var resp models.Respuesta
  resp.EstadoGral = "ERROR"
  resp.Mensaje = err.Error()
  respuesta, error := json.Marshal(resp)
  FatalErr(error)
  RespuestaJSON(w, req, start, respuesta, httpStat)
}

func RespOkJSON(w http.ResponseWriter, req *http.Request, start time.Time, mensaje string, httpStat int) {
  var resp models.Respuesta
  resp.EstadoGral = "OK"
  resp.Mensaje = mensaje
  respuesta, error := json.Marshal(resp)
  FatalErr(error)
  RespuestaJSON(w, req, start, respuesta, httpStat)
}
