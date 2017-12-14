package handlers

import (
  "net/http"
  "encoding/json"
  "time"
  "strings"
  "fmt"

  "github.com/pgmonzon/Yangee/models"
  "github.com/pgmonzon/Yangee/core"
  "github.com/pgmonzon/Yangee/config"

  "github.com/dgrijalva/jwt-go"
  "github.com/mitchellh/mapstructure"
)

func Autorizar(w http.ResponseWriter, req *http.Request) {
  start := time.Now()

  token, err, httpStat := core.ObtenerToken(req.Header.Get("authorization"))
  if err != nil {
    var resp models.Error
    resp.Estado = "ERROR"
    resp.Detalle = err.Error()
    respuesta, error := json.Marshal(resp)
    core.FatalErr(error)
    core.RespuestaJSON(w, req, start, respuesta, httpStat)
  } else {
    respuesta, error := json.Marshal(token)
    core.FatalErr(error)
    core.RespuestaJSON(w, req, start, respuesta, httpStat)
  }
/**
  var aut models.AutorizarToken

  aut, err, httpStat := core.ValidarAutorizacion(req.Header.Get("authorization"))
  if err != nil {
    // No se pudo validar el token de autorización
    var resp models.Error
    resp.Estado = "ERROR"
    resp.Detalle = err.Error()
    respuesta, error := json.Marshal(resp)
    core.FatalErr(error)
    core.RespuestaJSON(w, req, start, respuesta, httpStat)
  } else {
    token, err, httpStat := core.GenerarToken(aut)
    if err != nil {
      var resp models.Error
      resp.Estado = "ERROR"
      resp.Detalle = err.Error()
      respuesta, error := json.Marshal(resp)
      core.FatalErr(error)
      core.RespuestaJSON(w, req, start, respuesta, httpStat)
    }
    respuesta, error := json.Marshal(token)
    core.FatalErr(error)
    core.RespuestaJSON(w, req, start, respuesta, httpStat)
  }
**/
}

// Valida el token con las credenciales y devuelve el token para operar
func Autorizar_bis(w http.ResponseWriter, req *http.Request) {
  start := time.Now()

  authorizationHeader := req.Header.Get("authorization")
  if authorizationHeader != "" {
    bearerToken := strings.Split(authorizationHeader, " ")
    if len(bearerToken) == 2 {
      token, error := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
          if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
              return nil, fmt.Errorf("There was an error")
          }
          return []byte(config.SecretKey), nil
      })
      if error != nil {
        var resp models.Error
        resp.Estado = "ERROR"
        resp.Detalle = "INVALID_PARAMS_04"
        respuesta, error := json.Marshal(resp)
        core.FatalErr(error)
        core.RespuestaJSON(w, req, start, respuesta, http.StatusBadRequest)
      }
      if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        var usuario models.AutorizarToken
        mapstructure.Decode(claims, &usuario)
        token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
            "usr": usuario.User,
            "iat": time.Now(),
            "exp": time.Now().Add(time.Minute * config.ExpiraToken),
        })
        tokenString, error := token.SignedString([]byte(config.SecretKey))
        if error != nil {
          var resp models.Error
          resp.Estado = "ERROR"
          resp.Detalle = "INTERNAL_SERVER_ERROR_01"
          respuesta, error := json.Marshal(resp)
          core.FatalErr(error)
          core.RespuestaJSON(w, req, start, respuesta, http.StatusInternalServerError)
        }
        var resp models.Token
        resp.Token = tokenString
        respuesta, error := json.Marshal(resp)
        core.FatalErr(error)
        core.RespuestaJSON(w, req, start, respuesta, http.StatusOK)
      } else {
        var resp models.Error
        resp.Estado = "ERROR"
        resp.Detalle = "INVALID_PARAMS_02"
        respuesta, error := json.Marshal(resp)
        core.FatalErr(error)
        core.RespuestaJSON(w, req, start, respuesta, http.StatusBadRequest)
      }
    } else {
      var resp models.Error
      resp.Estado = "ERROR"
      resp.Detalle = "INVALID_PARAMS_03"
      respuesta, error := json.Marshal(resp)
      core.FatalErr(error)
      core.RespuestaJSON(w, req, start, respuesta, http.StatusBadRequest)
    }
  } else {
    var resp models.Error
    resp.Estado = "ERROR"
    resp.Detalle = "INVALID_PARAMS_01"
    respuesta, error := json.Marshal(resp)
    core.FatalErr(error)
    core.RespuestaJSON(w, req, start, respuesta, http.StatusBadRequest)
  }
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
