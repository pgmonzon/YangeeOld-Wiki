package config

import (
  "io/ioutil"
  "log"

  "crypto/rsa"
  "github.com/dgrijalva/jwt-go"
  "gopkg.in/mgo.v2/bson"
)

const(
  // Consola
  MostarEnConsola = true
  RegistrarCicloDeVida = true

  // Base de datos
  DB_Host = "localhost"
  //DB_Host = "mongodb://127.0.0.1:27017"
  //DB_Host = "mongodb://yng_user:laser@ds021326.mlab.com:21326/yangee"
  DB_Name = "yangee"
  DB_User = "yangee"
  DB_Pass = "1331"
  DB_Timeout = 10 // valor en minutos
  DB_Transaction = "transaction"

  DB_CicloVida = "cicloVida" // ciclo de vida
  DB_Audit = "audit" // auditoría
  DB_Modulo = "modulos" // módulos del sistema
  DB_Empresa = "empresas" // empresa del usuario
  DB_Usuario = "usuarios" // tabla de usuarios
  DB_ClienteAPI = "clientesApi" // tabla de los clientes de la API
  DB_Permiso = "permisos" // tabla de permisos
  DB_Rol = "roles" // tabla de roles
  DB_Filosofo = "filosofos"
  DB_TipoUnidad = "tipoUnidad"

  // jwt
  privKeyPath = "C:/Users/Patricio/Google Drive/proyectoYangee/codigoGo/src/github.com/pgmonzon/Yangee/config/keys/app.rsa"
  pubKeyPath = "C:/Users/Patricio/Google Drive/proyectoYangee/codigoGo/src/github.com/pgmonzon/Yangee/config/keys/app.rsa.pub"
  //privKeyPath = "/home/pgmonzon/work/src/github.com/pgmonzon/Yangee/config/keys/app.rsa"
  //pubKeyPath = "/home/pgmonzon/work/src/github.com/pgmonzon/Yangee/config/keys/app.rsa.pub"
  ExpiraToken   = 100000 // en minutos - Expiración del token para operar
  ExpiraTokenAut = 100000 // en minutos - Expiración del token de autorización
  Aes = "AES256Key-32Characters1234567890"

  // Timestamp layout
  TimestampLayout = "2006-01-02T15:04:05.000Z"
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
