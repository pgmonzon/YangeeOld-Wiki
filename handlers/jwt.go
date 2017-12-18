package handlers

import (
  "time"
  "strings"
  "fmt"
  "net/http"

  "github.com/pgmonzon/Yangee/models"
  "github.com/pgmonzon/Yangee/config"

  "github.com/dgrijalva/jwt-go"
  "github.com/mitchellh/mapstructure"
)

// Valida el token de autorizacion y devuelve el token para operar
func ObtenerToken(authorizationHeader string, clienteAPIHeader string) (models.Token, error, int) {
  var token models.Token
  aut, err, httpStat := ValidarAutorizacion(authorizationHeader, clienteAPIHeader)
  if err != nil {
    return token, err, httpStat
  } else {
    token, err, httpStat := GenerarToken(aut)
    return token, err, httpStat
  }
}

// Valida el token de autorización y devuelve el usuario
func ValidarAutorizacion(authorizationHeader string, clienteAPIHeader string) (models.AutorizarToken, error, int) {
  var aut models.AutorizarToken
  firma, err, httpStat := ClienteAPITraerFirma(clienteAPIHeader)
  if err != nil {
    return aut, err, httpStat
  }

  if authorizationHeader != "" {
    bearerToken := strings.Split(authorizationHeader, " ")
    if len(bearerToken) == 2 {
      token, error := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
          if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
              return nil, fmt.Errorf("INVALID_PARAMS: El token no es válido")
          }
          return []byte(firma), nil
      })
      if error != nil {
        s := []string{"INVALID_PARAMS:", error.Error()}
        return aut, fmt.Errorf(strings.Join(s, " ")), http.StatusBadRequest
      }
      if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        mapstructure.Decode(claims, &aut)
        // valida la fecha de creación del token de autorización
        if aut.Iat >= time.Now().Add(-time.Minute * config.ExpiraTokenAut).Unix() && aut.Iat <= time.Now().Add(time.Minute * config.ExpiraTokenAut).Unix() {
          return aut, nil, http.StatusOK
        } else {
          return aut, fmt.Errorf("INVALID_PARAMS: La fecha del token no es válida"), http.StatusBadRequest
        }
      } else {
        return aut, fmt.Errorf("INVALID_PARAMS: El token no es válido"), http.StatusBadRequest
      }
    } else {
      return aut, fmt.Errorf("INVALID_PARAMS: La key Authorization no tiene el prefijo Bearer  y un espacio antes del token"), http.StatusBadRequest
    }
  } else {
    return aut, fmt.Errorf("INVALID_PARAMS: Está vacía la key Authorization en el header"), http.StatusBadRequest
  }
}

// Genera el token para el usuario autorizado
func GenerarToken(aut models.AutorizarToken) (models.Token, error, int) {
  var token models.Token

  jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
      "usr": aut.User,
      "iat": time.Now().Unix(),
      "exp": time.Now().Add(time.Minute * config.ExpiraToken).Unix(),
  })

  tokenString, error := jwtToken.SignedString([]byte(config.SecretKey))
  if error != nil {
    return token, fmt.Errorf("INTERNAL_SERVER_ERROR: No pudimos firmar el token"), http.StatusInternalServerError
  }
  token.Token = tokenString
  return token, nil, http.StatusOK
}
