package handlers

import (
  "encoding/json"
  "net/http"
  "strings"
  "strconv"
  "time"

  "github.com/pgmonzon/Yangee/models"
  "github.com/pgmonzon/Yangee/core"
  "github.com/pgmonzon/Yangee/config"

  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "github.com/gorilla/context"
  "github.com/gorilla/mux"
)

/*
!!! PONER EL FIND EN CASE SENSITIVE !!!
*** Reemplazos automáticos ***
reemplazar Personales (mayúscula plural) 4 apariciones
reemplazar Personal (mayúscula singular) 76 apariciones
reemplazar personal (minúscula singular) 7 apariciones IMPORANTE: no puede tener mayúsculas

*** Reemplazos manuales ***
reemplazar "No podés dejar vacío" 3 apariciones
reemplazar "en forma manual" 7 apariciones
reemplazar "en orden" 2 apariciones
*/

func PersonalCrear(w http.ResponseWriter, req *http.Request) {
  //-------------------Modificar ###### estas 2 variables
	var documento models.Personal
  audit := "Crear"

  // Decode del JSON
  // ***************
  decoder := json.NewDecoder(req.Body)
  err := decoder.Decode(&documento)
  if err != nil {
    core.RspMsgJSON(w, req, "ERROR", "JSON", "INVALID_PARAMS: JSON decode erróneo", http.StatusBadRequest)
    return
  }

  // Doy de alta
  // ***********
  //------------------------------------------------------Modificar ######
  estado, valor, mensaje, httpStat, documento, existia := PersonalAlta(documento, req, audit)
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
  //------------------------------------Modificar ######
  s := []string{"Agregaste ", documento.Nombre, " ", documento.Apellido}
  //--------------------------------------Modificar ######
  core.RspMsgJSON(w, req, "OK", documento.Apellido, strings.Join(s, ""), http.StatusCreated)
  return
}

// Devuelve Estado, Valor, Mensaje, HttpStat, Collection, Existía
func PersonalAlta(documentoAlta models.Personal, req *http.Request, audit string) (string, string, string, int, models.Personal, bool) {
  //-------------------Modificar ###### las 3 variables
	var documento models.Personal
  camposVacios := "No podés dejar vacío apellido y nombre"
  coll := config.DB_Personal
  empresaID := context.Get(req, "Empresa_id").(bson.ObjectId)

  // Verifico los campos obligatorios
  // ********************************
  //---------------Modificar ######
  if documentoAlta.Apellido == "" || documentoAlta.Nombre == "" {
    s := []string{camposVacios}
    return "ERROR", "Alta", strings.Join(s, ""), http.StatusNonAuthoritativeInfo, documento, false
  }

  // Me fijo si ya Existe
  // ********************
  //-----------------------------------------------------Modificar ######--------------Modificar ######
  estado, valor, mensaje, httpStat, documento, existia := PersonalExiste(documentoAlta.Apellido, documentoAlta.Nombre, req)
  if httpStat != http.StatusOK || existia == true {
    return estado, valor, mensaje, httpStat, documento, existia
  }

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    return "ERROR", "GetMongoSession", err.Error(), httpStat, documento, false
  }
  defer session.Close()

  // Intento el alta
  // ***************
  documento = documentoAlta
  objID := bson.NewObjectId()
  documento.ID = objID
  documento.Empresa_id = empresaID
  documento.Timestamp = time.Now()
  documento.Borrado = false
  collection := session.DB(config.DB_Name).C(coll)
  err = collection.Insert(documento)
  if err != nil {
    s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
    return "ERROR", audit, strings.Join(s, ""), http.StatusInternalServerError, documento, false
  }

  // Está todo Ok
  // ************
  core.Audit(req, coll, documento.ID, audit, documento)
  return "OK", audit, "Ok", http.StatusOK, documento, false
}

