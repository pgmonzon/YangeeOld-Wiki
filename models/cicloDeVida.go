package models

import (
  "time"

	"gopkg.in/mgo.v2/bson"
)

type CicloDeVida struct {
	ID				    bson.ObjectId   `bson:"_id" json:"id,omitempty"`
	Fecha    	    time.Time       `json:"fecha"`
  RemoteAddr    string          `json:"remoteAddr"`
  Metodo        string          `json:"metodo"`
  RequestURI    string          `json:"requestURI"`
  Protocolo     string          `json:"protocolo"`
  Codigo        int             `json:"codigo"`
  Duracion      time.Duration   `json:"duracion"`
  ClienteAPI_id bson.ObjectId   `bson:"clienteAPI_id" json:"clienteAPI_id,omitempty"`
  ClienteAPI    string          `json:"clienteAPI"`
  Usuario_id    bson.ObjectId   `bson:"usuario_id" json:"usuario_id,omitempty"`
  Usuario       string          `json:"usuario"`
  Empresa_id    bson.ObjectId   `bson:"empresa_id" json:"empresa_id,omitempty"`
  Empresa       string          `json:"empresa"`
}

type Audit struct {
  ID				    bson.ObjectId   `bson:"_id" json:"id,omitempty"`
  CicloVida_id  bson.ObjectId   `bson:"cicloVida_id" json:"cicloVida_id,omitempty"`
	Fecha    	    time.Time       `json:"fecha"`
  ClienteAPI_id bson.ObjectId   `bson:"clienteAPI_id" json:"clienteAPI_id,omitempty"`
  ClienteAPI    string          `json:"clienteAPI"`
  Usuario_id    bson.ObjectId   `bson:"usuario_id" json:"usuario_id,omitempty"`
  Usuario       string          `json:"usuario"`
  Empresa_id    bson.ObjectId   `bson:"empresa_id" json:"empresa_id,omitempty"`
  Empresa       string          `json:"empresa"`
  Collection    string          `json:"collection"`
  Objeto_id     bson.ObjectId   `bson:"objeto_id" json:"objeto_id,omitempty"`
  TipoOper      string          `json:"tipoOper"`
  Auditado      interface{}     `json:"audit,omitempty"`
}
