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


func PermisoCrear(w http.ResponseWriter, req *http.Request) {
	var permiso models.Permiso

  // Decode del JSON
  // ***************
  decoder := json.NewDecoder(req.Body)
  err := decoder.Decode(&permiso)
  if err != nil {
    core.RspMsgJSON(w, req, "ERROR", "JSON", "INVALID_PARAMS: JSON decode erróneo", http.StatusBadRequest)
    return
  }

  // Doy de alta
  // ***********
  estado, valor, mensaje, httpStat, permiso, existia := PermisoAlta(permiso, req)
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
  s := []string{"Agregó el permiso ", permiso.Permiso}
  core.RspMsgJSON(w, req, "OK", permiso.Permiso, strings.Join(s, ""), http.StatusCreated)
  return
}

// Devuelve Estado, Valor, Mensaje, HttpStat, Permiso, Existía
func PermisoAlta(permisoAlta models.Permiso, req *http.Request) (string, string, string, int, models.Permiso, bool) {
	var permiso models.Permiso

  // Verifico los campos obligatorios
  // ********************************
  if permisoAlta.Permiso == "" {
    s := []string{"INVALID_PARAMS: permiso no puede estar vacío"}
    return "ERROR", "PermisoAlta", strings.Join(s, ""), http.StatusBadRequest, permiso, false
  }

  // Me fijo si ya Existe
  // ********************
  estado, valor, mensaje, httpStat, permiso, existia := PermisoExiste(permisoAlta.Permiso)
  if httpStat != http.StatusOK || existia == true {
    return estado, valor, mensaje, httpStat, permiso, existia
  }

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    return "ERROR", "GetMongoSession", err.Error(), httpStat, permiso, false
  }
  defer session.Close()

  // Intento el alta
  // ***************
  objID := bson.NewObjectId()
  permiso.ID = objID
  permiso.Permiso = permisoAlta.Permiso
  permiso.Modulo_id = permisoAlta.Modulo_id
  permiso.Activo = permisoAlta.Activo
  permiso.Timestamp = time.Now()
  permiso.Borrado = false
  collection := session.DB(config.DB_Name).C(config.DB_Permiso)
  err = collection.Insert(permiso)
  if err != nil {
    s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
    return "ERROR", "Insert Permiso", strings.Join(s, ""), http.StatusInternalServerError, permiso, false
  }

  // Está todo Ok
  // ************
  core.Audit(req, config.DB_Permiso, permiso.ID, "Alta", permiso)
  return "OK", "PermisoAlta", "Ok", http.StatusOK, permiso, false
}

// Devuelve Estado, Valor, Mensaje, HttpStat, Empresa, Existía
func PermisoExiste(permisoExiste string) (string, string, string, int, models.Permiso, bool) {
  var permiso models.Permiso

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    return "ERROR", "GetMongoSession", err.Error(), httpStat, permiso, false
  }
  defer session.Close()

  // Me aseguro el índice
  // ********************
  collection := session.DB(config.DB_Name).C(config.DB_Permiso)
  index := mgo.Index{
    Key:        []string{"permiso"},
    Unique:     true,
    DropDups:   false,
    Background: true,
    Sparse:     true,
  }
  err = collection.EnsureIndex(index)
  if err != nil {
    s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
    return "ERROR", "EnsureIndex", strings.Join(s, ""), http.StatusInternalServerError, permiso, false
  }

  // Verifico si Existe
  // ******************
  collection.Find(bson.M{"permiso": permisoExiste}).One(&permiso)
  // No existe
  if permiso.ID == "" {
    return "OK", "BuscarPermiso", "Ok", http.StatusOK, permiso, false
  }
  // Existe borrado
  if permiso.Borrado == true {
    s := []string{"INVALID_PARAMS: El permiso ", permisoExiste," ya existe borrada"}
    return "ERROR", "BuscarPermiso", strings.Join(s, ""), http.StatusBadRequest, permiso, true
  }
  // Existe inactivo
  if permiso.Activo == false {
    s := []string{"INVALID_PARAMS: El permiso ", permisoExiste," ya existe inactiva"}
    return "ERROR", "BuscarPermiso", strings.Join(s, ""), http.StatusBadRequest, permiso, true
  }
  // Existe
  s := []string{"INVALID_PARAMS: El permiso ", permisoExiste," ya existe"}
  return "ERROR", "BuscarPermiso", strings.Join(s, ""), http.StatusBadRequest, permiso, true
}

// Devuelve Estado, Valor, Mensaje, HttpStat, Empresa
func Permiso_X_ID(permisoID bson.ObjectId) (string, string, string, int, models.Permiso) {
  var permiso models.Permiso

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    return "ERROR", "GetMongoSession", err.Error(), httpStat, permiso
  }
  defer session.Close()

  // Trato de traerlo
  // ****************
  collection := session.DB(config.DB_Name).C(config.DB_Permiso)
  collection.Find(bson.M{"_id": permisoID}).One(&permiso)
  // No existe
  if permiso.ID == "" {
    s := []string{"INVALID_PARAMS: El permiso no existe"}
    return "ERROR", "BuscarPermiso", strings.Join(s, ""), http.StatusBadRequest, permiso
  }
  // Existe
  return "OK", "BuscarPermiso", "Ok", http.StatusOK, permiso
}