// Devuelve Estado, Valor, Mensaje, HttpStat, Collection, Existía
func PersonalExiste(documentoApellidoExiste string, documentoNombreExiste string, req *http.Request) (string, string, string, int, models.Personal, bool) {
  //-------------------Modificar ###### las 3 variables
  var documento models.Personal
  indice := []string{"empresa_id", "apellido", "nombre"}
  coll := config.DB_Personal
  empresaID := context.Get(req, "Empresa_id").(bson.ObjectId)

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    return "ERROR", "GetMongoSession", err.Error(), httpStat, documento, false
  }
  defer session.Close()

  // Me aseguro el índice
  // ********************
  collection := session.DB(config.DB_Name).C(coll)
  index := mgo.Index{
    Key:        indice,
    Unique:     true,
    DropDups:   false,
    Background: true,
    Sparse:     true,
  }
  err = collection.EnsureIndex(index)
  if err != nil {
    s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
    return "ERROR", "EnsureIndex", strings.Join(s, ""), http.StatusInternalServerError, documento, false
  }

  // Verifico si Existe
  // ******************
  //----------------------------------------------Modificar ######
  collection.Find(bson.M{"empresa_id": empresaID, "apellido": documentoApellidoExiste, "nombre": documentoNombreExiste}).One(&documento)
  // No existe
  if documento.ID == "" {
    return "OK", "Buscar", "Ok", http.StatusOK, documento, false
  }
  // Existe borrado
  if documento.Borrado == true {
    s := []string{documentoApellidoExiste," existe borrado"}
    return "ERROR", "Buscar", strings.Join(s, ""), http.StatusNonAuthoritativeInfo, documento, true
  }
  // Existe inactivo
  if documento.Activo == false {
    s := []string{documentoApellidoExiste," existe inactivo"}
    return "ERROR", "Buscar", strings.Join(s, ""), http.StatusNonAuthoritativeInfo, documento, true
  }
  // Existe
  s := []string{documentoApellidoExiste," ya existe"}
  return "ERROR", "Buscar", strings.Join(s, ""), http.StatusNonAuthoritativeInfo, documento, true
}

func PersonalesTraer(w http.ResponseWriter, req *http.Request) {
  //-------------------Modificar ###### estas 2 variables
  var documento models.Personal
  var documentos []models.Personal
  vars := mux.Vars(req)
  orden := vars["orden"]
  limite := vars["limite"]

  // Verifico el formato del campo limite
  // ************************************
  limiteInt, err := strconv.Atoi(limite)
  if err != nil {
    core.RspMsgJSON(w, req, "ERROR", "Límite debe ser numérico", err.Error(), http.StatusBadRequest)
    return
  }

  // Decode del JSON
  // ***************
  decoder := json.NewDecoder(req.Body)
  err = decoder.Decode(&documento)
  if err != nil {
    core.RspMsgJSON(w, req, "ERROR", "JSON", "INVALID_PARAMS: JSON decode erróneo", http.StatusBadRequest)
    return
  }

  // Busco
  // *****
  //----------------------------------------------Modificar ######
  estado, valor, mensaje, httpStat, documentos := PersonalesBuscar(documento, orden, limiteInt, false, "Buscar", req)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }

  // Está todo Ok
  // ************
  respuesta, error := json.Marshal(documentos)
  core.FatalErr(error)
  core.RspJSON(w, req, respuesta, http.StatusOK)
  return
}

