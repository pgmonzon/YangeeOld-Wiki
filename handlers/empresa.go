package handlers

import (
  "encoding/json"
  "net/http"
  "strings"
  "time"
  "strconv"
  "math/rand"

  "github.com/pgmonzon/Yangee/models"
  "github.com/pgmonzon/Yangee/core"
  "github.com/pgmonzon/Yangee/config"

  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

func EmpresaCrear(w http.ResponseWriter, req *http.Request) {
	var empresa models.Empresa

  // Decode del JSON
  // ***************
  decoder := json.NewDecoder(req.Body)
  err := decoder.Decode(&empresa)
  if err != nil {
    core.RspMsgJSON(w, req, "ERROR", "JSON", "INVALID_PARAMS: JSON decode erróneo", http.StatusBadRequest)
    return
  }

  // Doy de alta la empresa
  // **********************
  estado, valor, mensaje, httpStat, empresa, existia := EmpresaAlta(empresa, req)
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
  s := []string{"Agregó la empresa ", empresa.Empresa}
  core.RspMsgJSON(w, req, "OK", empresa.Empresa, strings.Join(s, ""), http.StatusCreated)
  return
}

// Devuelve Estado, Valor, Mensaje, HttpStat, Empresa, Existía
func EmpresaAlta(empresaAlta models.Empresa, req *http.Request) (string, string, string, int, models.Empresa, bool) {
	var empresa models.Empresa

  // Verifico los campos obligatorios
  // ********************************
  if empresaAlta.Empresa == "" {
    s := []string{"INVALID_PARAMS: empresa no puede estar vacía"}
    return "ERROR", "EmpresaAlta", strings.Join(s, ""), http.StatusBadRequest, empresa, false
  }

  // Me fijo si ya Existe
  // ********************
  estado, valor, mensaje, httpStat, empresa, existia := EmpresaExiste(empresaAlta.Empresa)
  if httpStat != http.StatusOK || existia == true {
    return estado, valor, mensaje, httpStat, empresa, existia
  }

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    return "ERROR", "GetMongoSession", err.Error(), httpStat, empresa, false
  }
  defer session.Close()

  // Intento el alta
  // ***************
  objID := bson.NewObjectId()
  empresa.ID = objID
  empresa.Empresa = empresaAlta.Empresa
  empresa.Logo = empresaAlta.Logo
  empresa.Modulos = empresaAlta.Modulos
  empresa.Activo = empresaAlta.Activo
  empresa.Timestamp = time.Now()
  empresa.Borrado = false
  collection := session.DB(config.DB_Name).C(config.DB_Empresa)
  err = collection.Insert(empresa)
  if err != nil {
    s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
    return "ERROR", "Insert Empresa", strings.Join(s, ""), http.StatusInternalServerError, empresa, false
  }

  // Está todo Ok
  // ************
  core.Audit(req, config.DB_Empresa, empresa.ID, "Alta", empresa)
  return "OK", "EmpresaAlta", "Ok", http.StatusOK, empresa, false
}

// Devuelve Estado, Valor, Mensaje, HttpStat, Empresa, Existía
func EmpresaExiste(empresaExiste string) (string, string, string, int, models.Empresa, bool) {
  var empresa models.Empresa

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    return "ERROR", "GetMongoSession", err.Error(), httpStat, empresa, false
  }
  defer session.Close()

  // Me aseguro el índice
  // ********************
  collection := session.DB(config.DB_Name).C(config.DB_Empresa)
  index := mgo.Index{
    Key:        []string{"empresa"},
    Unique:     true,
    DropDups:   false,
    Background: true,
    Sparse:     true,
  }
  err = collection.EnsureIndex(index)
  if err != nil {
    s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
    return "ERROR", "EnsureIndex", strings.Join(s, ""), http.StatusInternalServerError, empresa, false
  }

  // Verifico si Existe
  // ******************
  collection.Find(bson.M{"empresa": empresaExiste}).One(&empresa)
  // No existe
  if empresa.ID == "" {
    return "OK", "BuscarEmpresa", "Ok", http.StatusOK, empresa, false
  }
  // Existe borrado
  if empresa.Borrado == true {
    s := []string{"INVALID_PARAMS: La empresa ", empresaExiste," ya existe borrada"}
    return "ERROR", "BuscarEmpresa", strings.Join(s, ""), http.StatusBadRequest, empresa, true
  }
  // Existe inactivo
  if empresa.Activo == false {
    s := []string{"INVALID_PARAMS: La empresa ", empresaExiste," ya existe inactiva"}
    return "ERROR", "BuscarEmpresa", strings.Join(s, ""), http.StatusBadRequest, empresa, true
  }
  // Existe
  s := []string{"INVALID_PARAMS: La empresa ", empresaExiste," ya existe"}
  return "ERROR", "BuscarEmpresa", strings.Join(s, ""), http.StatusBadRequest, empresa, true
}

