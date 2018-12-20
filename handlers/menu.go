package handlers

import (
  //"encoding/json"
  "net/http"
  //"fmt"
  "strings"
  //"strconv"
  //"time"

  "github.com/pgmonzon/Yangee/models"
  "github.com/pgmonzon/Yangee/core"
  "github.com/pgmonzon/Yangee/config"

  //"gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  //"github.com/gorilla/context"
)

func Menu_X_Empresa(documentoID bson.ObjectId, audit string, req *http.Request) (string, string, string, int, models.Menu) {
  var documento models.Menu
  coll := config.DB_Menu

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    return "ERROR", "GetMongoSession", err.Error(), httpStat, documento
  }
  defer session.Close()

  // Trato de traerlo
  // ****************
  collection := session.DB(config.DB_Name).C(coll)
  collection.Find(bson.M{"empresa_id": documentoID}).Select(bson.M{"empresa_id":0}).One(&documento)
  // No existe
  if documento.ID == "" {
    s := []string{"No encuentro el documento"}
    return "ERROR", audit, strings.Join(s, ""), http.StatusNonAuthoritativeInfo, documento
  }
  // Existe
  return "OK", audit, "Ok", http.StatusOK, documento
}