// Devuelve Estado, Valor, Mensaje, HttpStat, Collection
func PersonalesBuscar(documento models.Personal, orden string, limiteInt int, borrados bool, audit string, req *http.Request) (string, string, string, int, []models.Personal) {
  //----------------------Modificar ###### estas 2 variables
  var documentos []models.Personal
  coll := config.DB_Personal
  empresaID := context.Get(req, "Empresa_id").(bson.ObjectId)

  // Verifico que el campo orden sea Unique
  // **************************************
  //-----------Modificar ###### en orden
  if orden != "apellido" && orden != "-apellido" && orden != "nombre" && orden != "-nombre" {
    s := []string{"No puedo ordenar por ", orden}
    return "ERROR", "Buscar", strings.Join(s, ""), http.StatusNonAuthoritativeInfo, documentos
  }

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    return "ERROR", "GetMongoSession", err.Error(), httpStat, documentos
  }
  defer session.Close()

  // Trato de traerlos
  // *****************
  //----------Modificar ###### en forma manual
  selector := bson.M{
    "empresa_id": empresaID,
    "apellido": bson.M{"$regex": bson.RegEx{documento.Apellido, "i"}},
    "nombre": bson.M{"$regex": bson.RegEx{documento.Nombre, "i"}},
    "borrado": borrados,
  }
  collection := session.DB(config.DB_Name).C(coll)
  collection.Find(selector).Select(bson.M{"empresa_id":0}).Sort(orden).Limit(limiteInt).All(&documentos)

  // Si el resultado es vacío devuelvo ERROR
  // ***************************************
  if len(documentos) == 0 {
    s := []string{"No encontré documentos"}
    return "ERROR", audit, strings.Join(s, ""), http.StatusNonAuthoritativeInfo, documentos
  }

  // Está todo Ok
  // ************
  return "OK", audit, "Ok", http.StatusOK, documentos
}

func PersonalTraer(w http.ResponseWriter, req *http.Request) {
  vars := mux.Vars(req)
  ID := vars["docID"]

  // Verifico el formato del campo ID
  // ********************************
  if bson.IsObjectIdHex(ID) != true {
    core.RspMsgJSON(w, req, "ERROR", ID, "INVALID_PARAMS: Formato ID incorrecto", http.StatusBadRequest)
    return
  }
  documentoID := bson.ObjectIdHex(ID)

  // Busco
  // *****
  //----------------------------------------------Modificar ######
  estado, valor, mensaje, httpStat, documento := Personal_X_ID(documentoID, "Buscar ID", req)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }

  // Está todo Ok
  // ************
  respuesta, error := json.Marshal(documento)
  core.FatalErr(error)
  core.RspJSON(w, req, respuesta, http.StatusOK)
  return
}

// Devuelve Estado, Valor, Mensaje, HttpStat, collection
func Personal_X_ID(documentoID bson.ObjectId, audit string, req *http.Request) (string, string, string, int, models.Personal) {
  //-------------------Modificar ######
  var documento models.Personal
  coll := config.DB_Personal
  empresaID := context.Get(req, "Empresa_id").(bson.ObjectId)

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
  collection.Find(bson.M{"_id": documentoID, "empresa_id": empresaID}).Select(bson.M{"empresa_id":0}).One(&documento)
  // No existe
  if documento.ID == "" {
    s := []string{"No encuentro el documento"}
    return "ERROR", audit, strings.Join(s, ""), http.StatusNonAuthoritativeInfo, documento
  }
  // Existe
  return "OK", audit, "Ok", http.StatusOK, documento
}

func PersonalGuardar(w http.ResponseWriter, req *http.Request) {
  //-------------------Modificar ######
  var documento models.Personal
  vars := mux.Vars(req)
  ID := vars["docID"]
  audit := "Guardar"

  // Verifico el formato del campo ID
  // ********************************
  if bson.IsObjectIdHex(ID) != true {
    core.RspMsgJSON(w, req, "ERROR", ID, "INVALID_PARAMS: Formato ID incorrecto", http.StatusBadRequest)
    return
  }
  documentoID := bson.ObjectIdHex(ID)

  // Decode del JSON
  // ***************
  decoder := json.NewDecoder(req.Body)
  err := decoder.Decode(&documento)
  if err != nil {
    core.RspMsgJSON(w, req, "ERROR", "JSON", "INVALID_PARAMS: JSON decode erróneo", http.StatusBadRequest)
    return
  }

  // Busco para obtener los campos faltantes
  // ***************************************
  //------------------------------------------------------Modificar ######
  estado, valor, mensaje, httpStat, documentoExistente := Personal_X_ID(documentoID, "Buscar ID", req)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }
  documento.Borrado = documentoExistente.Borrado

  // Modifico
  // ********
  //----------------------------------Modificar ######
  estado, valor, mensaje, httpStat = PersonalModificar(documentoID, documento, req, audit)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }

  // Está todo Ok
  // ************
  //------------------------------------Modificar ######
  s := []string{"Guardaste ", documento.Apellido}
  //--------------------------------------Modificar ######
  core.RspMsgJSON(w, req, "OK", documento.Apellido, strings.Join(s, ""), http.StatusAccepted)
  return
}

