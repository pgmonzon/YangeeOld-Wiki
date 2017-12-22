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

func Autorizar(w http.ResponseWriter, req *http.Request) {
  start := time.Now()

  token, err, httpStat := ObtenerToken(req.Header.Get("authorization"), req.Header.Get("API_ClienteID"))
  if err != nil {
    core.RespErrorJSON(w, req, start, err, httpStat)
  } else {
    respuesta, error := json.Marshal(token)
    if error != nil {
      core.RespErrorJSON(w, req, start, error, httpStat)
    } else {
      core.RespuestaJSON(w, req, start, respuesta, httpStat)
    }
  }
  return
}

func TokenAutorizar(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	var tokenAutorizar models.TokenAutorizar

  decoder := json.NewDecoder(req.Body)
  err := decoder.Decode(&tokenAutorizar)
  if err != nil {
    core.RespErrorJSON(w, req, start, fmt.Errorf("INVALID_PARAMS: JSON decode erróneo"), http.StatusBadRequest)
  } else {
    // campos obligatorios
    if tokenAutorizar.Usuario == "" || tokenAutorizar.Clave == "" || tokenAutorizar.Audience == "" {
      core.RespErrorJSON(w, req, start, fmt.Errorf("INVALID_PARAMS: Usuario, clave y audience no pueden estar vacíos"), http.StatusBadRequest)
    } else {
      var token models.Token
      var clienteAPI models.ClienteAPI

      clienteAPI, err, httpStat := ClienteAPITraer(tokenAutorizar.Audience)
      if err != nil {
        core.RespErrorJSON(w, req, start, err, httpStat)
        return
      }

      claveEncriptada, err, httpStat := core.Encriptar(clienteAPI.Aes, tokenAutorizar.Clave)
      if err != nil {
        core.RespErrorJSON(w, req, start, err, httpStat)
        return
      }

      // Create the Claims
      claims := models.TokenAutorizarClaims{
          tokenAutorizar.Usuario,
          claveEncriptada,
          jwt.StandardClaims{
              IssuedAt: time.Now().Unix(),
              Audience: tokenAutorizar.Audience,
          },
      }
      jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
      tokenString, error := jwtToken.SignedString([]byte(clienteAPI.Firma))
      if error != nil {
        core.RespErrorJSON(w, req, start, fmt.Errorf("INTERNAL_SERVER_ERROR: No pudimos firmar el token"), http.StatusInternalServerError)
      } else {
        token.Token = tokenString
        respuesta, error := json.Marshal(token)
        if error != nil {
          s := []string{"INTERNAL_SERVER_ERROR:", error.Error()}
          core.RespErrorJSON(w, req, start, fmt.Errorf(strings.Join(s, " ")), http.StatusInternalServerError)
        } else {
          core.RespuestaJSON(w, req, start, respuesta, http.StatusOK)
        }
      }
    }
  }
  return
}
