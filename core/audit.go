package core

import (
  "log"
  "net/http"
  "time"

  "github.com/pgmonzon/Yangee/models"
  "github.com/pgmonzon/Yangee/config"

  "github.com/gorilla/context"
  "gopkg.in/mgo.v2/bson"
)

func Audit(req *http.Request, collection string, objeto_id bson.ObjectId, tipoOper string, auditado interface{}) {
  session, err, _ := GetMongoSession()
  if err == nil {
    defer session.Close()
    var audit models.Audit
    objID := bson.NewObjectId()
    audit.ID = objID
    audit.CicloVida_id = context.Get(req, "CicloDeVida_id").(bson.ObjectId)
    audit.Fecha = time.Now()
    audit.ClienteAPI_id = context.Get(req, "ClienteAPI_id").(bson.ObjectId)
    audit.ClienteAPI = context.Get(req, "ClienteAPI").(string)
    audit.Usuario_id = context.Get(req, "Usuario_id").(bson.ObjectId)
    audit.Usuario = context.Get(req, "Usuario").(string)
    audit.Empresa_id = context.Get(req, "Empresa_id").(bson.ObjectId)
    audit.Empresa = context.Get(req, "Empresa").(string)
    audit.Collection = collection
    audit.Objeto_id = objeto_id
    audit.TipoOper = tipoOper
    audit.Auditado = auditado

    collection := session.DB(config.DB_Name).C(config.DB_Audit)
    err :=collection.Insert(audit)
    if err != nil {
      log.Printf("Falló Insert auditoría")
    }
  } else {
    log.Printf("Falló GetMongoSession auditoría")
  }
}
