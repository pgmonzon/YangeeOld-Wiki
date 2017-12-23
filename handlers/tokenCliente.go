package handlers

import (
  "net/http"
  "encoding/json"
  "time"
  "fmt"
  "strings"

  "github.com/pgmonzon/Yangee/core"
  "github.com/pgmonzon/Yangee/models"

  "github.com/dgrijalva/jwt-go"
)

// Genera el token que debe generar el cliente, es a los efectos de testing
func TokenCliente(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
  var reqCliente models.ReqCliente

  decoder := json.NewDecoder(req.Body)
  err := decoder.Decode(&reqCliente)
  if err != nil {
    s := []string{"INVALID_PARAMS:", err.Error()}
    core.RespErrorJSON(w, req, start, fmt.Errorf(strings.Join(s, " ")), http.StatusBadRequest)
  } else {
    if reqCliente.Usuario == "" || reqCliente.Clave == "" || reqCliente.Audience == "" {
      s := []string{"INVALID_PARAMS:", "Usuario, clave y audience no pueden estar vacíos"}
      core.RespErrorJSON(w, req, start, fmt.Errorf(strings.Join(s, " ")), http.StatusBadRequest)
    } else {
      var clienteAPI models.ClienteAPI

      clienteAPI, err, httpStat := ClienteAPITraer(reqCliente.Audience)
      if err != nil {
        core.RespErrorJSON(w, req, start, err, httpStat)
        return
      }

      claveEncriptada, err, httpStat := core.Encriptar(clienteAPI.Aes, reqCliente.Clave)
      if err != nil {
        core.RespErrorJSON(w, req, start, err, httpStat)
        return
      }
      token := jwt.New(jwt.SigningMethodHS256)
    	claims := make(jwt.MapClaims)
      claims["usr"] = reqCliente.Usuario
      claims["pas"] = claveEncriptada
    	claims["iat"] = time.Now().Unix()
      claims["aud"] = reqCliente.Audience
    	token.Claims = claims

      tokenString, err := token.SignedString([]byte(clienteAPI.Firma))
      if err != nil {
        s := []string{"INTERNAL_SERVER_ERROR:", err.Error()}
        core.RespErrorJSON(w, req, start, fmt.Errorf(strings.Join(s, " ")), http.StatusInternalServerError)
      } else {
        var tokenCliente models.TokenCliente
        tokenCliente.Token = tokenString
        respuesta, err := json.Marshal(tokenCliente)
        if err != nil {
          s := []string{"INTERNAL_SERVER_ERROR:", err.Error()}
          core.RespErrorJSON(w, req, start, fmt.Errorf(strings.Join(s, " ")), http.StatusInternalServerError)
        } else {
          core.RespuestaJSON(w, req, start, respuesta, http.StatusOK)
        }
      }
    }
  }
}
