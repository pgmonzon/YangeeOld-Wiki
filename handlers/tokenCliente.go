package handlers

import (
/**
  "net/http"
  "encoding/json"
  "time"
  "strings"

  "github.com/pgmonzon/Yangee/core"
  "github.com/pgmonzon/Yangee/models"

  "github.com/dgrijalva/jwt-go"
**/
)
/**
// Es el token que debe generar el cliente, es a los efectos de ejemplo
func TokenCliente(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
  var reqCliente models.ReqCliente
  var rsp models.Respuesta

  decoder := json.NewDecoder(req.Body)
  err := decoder.Decode(&reqCliente)
  if err != nil {
    rsp.Estado = "ERROR"
    rsp.Valor = "JSON decode"
    s := []string{"INVALID_PARAMS: ", err.Error()}
    rsp.Mensaje = strings.Join(s, "")
    respuesta, error := json.Marshal(rsp)
    core.FatalErr(error)
    core.RespuestaJSON(w, req, start, respuesta, http.StatusBadRequest)
    return
  } else {
    if reqCliente.Usuario == "" || reqCliente.Clave == "" || reqCliente.Audience == "" {
      rsp.Estado = "ERROR"
      rsp.Valor = "Campos obligatorios en vacío"
      s := []string{"INVALID_PARAMS: ", "Usuario, clave y audience no pueden estar vacíos"}
      rsp.Mensaje = strings.Join(s, "")
      respuesta, error := json.Marshal(rsp)
      core.FatalErr(error)
      core.RespuestaJSON(w, req, start, respuesta, http.StatusBadRequest)
      return
    } else {
      var clienteAPI models.ClienteAPI

      clienteAPI, err, httpStat := ClienteAPITraer(reqCliente.Audience)
      if err != nil {
        rsp.Estado = "ERROR"
        rsp.Valor = "ClienteAPITraer"
        s := []string{err.Error()}
        rsp.Mensaje = strings.Join(s, "")
        respuesta, error := json.Marshal(rsp)
        core.FatalErr(error)
        core.RespuestaJSON(w, req, start, respuesta, httpStat)
        return
      }

      if clienteAPI.ClienteAPI == "" {
        rsp.Estado = "ERROR"
        rsp.Valor = "ClienteAPI"
        s := []string{"INVALID_PARAMS: ", "ClienteAPI inexistente"}
        rsp.Mensaje = strings.Join(s, "")
        respuesta, error := json.Marshal(rsp)
        core.FatalErr(error)
        core.RespuestaJSON(w, req, start, respuesta, http.StatusBadRequest)
        return
      }

      claveEncriptada, err, httpStat := core.Encriptar(clienteAPI.Aes, reqCliente.Clave)
      if err != nil {
        rsp.Estado = "ERROR"
        rsp.Valor = "Encriptar"
        s := []string{err.Error()}
        rsp.Mensaje = strings.Join(s, "")
        respuesta, error := json.Marshal(rsp)
        core.FatalErr(error)
        core.RespuestaJSON(w, req, start, respuesta, httpStat)
        return
      }
      token := jwt.New(jwt.SigningMethodHS256)
    	claims := make(jwt.MapClaims)
      claims["usr"] = reqCliente.Usuario
      claims["pas"] = claveEncriptada
    	claims["iat"] = time.Now().Unix()
    	token.Claims = claims

      tokenString, err := token.SignedString([]byte(clienteAPI.Firma))
      if err != nil {
        rsp.Estado = "ERROR"
        rsp.Valor = "SignedString"
        s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
        rsp.Mensaje = strings.Join(s, "")
        respuesta, error := json.Marshal(rsp)
        core.FatalErr(error)
        core.RespuestaJSON(w, req, start, respuesta, http.StatusInternalServerError)
        return
      } else {
        var tokenCliente models.TokenCliente
        tokenCliente.Token = tokenString
        respuesta, err := json.Marshal(tokenCliente)
        if err != nil {
          rsp.Estado = "ERROR"
          rsp.Valor = "Marshal"
          s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
          rsp.Mensaje = strings.Join(s, "")
          respuesta, error := json.Marshal(rsp)
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
}
**/
