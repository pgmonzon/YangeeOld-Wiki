package handlers

import (
  "net/http"
  "encoding/json"
  "time"
  "strings"

  "github.com/pgmonzon/Yangee/core"
  "github.com/pgmonzon/Yangee/models"
)

func Autorizar(w http.ResponseWriter, req *http.Request) {
  start := time.Now()
  var resp models.Resp
  var mensaje models.Mensaje

  aut, err, httpStat := ValidarTokenCliente(w, req)
  if err != nil {
    resp.EstadoGral = "ERROR"
    mensaje.Valor = "ValidarTokenCliente"
    mensaje.Estado = "ERROR"
    s := []string{err.Error()}
    mensaje.Mensaje = strings.Join(s, "")
    resp.Mensajes = append(resp.Mensajes, mensaje)
    respuesta, error := json.Marshal(resp)
    core.FatalErr(error)
    core.RespuestaJSON(w, req, start, respuesta, httpStat)
    return
  } else {
    token, err, httpStat := GenerarToken(aut)
    if err != nil {
      resp.EstadoGral = "ERROR"
      mensaje.Valor = "GenerarToken"
      mensaje.Estado = "ERROR"
      s := []string{err.Error()}
      mensaje.Mensaje = strings.Join(s, "")
      resp.Mensajes = append(resp.Mensajes, mensaje)
      respuesta, error := json.Marshal(resp)
      core.FatalErr(error)
      core.RespuestaJSON(w, req, start, respuesta, httpStat)
      return
    } else {
      respuesta, err := json.Marshal(token)
      if err != nil {
        resp.EstadoGral = "ERROR"
        mensaje.Valor = "Marshal"
        mensaje.Estado = "ERROR"
        s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
        mensaje.Mensaje = strings.Join(s, "")
        resp.Mensajes = append(resp.Mensajes, mensaje)
        respuesta, error := json.Marshal(resp)
        core.FatalErr(error)
        core.RespuestaJSON(w, req, start, respuesta, http.StatusInternalServerError)
        return
      } else {
        core.RespuestaJSON(w, req, start, respuesta, http.StatusOK)
        return
      }
    }
  }
}
