package handlers

import (
  "net/http"
  "encoding/json"
  "strings"

  "github.com/pgmonzon/Yangee/core"
  "github.com/pgmonzon/Yangee/models"
)


func Autorizar(w http.ResponseWriter, req *http.Request) {

  // Valido el token del clienteAPI
  // ******************************
  estado, valor, mensaje, httpStat, aut, usuario, empresa := ValidarTokenCliente(w, req)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
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

  // Está todo Ok
  // ************
  core.RspJSON(w, req, respuesta, http.StatusOK)
  return
}
/*
func InvitacionEmpresa(w http.ResponseWriter, req *http.Request) {
	var invitacionEmpresa models.InvitacionEmpresa

  // Decode del JSON
  // ***************
  decoder := json.NewDecoder(req.Body)
  err := decoder.Decode(&invitacionEmpresa)
  if err != nil {
    core.RspMsgJSON(w, req, "ERROR", "JSON", "INVALID_PARAMS: JSON decode erróneo", http.StatusBadRequest)
    return
  }

  // Verifico los campos obligatorios
  // ********************************
  if invitacionEmpresa.Empresa == "" || invitacionEmpresa.Rol == "" || invitacionEmpresa.Mail == "" {
    core.RspMsgJSON(w, req, "ERROR", "InvitacionEmpresa", "INVALID_PARAMS: empresa, rol y mail no pueden estar vacíos", http.StatusBadRequest)
    return
  }

  // Doy de alta la empresa
  // **********************
  empresa.Empresa = invitacionEmpresa.Empresa
  s := []string{invitacionEmpresa.Empresa, ".jpg"}
  empresa.Logo = strings.Join(s, "")
  empresa.Activo = true
  estado, valor, mensaje, httpStat, empresa, _ := EmpresaAlta(empresa)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
  }

  // Establezco las variables
  // ************************
  context.Set(req, "TipoOper", "#Novedad#")
  context.Set(req, "Coleccion", config.DB_ClienteAPI)
  context.Set(req, "Objeto_id", clienteAPI.ID)
  context.Set(req, "Audit", clienteAPI)

  // Está todo Ok
  // ************
  core.RspMsgJSON(w, req, "OK", clienteAPI.ClienteAPI, "Ok", http.StatusCreated)
  return
}
*/
