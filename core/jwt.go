package core

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
func ObtenerToken(authorizationHeader string) (models.Token, error, int) {
  var token models.Token
  aut, err, httpStat := ValidarAutorizacion(authorizationHeader)
  if err != nil {
    return token, err, httpStat
  } else {
    token, err, httpStat := GenerarToken(aut)
    return token, err, httpStat
  }
}

// Valida el token de autorización y devuelve el usuario
func ValidarAutorizacion(authorizationHeader string) (models.AutorizarToken, error, int) {
  var aut models.AutorizarToken
  if authorizationHeader != "" {
    bearerToken := strings.Split(authorizationHeader, " ")
    if len(bearerToken) == 2 {
      token, error := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
          if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
              return nil, fmt.Errorf("INVALID_PARAMS_01")
          }
          return []byte(config.SecretKey), nil
      })
      if error != nil {
        return aut, error, http.StatusBadRequest
      }
      if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        mapstructure.Decode(claims, &aut)
        return aut, nil, http.StatusOK
      } else {
        return aut, fmt.Errorf("INVALID_PARAMS_03"), http.StatusBadRequest
      }
    } else {
      return aut, fmt.Errorf("INVALID_PARAMS_04"), http.StatusBadRequest
    }
  } else {
    return aut, fmt.Errorf("INVALID_PARAMS_05"), http.StatusBadRequest
  }
}

// Genera el token para el usuario autorizado
func GenerarToken(aut models.AutorizarToken) (models.Token, error, int) {
  var token models.Token

  jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
      "usr": aut.User,
      "iat": time.Now(),
      "exp": time.Now().Add(time.Minute * config.ExpiraToken),
  })
  tokenString, error := jwtToken.SignedString([]byte(config.SecretKey))
  if error != nil {
    return token, fmt.Errorf("INTERNAL_SERVER_ERROR_01"), http.StatusInternalServerError
  }
  token.Token = tokenString
  return token, nil, http.StatusOK
}

func CrearToken(aut models.Autorizar) (string, error) {
  token := jwt.New(jwt.SigningMethodRS256)

  token.Claims = &models.TokenClaims{
    &jwt.StandardClaims{
      ExpiresAt: time.Now().Add(time.Minute * config.ExpiraToken).Unix(),
    },
    aut.Usuario,
    aut.Clave,
  }

  tokenString, err := token.SignedString(config.SignKey)
  if err != nil {
    return "Error firmando el token", err
  }

  return tokenString, nil
}
