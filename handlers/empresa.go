package handlers

import (
  "net/http"
  "fmt"
  "strings"

  "github.com/pgmonzon/Yangee/models"
  "github.com/pgmonzon/Yangee/core"
  "github.com/pgmonzon/Yangee/config"

  "gopkg.in/mgo.v2/bson"
)

func Empresa_X_ID(empresaID bson.ObjectId) (models.Empresa, error, int) {
  var empresa models.Empresa

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, _ := core.GetMongoSession()
  if err != nil {
    s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
    return empresa, fmt.Errorf(strings.Join(s, "")), http.StatusInternalServerError
  }
  defer session.Close()

  // Trato de traerlo
  // ****************
  collection := session.DB(config.DB_Name).C(config.DB_Empresa)
  collection.Find(bson.M{"_id": empresaID}).One(&empresa)
  // No existe
  if empresa.ID == "" {
    s := []string{"INVALID_PARAMS: La empresa no existe"}
    return empresa, fmt.Errorf(strings.Join(s, "")), http.StatusBadRequest
  }
  // Existe
  return empresa, nil, http.StatusOK
}