func PersonalHabilitar(w http.ResponseWriter, req *http.Request) {
  //-------------------Modificar ######
  var documento models.Personal
  vars := mux.Vars(req)
  ID := vars["docID"]
  audit := "Habilitar"

  // Verifico el formato del campo ID
  // ********************************
  if bson.IsObjectIdHex(ID) != true {
    core.RspMsgJSON(w, req, "ERROR", ID, "INVALID_PARAMS: Formato ID incorrecto", http.StatusBadRequest)
    return
  }
  documentoID := bson.ObjectIdHex(ID)

  // Busco para obtener los campos faltantes
  // ***************************************
  //------------------------------------------------------Modificar ######
  estado, valor, mensaje, httpStat, documentoExistente := Personal_X_ID(documentoID, "Buscar ID", req)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }
  //-------Modificar ###### en forma manual
  documento.Apellido = documentoExistente.Apellido
  documento.Nombre = documentoExistente.Nombre
  documento.Celular = documentoExistente.Celular
  documento.Direccion = documentoExistente.Direccion
  documento.Latitud = documentoExistente.Latitud
  documento.Longitud = documentoExistente.Longitud
  documento.Categoria_id = documentoExistente.Categoria_id
  documento.Propio = documentoExistente.Propio
  documento.BasicoSindicato_id = documentoExistente.BasicoSindicato_id
  documento.Comision = documentoExistente.Comision
  documento.Curso = documentoExistente.Curso
  documento.LNH = documentoExistente.LNH
  documento.Registro = documentoExistente.Registro
  documento.FechaAlta = documentoExistente.FechaAlta
  documento.FechaBaja = documentoExistente.FechaBaja
  documento.Activo = true
  documento.Borrado = documentoExistente.Borrado

  // Modifico
  // ********
  //----------------------------------Modificar ######
  estado, valor, mensaje, httpStat = PersonalModificar(documentoID, documento, req, audit)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }

  // Está todo Ok
  // ************
  //------------------------------------Modificar ######
  s := []string{"Habilitaste ", documento.Apellido}
  //--------------------------------------Modificar ######
  core.RspMsgJSON(w, req, "OK", documento.Apellido, strings.Join(s, ""), http.StatusAccepted)
  return
}

func PersonalDeshabilitar(w http.ResponseWriter, req *http.Request) {
  //-------------------Modificar ######
  var documento models.Personal
  vars := mux.Vars(req)
  ID := vars["docID"]
  audit := "Deshabilitar"

  // Verifico el formato del campo ID
  // ********************************
  if bson.IsObjectIdHex(ID) != true {
    core.RspMsgJSON(w, req, "ERROR", ID, "INVALID_PARAMS: Formato ID incorrecto", http.StatusBadRequest)
    return
  }
  documentoID := bson.ObjectIdHex(ID)

  // Busco para obtener los campos faltantes
  // ***************************************
  //------------------------------------------------------Modificar ######
  estado, valor, mensaje, httpStat, documentoExistente := Personal_X_ID(documentoID, "Buscar ID", req)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }
  //-------Modificar ###### en forma manual
  documento.Apellido = documentoExistente.Apellido
  documento.Nombre = documentoExistente.Nombre
  documento.Celular = documentoExistente.Celular
  documento.Direccion = documentoExistente.Direccion
  documento.Latitud = documentoExistente.Latitud
  documento.Longitud = documentoExistente.Longitud
  documento.Categoria_id = documentoExistente.Categoria_id
  documento.Propio = documentoExistente.Propio
  documento.BasicoSindicato_id = documentoExistente.BasicoSindicato_id
  documento.Comision = documentoExistente.Comision
  documento.Curso = documentoExistente.Curso
  documento.LNH = documentoExistente.LNH
  documento.Registro = documentoExistente.Registro
  documento.FechaAlta = documentoExistente.FechaAlta
  documento.FechaBaja = documentoExistente.FechaBaja
  documento.Activo = false
  documento.Borrado = documentoExistente.Borrado

  // Modifico
  // ********
  //----------------------------------Modificar ######
  estado, valor, mensaje, httpStat = PersonalModificar(documentoID, documento, req, audit)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }

  // Está todo Ok
  // ************
  //------------------------------------Modificar ######
  s := []string{"Deshabilitaste ", documento.Apellido}
  //--------------------------------------Modificar ######
  core.RspMsgJSON(w, req, "OK", documento.Apellido, strings.Join(s, ""), http.StatusAccepted)
  return
}

