package handlers

import (
  "net/http"
  "encoding/json"
  "time"
  "strings"

  "github.com/pgmonzon/Yangee/core"
  "github.com/pgmonzon/Yangee/models"

  "github.com/dgrijalva/jwt-go"
)

// Es el token que debe generar el cliente, es a los efectos de ejemplo
func TokenCliente(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
  var reqCliente models.ReqCliente
  var resp models.Resp
  var mensaje models.Mensaje

  decoder := json.NewDecoder(req.Body)
  err := decoder.Decode(&reqCliente)
  if err != nil {
    resp.EstadoGral = "ERROR"
    mensaje.Valor = "JSON decode"
    mensaje.Estado = "ERROR"
    s := []string{"INVALID_PARAMS: ", err.Error()}
    mensaje.Mensaje = strings.Join(s, "")
    resp.Mensajes = append(resp.Mensajes, mensaje)
    respuesta, error := json.Marshal(resp)
    core.FatalErr(error)
    core.RespuestaJSON(w, req, start, respuesta, http.StatusBadRequest)
    return
  } else {
    if reqCliente.Usuario == "" || reqCliente.Clave == "" || reqCliente.Audience == "" {
      resp.EstadoGral = "ERROR"
      mensaje.Valor = "Campos obligatorios en vacío"
      mensaje.Estado = "ERROR"
      s := []string{"INVALID_PARAMS: ", "Usuario, clave y audience no pueden estar vacíos"}
      mensaje.Mensaje = strings.Join(s, "")
      resp.Mensajes = append(resp.Mensajes, mensaje)
      respuesta, error := json.Marshal(resp)
      core.FatalErr(error)
      core.RespuestaJSON(w, req, start, respuesta, http.StatusBadRequest)
      return
    } else {
      var clienteAPI models.ClienteAPI

      clienteAPI, err, httpStat := ClienteAPITraer(reqCliente.Audience)
      if err != nil {
        resp.EstadoGral = "ERROR"
        mensaje.Valor = "ClienteAPITraer"
        mensaje.Estado = "ERROR"
        s := []string{err.Error()}
        mensaje.Mensaje = strings.Join(s, "")
        resp.Mensajes = append(resp.Mensajes, mensaje)
        respuesta, error := json.Marshal(resp)
        core.FatalErr(error)
        core.RespuestaJSON(w, req, start, respuesta, httpStat)
        return
      }

      if clienteAPI.ClienteAPI == "" {
        resp.EstadoGral = "ERROR"
        mensaje.Valor = "ClienteAPI"
        mensaje.Estado = "ERROR"
        s := []string{"INVALID_PARAMS: ", "ClienteAPI inexistente"}
        mensaje.Mensaje = strings.Join(s, "")
        resp.Mensajes = append(resp.Mensajes, mensaje)
        respuesta, error := json.Marshal(resp)
        core.FatalErr(error)
        core.RespuestaJSON(w, req, start, respuesta, http.StatusBadRequest)
        return
      }

      claveEncriptada, err, httpStat := core.Encriptar(clienteAPI.Aes, reqCliente.Clave)
      if err != nil {
        resp.EstadoGral = "ERROR"
        mensaje.Valor = "Encriptar"
        mensaje.Estado = "ERROR"
        s := []string{err.Error()}
        mensaje.Mensaje = strings.Join(s, "")
        resp.Mensajes = append(resp.Mensajes, mensaje)
        respuesta, error := json.Marshal(resp)
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
        resp.EstadoGral = "ERROR"
        mensaje.Valor = "SignedString"
        mensaje.Estado = "ERROR"
        s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
        mensaje.Mensaje = strings.Join(s, "")
        resp.Mensajes = append(resp.Mensajes, mensaje)
        respuesta, error := json.Marshal(resp)
        core.FatalErr(error)
        core.RespuestaJSON(w, req, start, respuesta, http.StatusInternalServerError)
        return
      } else {
        var tokenCliente models.TokenCliente
        tokenCliente.Token = tokenString
        respuesta, err := json.Marshal(tokenCliente)
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
}
