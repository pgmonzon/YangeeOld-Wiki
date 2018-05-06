package handlers

import (
  "encoding/json"
  "net/http"
  "fmt"
  "strings"
  "strconv"
  "time"

  "github.com/pgmonzon/Yangee/models"
  "github.com/pgmonzon/Yangee/core"
  "github.com/pgmonzon/Yangee/config"

  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

func UsuarioCrear(w http.ResponseWriter, req *http.Request) {
	var usuario models.Usuario

  // Decode del JSON
  // ***************
  decoder := json.NewDecoder(req.Body)
  err := decoder.Decode(&usuario)
  if err != nil {
    core.RspMsgJSON(w, req, "ERROR", "JSON", "INVALID_PARAMS: JSON decode erróneo", http.StatusBadRequest)
    return
  }

  // Doy de alta
  // ***********
  estado, valor, mensaje, httpStat, usuario, existia := UsuarioAlta(usuario, req)
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
  s := []string{"Agregó el usuario ", usuario.Usuario}
  core.RspMsgJSON(w, req, "OK", usuario.Usuario, strings.Join(s, ""), http.StatusCreated)
  return
}

// Devuelve Estado, Valor, Mensaje, HttpStat, collection, Existía
func UsuarioAlta(usuarioAlta models.Usuario, req *http.Request) (string, string, string, int, models.Usuario, bool) {
	var usuario models.Usuario

  // Verifico los campos obligatorios
  // ********************************
  if usuarioAlta.Usuario == "" || usuarioAlta.Clave == "" || usuarioAlta.Mail == "" {
    s := []string{"INVALID_PARAMS: usuario, clave y mail no pueden estar vacíos"}
    return "ERROR", "UsuarioAlta", strings.Join(s, ""), http.StatusBadRequest, usuario, false
  }

  // Me fijo si ya Existe
  // ********************
  estado, valor, mensaje, httpStat, usuario, existia := UsuarioExiste(usuarioAlta.Usuario)
  if httpStat != http.StatusOK || existia == true {
    return estado, valor, mensaje, httpStat, usuario, existia
  }

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    return "ERROR", "GetMongoSession", err.Error(), httpStat, usuario, false
  }
  defer session.Close()

  // Intento el alta
  // ***************
  objID := bson.NewObjectId()
  usuario.ID = objID
  usuario.Usuario = usuarioAlta.Usuario
  usuario.Clave = strconv.FormatInt(core.HashSha512(usuarioAlta.Clave),16)
  usuario.Mail = usuarioAlta.Mail
  usuario.Apellido = usuarioAlta.Apellido
  usuario.Nombre = usuarioAlta.Nombre
  usuario.Empresa_id = usuarioAlta.Empresa_id
  usuario.Roles = usuarioAlta.Roles
  usuario.Menu = usuarioAlta.Menu
  usuario.Activo = usuarioAlta.Activo
  usuario.Timestamp = time.Now()
  usuario.Borrado = false
  collection := session.DB(config.DB_Name).C(config.DB_Usuario)
  err = collection.Insert(usuario)
  if err != nil {
    s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
    return "ERROR", "Insert Usuario", strings.Join(s, ""), http.StatusInternalServerError, usuario, false
  }

  // Está todo Ok
  // ************
  core.Audit(req, config.DB_Usuario, usuario.ID, "Alta", usuario)
  return "OK", "UsuarioAlta", "Ok", http.StatusOK, usuario, false
}

// Devuelve Estado, Valor, Mensaje, HttpStat, collection, Existía
func UsuarioExiste(usuarioExiste string) (string, string, string, int, models.Usuario, bool) {
  var usuario models.Usuario

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    return "ERROR", "GetMongoSession", err.Error(), httpStat, usuario, false
  }
  defer session.Close()

  // Me aseguro el índice
  // ********************
  collection := session.DB(config.DB_Name).C(config.DB_Usuario)
  index := mgo.Index{
    Key:        []string{"usuario"},
    Unique:     true,
    DropDups:   false,
    Background: true,
    Sparse:     true,
  }
  err = collection.EnsureIndex(index)
  if err != nil {
    s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
    return "ERROR", "EnsureIndex", strings.Join(s, ""), http.StatusInternalServerError, usuario, false
  }

  // Verifico si Existe
  // ******************
  collection.Find(bson.M{"usuario": usuarioExiste}).One(&usuario)
  // No existe
  if usuario.ID == "" {
    return "OK", "BuscarUsuario", "Ok", http.StatusOK, usuario, false
  }
  // Existe borrado
  if usuario.Borrado == true {
    s := []string{"INVALID_PARAMS: El usuario ", usuarioExiste," ya existe borrado"}
    return "ERROR", "BuscarUsuario", strings.Join(s, ""), http.StatusBadRequest, usuario, true
  }
  // Existe inactivo
  if usuario.Activo == false {
    s := []string{"INVALID_PARAMS: El usuario ", usuarioExiste," ya existe inactivo"}
    return "ERROR", "BuscarUsuario", strings.Join(s, ""), http.StatusBadRequest, usuario, true
  }
  // Existe
  s := []string{"INVALID_PARAMS: El usuario ", usuarioExiste," ya existe"}
  return "ERROR", "BuscarUsuario", strings.Join(s, ""), http.StatusBadRequest, usuario, true
}

