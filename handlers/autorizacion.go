package handlers

import (
  "net/http"
  "encoding/json"
  "strings"

  "github.com/pgmonzon/Yangee/core"
  "github.com/pgmonzon/Yangee/models"
)

func TestPostBody(w http.ResponseWriter, req *http.Request) {
  var test models.Test

  // Decode del JSON
  // ***************
  decoder := json.NewDecoder(req.Body)
  err := decoder.Decode(&test)
  if err != nil {
    core.RspMsgJSON(w, req, "ERROR", "JSON", "INVALID_PARAMS: JSON decode erróneo", http.StatusBadRequest)
    return
  }

  // Está todo Ok
  // ************
  test.Metodo = "TestPostBody"
  respuesta, _ := json.Marshal(test)
  core.RspJSON(w, req, respuesta, http.StatusOK)
  return
}

func TestOptionsBody(w http.ResponseWriter, req *http.Request) {
  var test models.Test

  // Decode del JSON
  // ***************
  decoder := json.NewDecoder(req.Body)
  err := decoder.Decode(&test)
  if err != nil {
    core.RspMsgJSON(w, req, "ERROR", "JSON", "INVALID_PARAMS: JSON decode erróneo", http.StatusBadRequest)
    return
  }

  // Está todo Ok
  // ************
  test.Metodo = "TestOptionsBody"
  respuesta, _ := json.Marshal(test)
  core.RspJSON(w, req, respuesta, http.StatusOK)
  return
}

func TestPostHeader(w http.ResponseWriter, req *http.Request) {
  var api_clienteID models.API_Cliente

  // Decode del JSON
  // ***************
  decoder := json.NewDecoder(req.Body)
  err := decoder.Decode(&api_clienteID)
  if err != nil {
    core.RspMsgJSON(w, req, "ERROR", "JSON", "INVALID_PARAMS: JSON decode erróneo", http.StatusBadRequest)
    return
  }

  var test models.Test
  test.Authorization = req.Header.Get("Authorization")
  test.API_ClienteID = api_clienteID.API_ClienteID

  // Está todo Ok
  // ************
  test.Metodo = "TestPostHeader"
  respuesta, _ := json.Marshal(test)
  core.RspJSON(w, req, respuesta, http.StatusOK)
  return
}

func TestOptionsHeader(w http.ResponseWriter, req *http.Request) {
  var api_clienteID models.API_Cliente

  // Decode del JSON
  // ***************
  decoder := json.NewDecoder(req.Body)
  err := decoder.Decode(&api_clienteID)
  if err != nil {
    core.RspMsgJSON(w, req, "ERROR", "JSON", "INVALID_PARAMS: JSON decode erróneo", http.StatusBadRequest)
    return
  }

  var test models.Test
  test.Authorization = req.Header.Get("Authorization")
  test.API_ClienteID = api_clienteID.API_ClienteID

  // Está todo Ok
  // ************
  test.Metodo = "TestOptionsHeader"
  respuesta, _ := json.Marshal(test)
  core.RspJSON(w, req, respuesta, http.StatusOK)
  return
}

func Autorizar(w http.ResponseWriter, req *http.Request) {
  var api_clienteID models.API_Cliente

  // Decode del JSON
  // ***************
  decoder := json.NewDecoder(req.Body)
  err := decoder.Decode(&api_clienteID)
  if err != nil {
    core.RspMsgJSON(w, req, "ERROR", "JSON", "INVALID_PARAMS: JSON decode erróneo", http.StatusBadRequest)
    return
  }

  // Valido el token del clienteAPI
  // ******************************
  estado, valor, mensaje, httpStat, aut, usuario, empresa := ValidarTokenCliente(w, req, api_clienteID)
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
