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
  Body          string          `json:"body"`
  //ClienteAPI_id bson.ObjectId   `bson:"_id" json:"clienteAPI_id,omitempty"`
  //ClienteAPI    string          `json:"clienteAPI"`
  //Usuario_id    bson.ObjectId   `bson:"_id" json:"usuario_id,omitempty"`
  //Usuario       string          `json:"usuario"`
/**
  Codigo        string          `json:"codigo"`
  Duracion      string          `json:"duracion"`
  Colleccion    string          `json:"coleccion"`
  Novedad       interface{}     `json:"novedad"`
**/
}
