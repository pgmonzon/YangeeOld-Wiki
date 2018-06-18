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
  Produccion = false

  // Ambiente producción
  // *******************
  //DB_Host = "198.100.45.12:27017"
  DB_Host = "localhost"
  DB_Name = "yangee"
  DB_User = "yngee"
  DB_Pass = "1962Laser"
  DB_Timeout = 10 // valor en minutos
  DB_Transaction = "transacciones"
  privKeyPath = "/usr/local/go/src/github.com/pgmonzon/Yangee/config/keys/app.rsa"
  pubKeyPath = "/usr/local/go/src/github.com/pgmonzon/Yangee/config/keys/app.rsa.pub"
/*
  // Ambiente Desarrollo
  // *******************
  DB_Host = "localhost"
  DB_Name = "yangee"
  DB_User = ""
  DB_Pass = ""
  DB_Timeout = 10 // valor en minutos
  DB_Transaction = "transacciones"
  privKeyPath = "C:/Users/Patricio/Google Drive/proyectoYangee/codigoGo/src/github.com/pgmonzon/Yangee/config/keys/app.rsa"
  pubKeyPath = "C:/Users/Patricio/Google Drive/proyectoYangee/codigoGo/src/github.com/pgmonzon/Yangee/config/keys/app.rsa.pub"
*/
  DB_CicloVida = "cicloVida" // ciclo de vida
  DB_Audit = "audit" // auditoría
  DB_Modulo = "modulos" // módulos del sistema
  DB_Empresa = "empresas" // empresa del usuario
  DB_Usuario = "usuarios" // tabla de usuarios
  DB_ClienteAPI = "clientesApi" // tabla de los clientes de la API
  DB_Permiso = "permisos" // tabla de permisos
  DB_Rol = "roles" // tabla de roles
  DB_Filosofo = "filosofos"
  DB_TipoUnidad = "tipoUnidades"
  DB_Categoria = "categorias"
  DB_CuentaGasto = "cuentaGastos"
  DB_BasicoSindicato = "basicoSindicatos"
  DB_Unidad = "unidades"
  DB_Personal = "personal"
  DB_Locacion = "locaciones"
  DB_Cliente = "clientes"
  DB_Transportista = "transportistas"
  DB_Viaje = "viajes"
  DB_Autorizacion = "autorizaciones"
  DB_Factura = "facturas"
  DB_Liquidacion = "liquidaciones"
  DB_Rendicion = "rendiciones"

  // jwt
  ExpiraToken   = 100000 // en minutos - Expiración del token para operar
  ExpiraTokenAut = 100000 // en minutos - Expiración del token de autorización
  Aes = "AES256Key-32Characters1234567890"
)

var (
	VerifyKey       *rsa.PublicKey
	SignKey         *rsa.PrivateKey
  FakeID          bson.ObjectId
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

  // config.FakeID
  FakeID = bson.ObjectIdHex("111111111111111111111111")

}
