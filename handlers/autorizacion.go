package handlers

import (
  "net/http"
  "encoding/json"
  "time"

  "github.com/pgmonzon/Yangee/models"
  "github.com/pgmonzon/Yangee/core"
)

// Valida las credenciales del usuario y devuelve un Token
func Autorizar(w http.ResponseWriter, req *http.Request) {
  start := time.Now()

  var aut models.Autorizar
  err := json.NewDecoder(req.Body).Decode(&aut)

  if err != nil || aut.Usuario == "" || aut.Clave == "" {
    var error models.Error
    error.Estado = "ERROR"
    error.Detalle = "Parámetros Incorrectos"
    respuesta, err := json.Marshal(error)
    core.FatalErr(err)
    core.RespuestaJSON(w, req, start, respuesta, http.StatusBadRequest)
  } else {
    var token models.Token
    token.Token, err = core.CrearToken(aut)
    if err != nil {
      var error models.Error
      error.Estado = "ERROR"
      error.Detalle = token.Token
      respuesta, err := json.Marshal(error)
      core.FatalErr(err)
      core.RespuestaJSON(w, req, start, respuesta, http.StatusInternalServerError)
    } else {
      respuesta, err := json.Marshal(token)
      core.FatalErr(err)
      core.RespuestaJSON(w, req, start, respuesta, http.StatusOK)
    }
  }
}