func PersonalBorrar(w http.ResponseWriter, req *http.Request) {
  //-------------------Modificar ######
  var documento models.Personal
  vars := mux.Vars(req)
  ID := vars["docID"]
  audit := "Borrar"

  // Verifico el formato del campo ID
  // ********************************
  if bson.IsObjectIdHex(ID) != true {
    core.RspMsgJSON(w, req, "ERROR", ID, "INVALID_PARAMS: Formato ID incorrecto", http.StatusBadRequest)
    return
  }
  documentoID := bson.ObjectIdHex(ID)

  // Busco para obtener los campos faltantes
  // ***************************************
  //------------------------------------------------------Modificar ######
  estado, valor, mensaje, httpStat, documentoExistente := Personal_X_ID(documentoID, "Buscar ID", req)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }
  //-------Modificar ###### en forma manual
  documento.Apellido = documentoExistente.Apellido
  documento.Nombre = documentoExistente.Nombre
  documento.Celular = documentoExistente.Celular
  documento.Direccion = documentoExistente.Direccion
  documento.Latitud = documentoExistente.Latitud
  documento.Longitud = documentoExistente.Longitud
  documento.Categoria_id = documentoExistente.Categoria_id
  documento.Propio = documentoExistente.Propio
  documento.BasicoSindicato_id = documentoExistente.BasicoSindicato_id
  documento.Comision = documentoExistente.Comision
  documento.Curso = documentoExistente.Curso
  documento.LNH = documentoExistente.LNH
  documento.Registro = documentoExistente.Registro
  documento.FechaAlta = documentoExistente.FechaAlta
  documento.FechaBaja = documentoExistente.FechaBaja
  documento.Activo = documentoExistente.Activo
  documento.Borrado = true

  // Modifico
  // ********
  //----------------------------------Modificar ######
  estado, valor, mensaje, httpStat = PersonalModificar(documentoID, documento, req, audit)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }

  // Está todo Ok
  // ************
  //------------------------------------Modificar ######
  s := []string{"Borraste ", documento.Apellido}
  //--------------------------------------Modificar ######
  core.RspMsgJSON(w, req, "OK", documento.Apellido, strings.Join(s, ""), http.StatusAccepted)
  return
}

func PersonalRecuperar(w http.ResponseWriter, req *http.Request) {
  //-------------------Modificar ######
  var documento models.Personal
  vars := mux.Vars(req)
  ID := vars["docID"]
  audit := "Recuperar"

  // Verifico el formato del campo ID
  // ********************************
  if bson.IsObjectIdHex(ID) != true {
    core.RspMsgJSON(w, req, "ERROR", ID, "INVALID_PARAMS: Formato ID incorrecto", http.StatusBadRequest)
    return
  }
  documentoID := bson.ObjectIdHex(ID)

  // Busco para obtener los campos faltantes
  // ***************************************
  //------------------------------------------------------Modificar ######
  estado, valor, mensaje, httpStat, documentoExistente := Personal_X_ID(documentoID, "Buscar ID", req)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }
  //-------Modificar ###### en forma manual
  documento.Apellido = documentoExistente.Apellido
  documento.Nombre = documentoExistente.Nombre
  documento.Celular = documentoExistente.Celular
  documento.Direccion = documentoExistente.Direccion
  documento.Latitud = documentoExistente.Latitud
  documento.Longitud = documentoExistente.Longitud
  documento.Categoria_id = documentoExistente.Categoria_id
  documento.Propio = documentoExistente.Propio
  documento.BasicoSindicato_id = documentoExistente.BasicoSindicato_id
  documento.Comision = documentoExistente.Comision
  documento.Curso = documentoExistente.Curso
  documento.LNH = documentoExistente.LNH
  documento.Registro = documentoExistente.Registro
  documento.FechaAlta = documentoExistente.FechaAlta
  documento.FechaBaja = documentoExistente.FechaBaja
  documento.Activo = documentoExistente.Activo
  documento.Borrado = false

  // Modifico
  // ********
  //----------------------------------Modificar ######
  estado, valor, mensaje, httpStat = PersonalModificar(documentoID, documento, req, audit)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }

  // Está todo Ok
  // ************
  //------------------------------------Modificar ######
  s := []string{"Recuperaste ", documento.Apellido}
  //--------------------------------------Modificar ######
  core.RspMsgJSON(w, req, "OK", documento.Apellido, strings.Join(s, ""), http.StatusAccepted)
  return
}

