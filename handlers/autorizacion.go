package handlers

import (
  "net/http"
  "encoding/json"
  "time"
  "fmt"
  "strings"

  "github.com/pgmonzon/Yangee/core"
)

func Autorizar(w http.ResponseWriter, req *http.Request) {
  start := time.Now()

  aut, err, httpStat := ValidarTokenCliente(w, req)
  if err != nil {
    core.RespErrorJSON(w, req, start, err, httpStat)
  } else {
    token, err, httpStat := GenerarToken(aut)
    if err != nil {
      core.RespErrorJSON(w, req, start, err, httpStat)
    } else {
      respuesta, err := json.Marshal(token)
      if err != nil {
        s := []string{"INTERNAL_SERVER_ERROR:", err.Error()}
        core.RespErrorJSON(w, req, start, fmt.Errorf(strings.Join(s, " ")), httpStat)
      } else {
        core.RespuestaJSON(w, req, start, respuesta, httpStat)
      }
    }
  }
  return
}
