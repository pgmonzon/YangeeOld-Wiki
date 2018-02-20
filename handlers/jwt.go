package handlers

import (
  "time"
  "strings"
  "fmt"
  "net/http"
  "encoding/json"

  "github.com/pgmonzon/Yangee/models"
  "github.com/pgmonzon/Yangee/config"
  "github.com/pgmonzon/Yangee/core"

  "github.com/dgrijalva/jwt-go"
  "github.com/dgrijalva/jwt-go/request"
  "github.com/mitchellh/mapstructure"
)

// Valida el token generado por el cliente API
func ValidarTokenCliente(w http.ResponseWriter, req *http.Request) (models.AutorizarTokenCliente, error, int) {
  var aut models.AutorizarTokenCliente
  var clienteAPI models.ClienteAPI

  // Busco la firma y la clave de encriptación que usa el cliente API
  clienteAPI, err, httpStat := ClienteAPITraer(req.Header.Get("API_ClienteID"))
  if err != nil {
    return aut, err, httpStat
  }

  token, err := request.ParseFromRequestWithClaims(req, request.AuthorizationHeaderExtractor, &models.AutorizarTokenCliente{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(clienteAPI.Firma), nil
		})

	if err == nil {
		if token.Valid {
      claims := token.Claims.(*models.AutorizarTokenCliente)
      mapstructure.Decode(claims, &aut)
      if claims.IssuedAt >= time.Now().Add(-time.Minute * config.ExpiraTokenAut).Unix() && claims.IssuedAt <= time.Now().Add(time.Minute * config.ExpiraTokenAut).Unix() {
        claveDesencriptada, err, httpStat := core.Desencriptar(clienteAPI.Aes, claims.Pas)
        if err != nil {
          return aut, err, httpStat
        }
        err, httpStat = UsuarioLogin(claims.Usr, claveDesencriptada)
        if err != nil {
          return aut, err, httpStat
        } else {
          return aut, nil, http.StatusOK
        }
      } else {
        s := []string{"INVALID_PARAMS:", "La fecha es inválida"}
        return aut, fmt.Errorf(strings.Join(s, " ")), http.StatusBadRequest
      }
		} else {
      s := []string{"INVALID_PARAMS:", err.Error()}
      return aut, fmt.Errorf(strings.Join(s, " ")), http.StatusBadRequest
		}
	} else {
    s := []string{"INVALID_PARAMS:", err.Error()}
    return aut, fmt.Errorf(strings.Join(s, " ")), http.StatusBadRequest
	}
}

// Genera el token para el usuario autorizado
func GenerarToken(aut models.AutorizarTokenCliente) (models.Token, error, int) {
  var tokenAutorizado models.Token

  permisos, err, httpStat := UsuarioPermisos(aut.Usr)
  if err != nil {
    return tokenAutorizado, err, httpStat
	}

  permisosEncriptados, err, httpStat := core.Encriptar(config.Aes, permisos)
  if err != nil {
    return tokenAutorizado, err, httpStat
  }

  token := jwt.New(jwt.SigningMethodRS256)
  claims := make(jwt.MapClaims)
  claims["usr"] = aut.Usr
  claims["exp"] = time.Now().Add(time.Minute * config.ExpiraToken).Unix()
	claims["iat"] = time.Now().Unix()
  claims["rbc"] = permisosEncriptados
  token.Claims = claims

	tokenString, err := token.SignedString(config.SignKey)
	if err != nil {
    s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
    return tokenAutorizado, fmt.Errorf(strings.Join(s, "")), http.StatusInternalServerError
	}

  tokenAutorizado.Token = tokenString
  return tokenAutorizado, nil, http.StatusOK
}