// Devuelve Estado, Valor, Mensaje, HttpStat, Collection, Existía
func PersonalModificar(documentoID bson.ObjectId, documentoModi models.Personal, req *http.Request, audit string) (string, string, string, int) {
  //-------------------Modificar ###### las 2 variables
  camposVacios := "No podés dejar vacío apellido y nombre"
  coll := config.DB_Personal
  empresaID := context.Get(req, "Empresa_id").(bson.ObjectId)

  // Verifico los campos obligatorios
  // ********************************
  //---------------Modificar ######
  if documentoModi.Apellido == "" || documentoModi.Nombre == "" {
    s := []string{camposVacios}
    return "ERROR", "Alta", strings.Join(s, ""), http.StatusNonAuthoritativeInfo
  }

  // Me fijo si ya Existe la clave única
  // ***********************************
  //------------------------------------------------------Modificar ######-------------Modificar ######
  estado, valor, mensaje, httpStat, documentoExiste, _ := PersonalExiste(documentoModi.Apellido, documentoModi.Nombre, req)
  if httpStat == http.StatusInternalServerError {
    return estado, valor, mensaje, httpStat
  }
  if httpStat != http.StatusOK && documentoExiste.ID != documentoID {
    return estado, valor, mensaje, httpStat
  }

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    return "ERROR", "GetMongoSession", err.Error(), httpStat
  }
  defer session.Close()

  // Intento la modificación
  // ***********************
  //-------------------Modificar ###### en forma manual
  documentoModi.ID = documentoID
  documentoModi.Empresa_id = empresaID
  collection := session.DB(config.DB_Name).C(coll)
  selector := bson.M{"_id": documentoID, "empresa_id": empresaID}
  updator := bson.M{
    "$set": bson.M{
      "apellido": documentoModi.Apellido,
      "nombre": documentoModi.Nombre,
      "celular": documentoModi.Celular,
      "direccion": documentoModi.Direccion,
      "latitud": documentoModi.Latitud,
      "longitud": documentoModi.Longitud,
      "categoria_id": documentoModi.Categoria_id,
      "propio": documentoModi.Propio,
      "basicoSindicato_id": documentoModi.BasicoSindicato_id,
      "comision": documentoModi.Comision,
      "curso": documentoModi.Curso,
      "lnh": documentoModi.LNH,
      "registro": documentoModi.Registro,
      "fechaAlta": documentoModi.FechaAlta,
      "fechaBaja": documentoModi.FechaBaja,
      "activo": documentoModi.Activo,
      "borrado": documentoModi.Borrado,
      "timestamp": time.Now(),
    },
  }
  err = collection.Update(selector, updator)
  if err != nil {
    s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
    return "ERROR", audit, strings.Join(s, ""), http.StatusInternalServerError
  }

  // Está todo Ok
  // ************
  core.Audit(req, coll, documentoID, audit, documentoModi)
  return "OK", audit, "Ok", http.StatusOK
}
