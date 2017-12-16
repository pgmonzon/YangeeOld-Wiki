package core

import (
  "log"
  "net/http"
  "time"
  "encoding/json"

  "github.com/pgmonzon/Yangee/models"
)

func RespuestaJSON(w http.ResponseWriter, req *http.Request, start time.Time, respuesta []byte, code int) {
  w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
  if string(respuesta) != "" {
		w.Write(respuesta)
	}

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

func RespErrorJSON(w http.ResponseWriter, req *http.Request, start time.Time, err error, httpStat int) {
  var resp models.Respuesta
  resp.Estado = "ERROR"
  resp.Detalle = err.Error()
  respuesta, error := json.Marshal(resp)
  FatalErr(error)
  RespuestaJSON(w, req, start, respuesta, httpStat)
}

func RespOkJSON(w http.ResponseWriter, req *http.Request, start time.Time, detalle string, httpStat int) {
  var resp models.Respuesta
  resp.Estado = "OK"
  resp.Detalle = detalle
  respuesta, error := json.Marshal(resp)
  FatalErr(error)
  RespuestaJSON(w, req, start, respuesta, httpStat)
}
