package handlers

import (
  "net/http"
  "encoding/json"
  "time"
  "strings"

  "github.com/pgmonzon/Yangee/core"
  "github.com/pgmonzon/Yangee/models"

  "github.com/dgrijalva/jwt-go"
  "github.com/gorilla/context"
)

// Es el token que debe generar el cliente, es a los efectos de ejemplo
func TokenCliente(w http.ResponseWriter, req *http.Request) {
  var reqCliente models.ReqCliente

  decoder := json.NewDecoder(req.Body)
  err := decoder.Decode(&reqCliente)
  if err != nil {
    s := []string{"INVALID_PARAMS: ", err.Error()}
    core.RspMsgJSON(w, req, "ERROR", "JSON decode", strings.Join(s, ""), http.StatusBadRequest)
    return
  } else {
    if reqCliente.Usuario == "" || reqCliente.Clave == "" || reqCliente.Audience == "" {
      s := []string{"INVALID_PARAMS: ", "Usuario, clave y audience no pueden estar vacíos"}
      core.RspMsgJSON(w, req, "ERROR", "Campos obligatorios en vacío", strings.Join(s, ""), http.StatusBadRequest)
      return
    } else {
      var clienteAPI models.ClienteAPI

      clienteAPI, err, httpStat := ClienteAPI_X_clienteAPI(reqCliente.Audience)
      if err != nil {
        s := []string{err.Error()}
        core.RspMsgJSON(w, req, "ERROR", "ClienteAPITraer", strings.Join(s, ""), httpStat)
        return
      }

      if clienteAPI.ClienteAPI == "" {
        s := []string{"INVALID_PARAMS: ", "ClienteAPI inexistente"}
        core.RspMsgJSON(w, req, "ERROR", "ClienteAPI", strings.Join(s, ""), http.StatusBadRequest)
        return
      }

      claveEncriptada, err, httpStat := core.Encriptar(clienteAPI.Aes, reqCliente.Clave)
      if err != nil {
        s := []string{err.Error()}
        core.RspMsgJSON(w, req, "ERROR", "Encriptar", strings.Join(s, ""), httpStat)
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
        s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
        core.RspMsgJSON(w, req, "ERROR", "SignedString", strings.Join(s, ""), http.StatusInternalServerError)
        return
      } else {
        var tokenCliente models.TokenCliente
        tokenCliente.Token = tokenString
        respuesta, err := json.Marshal(tokenCliente)
        if err != nil {
          s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
          core.RspMsgJSON(w, req, "ERROR", "Marshal", strings.Join(s, ""), http.StatusInternalServerError)
          return
        } else {
          // Establezco las variables
          // ************************
          context.Set(req, "TipoOper", "#TokenCliente#")
          s := []string{"Devolvió el token cliente"}
          context.Set(req, "Novedad", strings.Join(s, ""))

          // Está todo Ok
          // ************
          core.RspJSON(w, req, respuesta, http.StatusOK)
          return
        }
      }
    }
  }
}
