package handlers

import (
  "time"
  "strings"
  "fmt"
  "net/http"

  "github.com/pgmonzon/Yangee/models"
  "github.com/pgmonzon/Yangee/config"
  "github.com/pgmonzon/Yangee/core"

  "github.com/dgrijalva/jwt-go"
  "github.com/dgrijalva/jwt-go/request"
  "github.com/mitchellh/mapstructure"
  "github.com/gorilla/context"
  "gopkg.in/mgo.v2/bson"
)

// Valida el token generado por el cliente API
func ValidarTokenCliente(w http.ResponseWriter, req *http.Request, api_clienteID models.API_Cliente) (string, string, string, int, models.AutorizarTokenCliente, models.Usuario, models.Empresa) {
  var aut models.AutorizarTokenCliente
  var usuario models.Usuario
  var empresa models.Empresa
  var clienteAPI models.ClienteAPI

  // Busco la firma y la clave de encriptación que usa el cliente API
  // ****************************************************************
  clienteAPI, err, httpStat := ClienteAPI_X_clienteAPI(api_clienteID.API_ClienteID)
  if err != nil {
    return "ERROR", "ClienteAPI_X_clienteAPI", err.Error(), httpStat, aut, usuario, empresa
  }

  // Establezco las variables
  // ************************
  context.Set(req, "ClienteAPI_id", clienteAPI.ID)
  context.Set(req, "ClienteAPI", clienteAPI.ClienteAPI)

  // Parseo el request
  // *****************
  token, err := request.ParseFromRequestWithClaims(req, request.AuthorizationHeaderExtractor, &models.AutorizarTokenCliente{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(clienteAPI.Firma), nil
		})
  if err != nil {
    s := []string{"INVALID_PARAMS:", err.Error()}
    return "ERROR", "ParseFromRequestWithClaims", strings.Join(s, " "), http.StatusBadRequest, aut, usuario, empresa
  }

  // Me fijo si es válido el token
  // *****************************
  if !token.Valid {
    s := []string{"INVALID_PARAMS: Token inválido"}
    return "ERROR", "Token Valid", strings.Join(s, " "), http.StatusBadRequest, aut, usuario, empresa
  }

  // Hago el decode del token
  // ************************
  claims := token.Claims.(*models.AutorizarTokenCliente)
  mapstructure.Decode(claims, &aut)

  // Me fijo si expiró el token
  // **************************
  if claims.IssuedAt < time.Now().Add(-time.Minute * config.ExpiraTokenAut).Unix() {
    s := []string{"INVALID_PARAMS:", "Expiró el token"}
    return "ERROR", "Token Expiración", strings.Join(s, " "), http.StatusBadRequest, aut, usuario, empresa
  }

  // Intento desencriptar la clave
  // *****************************
  claveDesencriptada, err, httpStat := core.Desencriptar(clienteAPI.Aes, claims.Pas)
  if err != nil {
    return "ERROR", "Desencriptar clave", err.Error(), httpStat, aut, usuario, empresa
  }

  // Intento loguear el Usuario
  // **************************
  estado, valor, mensaje, httpStat, usuario, empresa := UsuarioLogin(claims.Usr, claveDesencriptada)
  if httpStat != http.StatusOK {
    return estado, valor, mensaje, httpStat, aut, usuario, empresa
  }

  // Establezco las variables
  // ************************
  context.Set(req, "Usuario_id", usuario.ID)
  context.Set(req, "Usuario", usuario.Usuario)
  context.Set(req, "Empresa_id", usuario.Empresa_id)
  context.Set(req, "Empresa", empresa.Empresa)

  // Está todo Ok
  // ************
  return "OK", "ValidarTokenCliente", "Ok", http.StatusOK, aut, usuario, empresa
}

// Genera el token para el usuario autorizado
func GenerarToken(aut models.AutorizarTokenCliente, req *http.Request) (models.Token, error, int) {
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
  claims["cid"] = context.Get(req, "ClienteAPI_id").(bson.ObjectId)
  claims["clt"] = context.Get(req, "ClienteAPI").(string)
  claims["uid"] = context.Get(req, "Usuario_id").(bson.ObjectId)
  claims["eid"] = context.Get(req, "Empresa_id").(bson.ObjectId)
  claims["emp"] = context.Get(req, "Empresa").(string)

  token.Claims = claims

	tokenString, err := token.SignedString(config.SignKey)
	if err != nil {
    s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
    return tokenAutorizado, fmt.Errorf(strings.Join(s, "")), http.StatusInternalServerError
	}

  tokenAutorizado.Token = tokenString
  return tokenAutorizado, nil, http.StatusOK
}

