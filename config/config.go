package config

import (
  "io/ioutil"
  "log"

  "crypto/rsa"
  "github.com/dgrijalva/jwt-go"
  "gopkg.in/mgo.v2/bson"
)

const(
  // Base de datos
  DB_Host = "localhost"
  //DB_Host = "mongodb://127.0.0.1:27017"
  //DB_Host = "mongodb://yng_user:laser@ds021326.mlab.com:21326/yangee"
  DB_Name = "yangee"
  DB_User = "yangee"
  DB_Pass = "1331"
  DB_Timeout = 10 // valor en minutos
  DB_Transaction = "transaction"

  DB_Usuario = "usuario" // tabla de usuarios
  DB_ClienteAPI = "clienteapi" // tabla de los clientes de la API
  DB_Permiso = "permiso" // tabla de permisos

  // jwt
  privKeyPath = "C:/Users/Patricio/Google Drive/proyectoYangee/codigoGo/src/github.com/pgmonzon/Yangee/config/keys/app.rsa"
  pubKeyPath = "C:/Users/Patricio/Google Drive/proyectoYangee/codigoGo/src/github.com/pgmonzon/Yangee/config/keys/app.rsa.pub"
  ExpiraToken   = 100000 // en minutos - Expiración del token para operar
  ExpiraTokenAut = 100000 // en minutos - Expiración del token de autorización
)

var (
	VerifyKey       *rsa.PublicKey
	SignKey         *rsa.PrivateKey
  UsuarioActivoID bson.ObjectId
)

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Inicializar() {
  signBytes, err := ioutil.ReadFile(privKeyPath)
	fatal(err)

	SignKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	fatal(err)

	verifyBytes, err := ioutil.ReadFile(pubKeyPath)
	fatal(err)

	VerifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	fatal(err)
}
