package handlers

import (
  "net/http"
  "encoding/json"
  "time"
  "strings"
  "fmt"

  "github.com/pgmonzon/Yangee/models"
  "github.com/pgmonzon/Yangee/core"
  //"github.com/pgmonzon/Yangee/config"

  "github.com/dgrijalva/jwt-go"
  //"github.com/dgrijalva/jwt-go/request"
  "github.com/mitchellh/mapstructure"
)

// Valida el token con las credenciales y devuelve el token para operar
func Autorizar(w http.ResponseWriter, req *http.Request) {
  start := time.Now()
  var err models.Error
  var httpStat int

  authorizationHeader := req.Header.Get("authorization")
  if authorizationHeader != "" {
    bearerToken := strings.Split(authorizationHeader, " ")
    if len(bearerToken) == 2 {
      token, error := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
          if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
              return nil, fmt.Errorf("There was an error")
          }
          return []byte("secret"), nil
      })
      if error != nil {
        err.Estado = "ERROR"
        err.Detalle = "INVALID_PARAMS_04"
        httpStat = http.StatusBadRequest
      }
      if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        var usuario models.AutorizarToken
        mapstructure.Decode(claims, &usuario)
        err.Estado = "OK"
        err.Detalle = usuario.User
        httpStat = http.StatusOK
      } else {
        err.Estado = "ERROR"
        err.Detalle = "INVALID_PARAMS_02"
        httpStat = http.StatusBadRequest
      }
    } else {
      err.Estado = "ERROR"
      err.Detalle = "INVALID_PARAMS_03"
      httpStat = http.StatusBadRequest
    }
  } else {
    err.Estado = "ERROR"
    err.Detalle = "INVALID_PARAMS_01"
    httpStat = http.StatusBadRequest
  }

  respuesta, error := json.Marshal(err)
  core.FatalErr(error)
  core.RespuestaJSON(w, req, start, respuesta, httpStat)
  return
/**
///////////////////////////////////////////////////////////
  if authorizationHeader != "" {
      bearerToken := strings.Split(authorizationHeader, " ")
      if len(bearerToken) == 2 {
          token, error := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
              if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                  return nil, fmt.Errorf("There was an error")
              }
              return []byte("secret"), nil
          })
          if error != nil {
              json.NewEncoder(w).Encode(Exception{Message: error.Error()})
              return
          }
          if token.Valid {
              context.Set(req, "decoded", token.Claims)
              next(w, req)
          } else {
              json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
          }
      }
  } else {
      json.NewEncoder(w).Encode(Exception{Message: "An authorization header is required"})
  }
  //////////////////////////////////////////////
  start := time.Now()
  var err models.Error
  var httpStat int

  token, error := request.ParseFromRequestWithClaims(req, request.AuthorizationHeaderExtractor, &models.AutorizarToken{}, func(token *jwt.Token) (interface{}, error) {
    return config.SecretKey, nil
  })

  if error == nil {
    err.Estado = "ERROR"
    err.Detalle = "INVALID_PARAMS_01"
    httpStat = http.StatusBadRequest
  }

  if token.Valid{
    err.Estado = "OK"
    err.Detalle = token.Claims.(*models.AutorizarToken).Pas
    httpStat = http.StatusOK
  } else {
    err.Estado = "ERROR"
    err.Detalle = "INVALID_PARAMS_02"
    httpStat = http.StatusBadRequest
  }
  respuesta, error := json.Marshal(err)
  core.FatalErr(error)
  core.RespuestaJSON(w, req, start, respuesta, httpStat)
  return
  ///////////////////////////////////////////////////////////////////
**/
}

// Valida las credenciales del usuario y devuelve un Token
func Autorizar_Anterior(w http.ResponseWriter, req *http.Request) {
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