func ValidarMiddleware(next http.HandlerFunc, permiso string) http.HandlerFunc {
  return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
  start := time.Now()
  var resp models.Resp
  var mensaje models.Mensaje

  token, err := request.ParseFromRequestWithClaims(req, request.AuthorizationHeaderExtractor, &models.TokenAutorizado{},
		func(token *jwt.Token) (interface{}, error) {
			return config.VerifyKey, nil
		})

  if err == nil {
    if token.Valid {
      claims := token.Claims.(*models.TokenAutorizado)
      if claims.ExpiresAt >= time.Now().Unix() {
        // si no está expirado me fijo si tiene el permiso
        if permiso == "" {
          resp.EstadoGral = "ERROR"
          mensaje.Valor = "Permiso a validar"
          mensaje.Estado = "ERROR"
          s := []string{"INVALID_PARAMS: ", "No tiene permiso a validar"}
          mensaje.Mensaje = strings.Join(s, "")
          resp.Mensajes = append(resp.Mensajes, mensaje)
          respuesta, error := json.Marshal(resp)
          core.FatalErr(error)
          core.RespuestaJSON(w, req, start, respuesta, http.StatusBadRequest)
          return
        } else {
          cadena := []string{"#", permiso, "#"}
          permisoBuscar := strings.Join(cadena, "")
          permisosDesencriptados, err, httpStat := core.Desencriptar(config.Aes, claims.Rbc)
          if err != nil {
            resp.EstadoGral = "ERROR"
            mensaje.Valor = "Desencriptar"
            mensaje.Estado = "ERROR"
            s := []string{err.Error()}
            mensaje.Mensaje = strings.Join(s, "")
            resp.Mensajes = append(resp.Mensajes, mensaje)
            respuesta, error := json.Marshal(resp)
            core.FatalErr(error)
            core.RespuestaJSON(w, req, start, respuesta, httpStat)
            return
          }

          if strings.Contains(permisosDesencriptados, permisoBuscar) == false {
            resp.EstadoGral = "ERROR"
            mensaje.Valor = "Permiso denegado"
            mensaje.Estado = "ERROR"
            s := []string{"FORBIDDEN: ", "No tenés permiso"}
            mensaje.Mensaje = strings.Join(s, "")
            resp.Mensajes = append(resp.Mensajes, mensaje)
            respuesta, error := json.Marshal(resp)
            core.FatalErr(error)
            core.RespuestaJSON(w, req, start, respuesta, http.StatusForbidden)
            return
          }
        }

        //context.Set(req, "decoded", token.Claims)
        next(w, req)
      } else {
        resp.EstadoGral = "ERROR"
        mensaje.Valor = "Token Expirado"
        mensaje.Estado = "ERROR"
        s := []string{"FORBIDDEN: ", "Token Expirado"}
        mensaje.Mensaje = strings.Join(s, "")
        resp.Mensajes = append(resp.Mensajes, mensaje)
        respuesta, error := json.Marshal(resp)
        core.FatalErr(error)
        core.RespuestaJSON(w, req, start, respuesta, http.StatusForbidden)
        return
      }
    } else {
      resp.EstadoGral = "ERROR"
      mensaje.Valor = "Token Inválido"
      mensaje.Estado = "ERROR"
      s := []string{"FORBIDDEN: ", "Token Inválido"}
      mensaje.Mensaje = strings.Join(s, "")
      resp.Mensajes = append(resp.Mensajes, mensaje)
      respuesta, error := json.Marshal(resp)
      core.FatalErr(error)
      core.RespuestaJSON(w, req, start, respuesta, http.StatusForbidden)
      return
    }
  } else {
    resp.EstadoGral = "ERROR"
    mensaje.Valor = "ParseFromRequestWithClaims"
    mensaje.Estado = "ERROR"
    s := []string{"INVALID_PARAMS: ", err.Error()}
    mensaje.Mensaje = strings.Join(s, "")
    resp.Mensajes = append(resp.Mensajes, mensaje)
    respuesta, error := json.Marshal(resp)
    core.FatalErr(error)
    core.RespuestaJSON(w, req, start, respuesta, http.StatusBadRequest)
    return
  }
  })
}