// Devuelve Estado, Valor, Mensaje, HttpStat, Empresa
func Empresa_X_ID(empresaID bson.ObjectId) (string, string, string, int, models.Empresa) {
  var empresa models.Empresa

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    return "ERROR", "GetMongoSession", err.Error(), httpStat, empresa
  }
  defer session.Close()

  // Trato de traerlo
  // ****************
  collection := session.DB(config.DB_Name).C(config.DB_Empresa)
  collection.Find(bson.M{"_id": empresaID}).One(&empresa)
  // No existe
  if empresa.ID == "" {
    s := []string{"INVALID_PARAMS: La empresa no existe"}
    return "ERROR", "Buscar Empresa", strings.Join(s, ""), http.StatusBadRequest, empresa
  }
  // Existe
  return "OK", "BuscarEmpresa", "Ok", http.StatusOK, empresa
}

func EmpresaInvitar(w http.ResponseWriter, req *http.Request) {
	var empresaInvitacion models.EmpresaInvitacion

  // Decode del JSON
  // ***************
  decoder := json.NewDecoder(req.Body)
  err := decoder.Decode(&empresaInvitacion)
  if err != nil {
    core.RspMsgJSON(w, req, "ERROR", "JSON", "INVALID_PARAMS: JSON decode erróneo", http.StatusBadRequest)
    return
  }

  // Verifico los campos obligatorios
  // ********************************
  if empresaInvitacion.Empresa == "" || empresaInvitacion.Rol == "" || empresaInvitacion.Mail == "" {
    core.RspMsgJSON(w, req, "ERROR", "empresaInvitacion", "INVALID_PARAMS: empresa, rol y mail no pueden estar vacíos", http.StatusBadRequest)
    return
  }

  // Doy de alta la empresa
  // **********************
  var empresa models.Empresa
  empresa.Empresa = empresaInvitacion.Empresa
  s := []string{empresaInvitacion.Empresa, ".jpg"}
  empresa.Logo = strings.Join(s, "")
  empresa.Modulos = empresaInvitacion.Modulos
  empresa.Activo = true
  estado, valor, mensaje, httpStat, empresa, existia := EmpresaAlta(empresa, req)
  if httpStat != http.StatusOK && existia == false {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }

  // Doy de alta el Rol con los permisos de los módulos
  // **************************************************
  estado, valor, mensaje, httpStat, permisosID := ModulosPermisos(empresaInvitacion.Modulos)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }
  var rol models.Rol
  rol.Rol = empresaInvitacion.Rol
  rol.Empresa_id = empresa.ID
  rol.Permisos = permisosID
  rol.Activo = true
  estado, valor, mensaje, httpStat, rol, existia = RolAlta(rol, req)
  if httpStat != http.StatusOK && existia == false {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    core.RspMsgJSON(w, req, "ERROR", "GetMongoSession", err.Error(), httpStat)
  }
  defer session.Close()

  // Traigo el rol
  // *************
  rolesID := make([]models.IdRol, 0)
  collection := session.DB(config.DB_Name).C(config.DB_Rol)
  collection.Find(bson.M{"_id": rol.ID}).Select(bson.M{"_id": 1}).All(&rolesID)

  // Doy de alta el usuario
  // **********************
  var usuario models.Usuario
  usuario.Usuario = strconv.FormatInt(rand.Int63(), 10)
  usuario.Clave = strconv.FormatInt(rand.Int63(), 10)
  usuario.Mail = empresaInvitacion.Mail
  usuario.Apellido = ""
  usuario.Nombre = ""
  usuario.Empresa_id = empresa.ID
  usuario.Activo = true
  usuario.Roles = rolesID
  estado, valor, mensaje, httpStat, usuario, existia = UsuarioAlta(usuario, req)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }

  // Está todo Ok
  // ************
  core.RspMsgJSON(w, req, "OK", "Invitación Empresa", "Ok", http.StatusCreated)
  return
}