// Todas las operaciones del router pasan por aca
func ValidarMiddleware(next http.HandlerFunc, permiso string) http.HandlerFunc {
  return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
    start := time.Now()

    // Establezco las variables
    // ************************
    objID := bson.NewObjectId()
    context.Set(req, "CicloDeVida_id", objID)
    context.Set(req, "Start", start)
    context.Set(req, "ClienteAPI_id", objID)
    context.Set(req, "ClienteAPI", "")
    context.Set(req, "Usuario_id", objID)
    context.Set(req, "Usuario", "")
    context.Set(req, "Empresa_id", objID)
    context.Set(req, "Empresa", "")

    // Si es AUTH redirecciono
    // ***********************
    if permiso == "AUTH" {
      next(w, req)
      return
    }

    // Me fijo si está vacío el permiso del Handler
    // ********************************************
    if permiso == "" {
      s := []string{"INTERNAL_SERVER_ERROR: ", "Está vacío del permiso del Handler"}
      core.RspMsgJSON(w, req, "ERROR", "Permiso del Handler", strings.Join(s, ""), http.StatusInternalServerError)
      return
    }

    // Parseo el Request
    // *****************
    token, err := request.ParseFromRequestWithClaims(req, request.AuthorizationHeaderExtractor, &models.TokenAutorizado{}, func(token *jwt.Token) (interface{}, error) {
			return config.VerifyKey, nil
		})
    if err != nil {
      s := []string{"INVALID_PARAMS: ", err.Error()}
      core.RspMsgJSON(w, req, "ERROR", "ParseFromRequestWithClaims", strings.Join(s, ""), http.StatusBadRequest)
      return
    }

    // Me fijo si es válido el token
    // *****************************
    if !token.Valid {
      s := []string{"INVALID_PARAMS: ", "Token Inválido"}
      core.RspMsgJSON(w, req, "ERROR", "Token Inválido", strings.Join(s, ""), http.StatusBadRequest)
      return
    }

    // Obtengo el Claims
    // *****************
    claims := token.Claims.(*models.TokenAutorizado)

    // Establezco las variables
    // ************************
    context.Set(req, "ClienteAPI_id", claims.Cid)
    context.Set(req, "ClienteAPI", claims.Clt)
    context.Set(req, "Usuario_id", claims.Uid)
    context.Set(req, "Usuario", claims.Usr)
    context.Set(req, "Empresa_id", claims.Eid)
    context.Set(req, "Empresa", claims.Emp)

    // Me fijo si está Expirado
    // ************************
    if claims.ExpiresAt < time.Now().Unix() {
      s := []string{"FORBIDDEN: ", "Token Expirado"}
      core.RspMsgJSON(w, req, "ERROR", "Token Expirado", strings.Join(s, ""), http.StatusForbidden)
      return
    }

    // Si es NO_VALIDAR redirecciono
    // *****************************
    if permiso == "NO_VALIDAR" {
      next(w, req)
      return
    }

    // Desencripto los permisos
    // ************************
    permisosDesencriptados, err, httpStat := core.Desencriptar(config.Aes, claims.Rbc)
    if err != nil {
      s := []string{err.Error()}
      core.RspMsgJSON(w, req, "ERROR", "Desencriptar", strings.Join(s, ""), httpStat)
      return
    }

    // Me fijo si tiene el permiso
    // ***************************
    cadena := []string{"#", permiso, "#"}
    permisoBuscar := strings.Join(cadena, "")
    // Si no lo tiene
    if strings.Contains(permisosDesencriptados, permisoBuscar) == false {
      s := []string{"FORBIDDEN: ", "No tenés permiso"}
      core.RspMsgJSON(w, req, "ERROR", "Permiso denegado", strings.Join(s, ""), http.StatusForbidden)
      return
    }

    // Está todo Ok
    // ************
    next(w, req)
    return
  })
}
