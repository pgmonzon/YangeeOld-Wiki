package handlers

import (
  "net/http"
  "encoding/json"
  "strings"

  "github.com/pgmonzon/Yangee/core"
  "github.com/pgmonzon/Yangee/models"

  "github.com/gorilla/context"
  "gopkg.in/mgo.v2/bson"
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

  // Busco el logo de la Empresa
  // ***************************
  empresaID := context.Get(req, "Empresa_id").(bson.ObjectId)
  empresa, err, httpStat := Empresa_X_ID(empresaID)
  if err != nil {
    core.RspMsgJSON(w, req, "ERROR", "Buscando Empresa", err.Error(), httpStat)
    return
  }

  // Busco los datos del usuario
  // ***************************
  usuarioID := context.Get(req, "Usuario_id").(bson.ObjectId)
  usuario, err, httpStat := Usuario_X_ID(usuarioID)
  if err != nil {
    core.RspMsgJSON(w, req, "ERROR", "Buscando Usuario", err.Error(), httpStat)
    return
  }

  // Hago el marshal de la autorización
  // **********************************
  var autorizacion models.Autorizacion
  autorizacion.Token = token.Token
  autorizacion.Logo = empresa.Logo
  usr := []string{usuario.Nombre, " ", usuario.Apellido}
  strUsr := strings.Join(usr, "")
  if strUsr == " " {
    autorizacion.Usuario = usuario.Usuario
  } else {
    autorizacion.Usuario = strUsr
  }
  autorizacion.Menu = usuario.Menu
  respuesta, err := json.Marshal(autorizacion)
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