// Devuelve Estado, Valor, Mensaje, HttpStat, collection
func Usuario_X_ID(usuarioID bson.ObjectId) (string, string, string, int, models.Usuario) {
  var usuario models.Usuario

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    return "ERROR", "GetMongoSession", err.Error(), httpStat, usuario
  }
  defer session.Close()

  // Trato de traerlo
  // ****************
  collection := session.DB(config.DB_Name).C(config.DB_Usuario)
  collection.Find(bson.M{"_id": usuarioID}).One(&usuario)
  // No existe
  if usuario.ID == "" {
    s := []string{"INVALID_PARAMS: El usuario no existe"}
    return "ERROR", "BuscarUsuario", strings.Join(s, ""), http.StatusBadRequest, usuario
  }
  // Existe
  return "OK", "BuscarUsuario", "Ok", http.StatusOK, usuario
}

func UsuarioPermisos(usuarioPermisos string) (string, error, int) {

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, _ := core.GetMongoSession()
  if err != nil {
    s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
    return "", fmt.Errorf(strings.Join(s, "")), http.StatusInternalServerError
  }

  // Busco el usuario y verifico que esté activo
  // *******************************************
  defer session.Close()
  cUsuario := session.DB(config.DB_Name).C(config.DB_Usuario)
  var usuario models.Usuario
  cUsuario.Find(bson.M{"usuario": usuarioPermisos, "activo": true, "borrado": false}).One(&usuario)
  if usuario.ID == "" {
    s := []string{"INVALID_PARAMS: El usuario no existe o está inactivo"}
    return "", fmt.Errorf(strings.Join(s, "")), http.StatusBadRequest
  }

  // Traigo la empresa del usuario
  // *****************************
  _, _, _, httpStat, empresa := Empresa_X_ID(usuario.Empresa_id)
  if httpStat != http.StatusOK || empresa.Activo == false || empresa.Borrado == true {
    s := []string{"INVALID_PARAMS: La empresa del usuario no existe o está inactiva"}
    return "", fmt.Errorf(strings.Join(s, "")), http.StatusBadRequest
  }

  // Busco los módulos de la empresa del usuario
  // *******************************************
  modulosArr := []bson.ObjectId{}
  for _, itemMod := range empresa.Modulos {
    if itemMod.ID != "" {
      modulosArr = append(modulosArr, itemMod.ID)
    }
  }

  // Obtengo los ID roles del usuario
  // ********************************
  rolesArr := []bson.ObjectId{}
  for _, item := range usuario.Roles {
    if item.ID != "" {
      rolesArr = append(rolesArr, item.ID)
    }
  }
  roles := make([]models.Rol, 0)
  cRoles := session.DB(config.DB_Name).C(config.DB_Rol)
  cRoles.Find(bson.M{"_id": bson.M{"$in": rolesArr}}).All(&roles)

  // Obtengo los ID permisos de los roles
  // ************************************
  permisosArr := []bson.ObjectId{}
  for _, itemRol := range roles {
    for _, itemPermiso := range itemRol.Permisos {
      if itemPermiso.ID != "" {
        permisosArr = append(permisosArr, itemPermiso.ID)
      }
    }
  }
  permisos := make([]models.Permiso, 0)
  cPermisos := session.DB(config.DB_Name).C(config.DB_Permiso)
  cPermisos.Find(bson.M{"_id": bson.M{"$in": permisosArr}, "modulo_id": bson.M{"$in": modulosArr}}).All(&permisos)

  // Junto los permisos en un string
  // *******************************
  permisosStr := []string{}
  permisosStr = append(permisosStr, "#")
  for _, itemItem := range permisos {
    if itemItem.Permiso != "" {
      permisosStr = append(permisosStr, itemItem.Permiso)
    }
  }
  permisosStr = append(permisosStr, "#")
  permisosUsuario := strings.Join(permisosStr, "#")

  // Está todo Ok
  // ************
  return permisosUsuario, nil, http.StatusOK
}

func UsuarioLogin(usuarioLogin string, claveLogin string) (string, string, string, int, models.Usuario, models.Empresa) {
  var usuario models.Usuario
  var empresa models.Empresa

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    return "ERROR", "GetMongoSession", err.Error(), httpStat, usuario, empresa
  }

  // Intento el login
  // ****************
  defer session.Close()
  collection := session.DB(config.DB_Name).C(config.DB_Usuario)
  collection.Find(bson.M{"usuario": usuarioLogin, "clave": strconv.FormatInt(core.HashSha512(claveLogin),16), "activo": true, "borrado": false}).One(&usuario)
  // Si no loguea
  if usuario.ID == "" {
    s := []string{"Usuario y clave incorrectos"}
    return "ERROR", "Login", strings.Join(s, ""), http.StatusNonAuthoritativeInfo, usuario, empresa
  }

  // Traigo la empresa del usuario
  // *****************************
  estado, valor, mensaje, httpStat, empresa := Empresa_X_ID(usuario.Empresa_id)
  if httpStat != http.StatusOK {
    return estado, valor, mensaje, httpStat, usuario, empresa
  }

  // Está todo Ok
  // ************
  return "OK", "Login", "Ok", http.StatusOK, usuario, empresa
}
