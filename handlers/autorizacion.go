package handlers

import (
  "net/http"
  "encoding/json"
  "strings"

  "github.com/pgmonzon/Yangee/core"
  "github.com/gorilla/context"
)


func Autorizar(w http.ResponseWriter, req *http.Request) {

  // Valido el token del clienteAPI
  // ******************************
  aut, err, httpStat := ValidarTokenCliente(w, req)
  if err != nil {
    s := []string{err.Error()}
    core.RspMsgJSON(w, req, "ERROR", "ValidarTokenCliente", strings.Join(s, ""), httpStat)
    return
  }

  // Genero el token para operar
  // ***************************
  token, err, httpStat := GenerarToken(aut, req)
  if err != nil {
    s := []string{err.Error()}
    core.RspMsgJSON(w, req, "ERROR", "GenerarToken", strings.Join(s, ""), httpStat)
    return
  }

  // Hago el marshal del Token
  // *************************
  respuesta, err := json.Marshal(token)
  if err != nil {
    s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
    core.RspMsgJSON(w, req, "ERROR", "Marshal", strings.Join(s, ""), http.StatusInternalServerError)
    return
  }

  // Establezco las variables
  // ************************
  context.Set(req, "TipoOper", "#Login#")
  s := []string{"Ingresó al sistema"}
  context.Set(req, "Novedad", strings.Join(s, ""))

  // Está todo Ok
  // ************
  core.RspJSON(w, req, respuesta, http.StatusOK)
  return
}
