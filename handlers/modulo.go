package handlers

import (
  "encoding/json"
  "net/http"
  "strings"
  "time"

  "github.com/pgmonzon/Yangee/models"
  "github.com/pgmonzon/Yangee/core"
  "github.com/pgmonzon/Yangee/config"

  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

func ModuloCrear(w http.ResponseWriter, req *http.Request) {
	var modulo models.Modulo

  // Decode del JSON
  // ***************
  decoder := json.NewDecoder(req.Body)
  err := decoder.Decode(&modulo)
  if err != nil {
    core.RspMsgJSON(w, req, "ERROR", "JSON", "INVALID_PARAMS: JSON decode erróneo", http.StatusBadRequest)
    return
  }

  // Doy de alta el módulo
  // *********************
  estado, valor, mensaje, httpStat, modulo, existia := ModuloAlta(modulo, req)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }
  if existia {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }

  // Está todo Ok
  // ************
  s := []string{"Agregó el módulo ", modulo.Modulo}
  core.RspMsgJSON(w, req, "OK", modulo.Modulo, strings.Join(s, ""), http.StatusCreated)
  return
}

// Devuelve Estado, Valor, Mensaje, HttpStat, Modulo, Existía
func ModuloAlta(moduloAlta models.Modulo, req *http.Request) (string, string, string, int, models.Modulo, bool) {
	var modulo models.Modulo

  // Verifico los campos obligatorios
  // ********************************
  if moduloAlta.Modulo == "" {
    s := []string{"INVALID_PARAMS: módulo no puede estar vacío"}
    return "ERROR", "ModuloAlta", strings.Join(s, ""), http.StatusBadRequest, modulo, false
  }

  // Me fijo si ya Existe
  // ********************
  estado, valor, mensaje, httpStat, modulo, existia := ModuloExiste(moduloAlta.Modulo)
  if httpStat != http.StatusOK || existia == true {
    return estado, valor, mensaje, httpStat, modulo, existia
  }

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    return "ERROR", "GetMongoSession", err.Error(), httpStat, modulo, false
  }
  defer session.Close()

  // Intento el alta
  // ***************
  objID := bson.NewObjectId()
  modulo.ID = objID
  modulo.Modulo = moduloAlta.Modulo
  modulo.Activo = moduloAlta.Activo
  modulo.Timestamp = time.Now()
  modulo.Borrado = false
  collection := session.DB(config.DB_Name).C(config.DB_Modulo)
  err = collection.Insert(modulo)
  if err != nil {
    s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
    return "ERROR", "Insert Módulo", strings.Join(s, ""), http.StatusInternalServerError, modulo, false
  }

  // Está todo Ok
  // ************
  core.Audit(req, config.DB_Modulo, modulo.ID, "Alta", modulo)
  return "OK", "ModuloAlta", "Ok", http.StatusOK, modulo, false
}

// Devuelve Estado, Valor, Mensaje, HttpStat, Collection, Existía
func ModuloExiste(moduloExiste string) (string, string, string, int, models.Modulo, bool) {
  var modulo models.Modulo

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    return "ERROR", "GetMongoSession", err.Error(), httpStat, modulo, false
  }
  defer session.Close()

  // Me aseguro el índice
  // ********************
  collection := session.DB(config.DB_Name).C(config.DB_Modulo)
  index := mgo.Index{
    Key:        []string{"modulo"},
    Unique:     true,
    DropDups:   false,
    Background: true,
    Sparse:     true,
  }
  err = collection.EnsureIndex(index)
  if err != nil {
    s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
    return "ERROR", "EnsureIndex", strings.Join(s, ""), http.StatusInternalServerError, modulo, false
  }

  // Verifico si Existe
  // ******************
  collection.Find(bson.M{"modulo": moduloExiste}).One(&modulo)
  // No existe
  if modulo.ID == "" {
    return "OK", "BuscarModulo", "Ok", http.StatusOK, modulo, false
  }
  // Existe borrado
  if modulo.Borrado == true {
    s := []string{"INVALID_PARAMS: El módulo ", moduloExiste," ya existe borrado"}
    return "ERROR", "BuscarModulo", strings.Join(s, ""), http.StatusBadRequest, modulo, true
  }
  // Existe inactivo
  if modulo.Activo == false {
    s := []string{"INVALID_PARAMS: El módulo ", moduloExiste," ya existe inactivo"}
    return "ERROR", "BuscarModulo", strings.Join(s, ""), http.StatusBadRequest, modulo, true
  }
  // Existe
  s := []string{"INVALID_PARAMS: El módulo ", moduloExiste," ya existe"}
  return "ERROR", "BuscarModulo", strings.Join(s, ""), http.StatusBadRequest, modulo, true
}

// Devuelve Estado, Valor, Mensaje, HttpStat, Collection
func Modulo_X_ID(moduloID bson.ObjectId) (string, string, string, int, models.Modulo) {
  var modulo models.Modulo

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    return "ERROR", "GetMongoSession", err.Error(), httpStat, modulo
  }
  defer session.Close()

  // Trato de traerlo
  // ****************
  collection := session.DB(config.DB_Name).C(config.DB_Modulo)
  collection.Find(bson.M{"_id": moduloID}).One(&modulo)
  // No existe
  if modulo.ID == "" {
    s := []string{"INVALID_PARAMS: El módulo no existe"}
    return "ERROR", "Buscar Módulo", strings.Join(s, ""), http.StatusBadRequest, modulo
  }
  // Existe
  return "OK", "Buscar Módulo", "Ok", http.StatusOK, modulo
}

// Devuelve Estado, Valor, Mensaje, HttpStat, Collection
func ModulosPermisos(modulosID []models.IdModulo) (string, string, string, int, []models.IdPermiso) {
  permisosID := make([]models.IdPermiso, 0)

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    return "ERROR", "GetMongoSession", err.Error(), httpStat, permisosID
  }
  defer session.Close()

  // recorro los módulos
  // *******************
  modulosArr := []bson.ObjectId{}
  for _, itemMod := range modulosID {
    if itemMod.ID != "" {
      _, _, _, httpStat, modulo := Modulo_X_ID(itemMod.ID)
      if httpStat == http.StatusOK && modulo.Activo == true && modulo.Borrado == false {
        modulosArr = append(modulosArr, itemMod.ID)
      }
    }
  }

  // Traigo los permisos pertenecientes a los módulos
  // ************************************************
  collection := session.DB(config.DB_Name).C(config.DB_Permiso)
  collection.Find(bson.M{"modulo_id": bson.M{"$in": modulosArr}, "activo": true, "borrado": false}).Select(bson.M{"_id": 1}).All(&permisosID)

  // Está todo Ok
  // ************
  return "OK", "ModulosPermisos", "Ok", http.StatusOK, permisosID
}
