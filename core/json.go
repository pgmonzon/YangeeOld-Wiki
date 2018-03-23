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

func RspMsgJSON(w http.ResponseWriter, req *http.Request, estado string, valor string, mensaje string, httpStat int) {
  var rsp models.Respuesta

  // Establezco las variables
  // ************************
  context.Set(req, "Novedad", mensaje)

  rsp.Estado = estado
  rsp.Valor = valor
  rsp.Mensaje = mensaje
  respuesta, error := json.Marshal(rsp)
  FatalErr(error)
  RspJSON(w, req, respuesta, httpStat)
}

func RspJSON(w http.ResponseWriter, req *http.Request, respuesta []byte, code int) {
  w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
  if string(respuesta) != "" {
		w.Write(respuesta)
	}

  start := context.Get(req, "Start").(time.Time)

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
    cicloVida.Fecha = start
    cicloVida.RemoteAddr = req.RemoteAddr
    cicloVida.Metodo = req.Method
    cicloVida.RequestURI = req.RequestURI
    cicloVida.Protocolo = req.Proto
    cicloVida.Codigo = code
    cicloVida.Duracion = time.Since(start)
    cicloVida.ClienteAPI_id = context.Get(req, "ClienteAPI_id").(bson.ObjectId)
    cicloVida.ClienteAPI = context.Get(req, "ClienteAPI").(string)
    cicloVida.Usuario_id = context.Get(req, "Usuario_id").(bson.ObjectId)
    cicloVida.Usuario = context.Get(req, "Usuario").(string)
    cicloVida.TipoOper = context.Get(req, "TipoOper").(string)
    cicloVida.Coleccion = context.Get(req, "Coleccion").(string)
    cicloVida.Novedad = context.Get(req, "Novedad").(string)
    cicloVida.Objeto_id = context.Get(req, "Objeto_id").(bson.ObjectId)
    cicloVida.Audit = context.Get(req, "Audit")

    collection := session.DB(config.DB_Name).C(config.DB_CicloVida)
    err :=collection.Insert(cicloVida)
    if err != nil {
      log.Printf("Falló Insert ciclo de vida")
    }
  } else {
    log.Printf("Falló GetMongoSession ciclo de vida")
  }
}
