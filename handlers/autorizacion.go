package handlers

import (
  "net/http"
  "encoding/json"
  "time"

  "github.com/pgmonzon/Yangee/core"
)

func Autorizar(w http.ResponseWriter, req *http.Request) {
  start := time.Now()

  token, err, httpStat := core.ObtenerToken(req.Header.Get("authorization"))
  if err != nil {
    core.RespErrorJSON(w, req, start, err, httpStat)
  } else {
    respuesta, error := json.Marshal(token)
    if error != nil {
      core.RespErrorJSON(w, req, start, error, httpStat)
    } else {
      core.RespuestaJSON(w, req, start, respuesta, httpStat)
    }
  }

  return
}
