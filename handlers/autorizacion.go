package handlers

import (
  "net/http"
  "encoding/json"
  "time"

  "github.com/pgmonzon/Yangee/models"
  "github.com/pgmonzon/Yangee/core"
)

func Autorizar(w http.ResponseWriter, req *http.Request) {
  start := time.Now()

  token, err, httpStat := core.ObtenerToken(req.Header.Get("authorization"))
  if err != nil {
    var resp models.Error
    resp.Estado = "ERROR"
    resp.Detalle = err.Error()
    respuesta, error := json.Marshal(resp)
    core.FatalErr(error)
    core.RespuestaJSON(w, req, start, respuesta, httpStat)
  } else {
    respuesta, error := json.Marshal(token)
    core.FatalErr(error)
    core.RespuestaJSON(w, req, start, respuesta, httpStat)
  }
}