func RolCrear(w http.ResponseWriter, req *http.Request) {
	var rol models.Rol

  // Decode del JSON
  // ***************
  decoder := json.NewDecoder(req.Body)
  err := decoder.Decode(&rol)
  if err != nil {
    core.RspMsgJSON(w, req, "ERROR", "JSON", "INVALID_PARAMS: JSON decode erróneo", http.StatusBadRequest)
    return
  }

  // Doy de alta
  // ***********
  estado, valor, mensaje, httpStat, rol, existia := RolAlta(rol, req)
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
  s := []string{"Agregó el rol ", rol.Rol}
  core.RspMsgJSON(w, req, "OK", rol.Rol, strings.Join(s, ""), http.StatusCreated)
  return
}

// Devuelve Estado, Valor, Mensaje, HttpStat, collection, Existía
func RolAlta(rolAlta models.Rol, req *http.Request) (string, string, string, int, models.Rol, bool) {
	var rol models.Rol

  // Verifico los campos obligatorios
  // ********************************
  if rolAlta.Rol == "" {
    s := []string{"INVALID_PARAMS: rol no puede estar vacío"}
    return "ERROR", "RolAlta", strings.Join(s, ""), http.StatusBadRequest, rol, false
  }

  // Me fijo si ya Existe
  // ********************
  estado, valor, mensaje, httpStat, rol, existia := RolExiste(rolAlta.Rol, rolAlta.Empresa_id)
  if httpStat != http.StatusOK || existia == true {
    return estado, valor, mensaje, httpStat, rol, existia
  }

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    return "ERROR", "GetMongoSession", err.Error(), httpStat, rol, false
  }
  defer session.Close()

  // Intento el alta
  // ***************
  objID := bson.NewObjectId()
  rol.ID = objID
  rol.Rol = rolAlta.Rol
  rol.Empresa_id = rolAlta.Empresa_id
  rol.Permisos = rolAlta.Permisos
  rol.Activo = rolAlta.Activo
  rol.Timestamp = time.Now()
  rol.Borrado = false
  collection := session.DB(config.DB_Name).C(config.DB_Rol)
  err = collection.Insert(rol)
  if err != nil {
    s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
    return "ERROR", "Insert Rol", strings.Join(s, ""), http.StatusInternalServerError, rol, false
  }

  // Está todo Ok
  // ************
  core.Audit(req, config.DB_Rol, rol.ID, "Alta", rol)
  return "OK", "RolAlta", "Ok", http.StatusOK, rol, false
}

// Devuelve Estado, Valor, Mensaje, HttpStat, collection, Existía
func RolExiste(rolExiste string, empresaExiste bson.ObjectId) (string, string, string, int, models.Rol, bool) {
  var rol models.Rol

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    return "ERROR", "GetMongoSession", err.Error(), httpStat, rol, false
  }
  defer session.Close()

  // Me aseguro el índice
  // ********************
  collection := session.DB(config.DB_Name).C(config.DB_Rol)
  index := mgo.Index{
    Key:        []string{"empresa_id", "rol"},
    Unique:     true,
    DropDups:   false,
    Background: true,
    Sparse:     true,
  }
  err = collection.EnsureIndex(index)
  if err != nil {
    s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
    return "ERROR", "EnsureIndex", strings.Join(s, ""), http.StatusInternalServerError, rol, false
  }

  // Verifico si Existe
  // ******************
  collection.Find(bson.M{"empresa_id": empresaExiste, "rol": rolExiste}).One(&rol)
  // No existe
  if rol.ID == "" {
    return "OK", "BuscarPermiso", "Ok", http.StatusOK, rol, false
  }
  // Existe borrado
  if rol.Borrado == true {
    s := []string{"INVALID_PARAMS: El rol ", rolExiste," ya existe borrado"}
    return "ERROR", "BuscarRol", strings.Join(s, ""), http.StatusBadRequest, rol, true
  }
  // Existe inactivo
  if rol.Activo == false {
    s := []string{"INVALID_PARAMS: El rol ", rolExiste," ya existe inactivo"}
    return "ERROR", "BuscarRol", strings.Join(s, ""), http.StatusBadRequest, rol, true
  }
  // Existe
  s := []string{"INVALID_PARAMS: El rol ", rolExiste," ya existe"}
  return "ERROR", "BuscarRol", strings.Join(s, ""), http.StatusBadRequest, rol, true
}

// Devuelve Estado, Valor, Mensaje, HttpStat, collection
func Rol_X_ID(rolID bson.ObjectId) (string, string, string, int, models.Rol) {
  var rol models.Rol

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    return "ERROR", "GetMongoSession", err.Error(), httpStat, rol
  }
  defer session.Close()

  // Trato de traerlo
  // ****************
  collection := session.DB(config.DB_Name).C(config.DB_Rol)
  collection.Find(bson.M{"_id": rolID}).One(&rol)
  // No existe
  if rol.ID == "" {
    s := []string{"INVALID_PARAMS: El rol no existe"}
    return "ERROR", "BuscarRol", strings.Join(s, ""), http.StatusBadRequest, rol
  }
  // Existe
  return "OK", "BuscarRol", "Ok", http.StatusOK, rol
}
