package core

import (
  "time"

  "github.com/pgmonzon/Yangee/models"
  "github.com/pgmonzon/Yangee/config"

  "github.com/dgrijalva/jwt-go"
)

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
