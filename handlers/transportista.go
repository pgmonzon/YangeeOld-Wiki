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
reemplazar Transportistas (mayúscula plural) 4 apariciones
reemplazar Transportista (mayúscula singular) 76 apariciones
reemplazar transportista (minúscula singular) 7 apariciones IMPORANTE: no puede tener mayúsculas

*** Reemplazos manuales ***
reemplazar "No podés dejar vacío" 3 apariciones
reemplazar "en forma manual" 7 apariciones
reemplazar "en orden" 2 apariciones
*/

func TransportistaCrear(w http.ResponseWriter, req *http.Request) {
  //-------------------Modificar ###### estas 2 variables
	var documento models.Transportista
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
  estado, valor, mensaje, httpStat, documento, existia := TransportistaAlta(documento, req, audit)
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
  s := []string{"Agregaste ", documento.Transportista}
  //--------------------------------------Modificar ######
  core.RspMsgJSON(w, req, "OK", documento.Transportista, strings.Join(s, ""), http.StatusCreated)
  return
}

// Devuelve Estado, Valor, Mensaje, HttpStat, Collection, Existía
func TransportistaAlta(documentoAlta models.Transportista, req *http.Request, audit string) (string, string, string, int, models.Transportista, bool) {
  //-------------------Modificar ###### las 3 variables
	var documento models.Transportista
  camposVacios := "No podés dejar vacío el campo Transportista"
  coll := config.DB_Transportista
  empresaID := context.Get(req, "Empresa_id").(bson.ObjectId)

  // Verifico los campos obligatorios
  // ********************************
  //---------------Modificar ######
  if documentoAlta.Transportista == "" {
    s := []string{camposVacios}
    return "ERROR", "Alta", strings.Join(s, ""), http.StatusNonAuthoritativeInfo, documento, false
  }

  // Me fijo si ya Existe
  // ********************
  //-----------------------------------------------------Modificar ######--------------Modificar ######
  estado, valor, mensaje, httpStat, documento, existia := TransportistaExiste(documentoAlta.Transportista, req)
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
func TransportistaExiste(documentoExiste string, req *http.Request) (string, string, string, int, models.Transportista, bool) {
  //-------------------Modificar ###### las 3 variables
  var documento models.Transportista
  indice := []string{"empresa_id", "transportista"}
  coll := config.DB_Transportista
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
  collection.Find(bson.M{"empresa_id": empresaID, "transportista": documentoExiste}).One(&documento)
  // No existe
  if documento.ID == "" {
    return "OK", "Buscar", "Ok", http.StatusOK, documento, false
  }
  // Existe borrado
  if documento.Borrado == true {
    s := []string{documentoExiste," existe borrado"}
    return "ERROR", "Buscar", strings.Join(s, ""), http.StatusNonAuthoritativeInfo, documento, true
  }
  // Existe inactivo
  if documento.Activo == false {
    s := []string{documentoExiste," existe inactivo"}
    return "ERROR", "Buscar", strings.Join(s, ""), http.StatusNonAuthoritativeInfo, documento, true
  }
  // Existe
  s := []string{documentoExiste," ya existe"}
  return "ERROR", "Buscar", strings.Join(s, ""), http.StatusNonAuthoritativeInfo, documento, true
}

func TransportistasTraer(w http.ResponseWriter, req *http.Request) {
  //-------------------Modificar ###### estas 2 variables
  var documento models.Transportista
  var documentos []models.Transportista
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
  estado, valor, mensaje, httpStat, documentos := TransportistasBuscar(documento, orden, limiteInt, false, "Buscar", req)
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
func TransportistasBuscar(documento models.Transportista, orden string, limiteInt int, borrados bool, audit string, req *http.Request) (string, string, string, int, []models.Transportista) {
  //----------------------Modificar ###### estas 2 variables
  var documentos []models.Transportista
  coll := config.DB_Transportista
  empresaID := context.Get(req, "Empresa_id").(bson.ObjectId)

  // Verifico que el campo orden sea Unique
  // **************************************
  //-----------Modificar ######
  if orden != "transportista" && orden != "-transportista" {
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
    "transportista": bson.M{"$regex": bson.RegEx{documento.Transportista, "i"}},
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

func TransportistaTraer(w http.ResponseWriter, req *http.Request) {
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
  estado, valor, mensaje, httpStat, documento := Transportista_X_ID(documentoID, "Buscar ID", req)
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
func Transportista_X_ID(documentoID bson.ObjectId, audit string, req *http.Request) (string, string, string, int, models.Transportista) {
  //-------------------Modificar ######
  var documento models.Transportista
  coll := config.DB_Transportista
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

func TransportistaGuardar(w http.ResponseWriter, req *http.Request) {
  //-------------------Modificar ######
  var documento models.Transportista
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
  estado, valor, mensaje, httpStat, documentoExistente := Transportista_X_ID(documentoID, "Buscar ID", req)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }
  documento.Borrado = documentoExistente.Borrado

  // Modifico
  // ********
  //----------------------------------Modificar ######
  estado, valor, mensaje, httpStat = TransportistaModificar(documentoID, documento, req, audit)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }

  // Está todo Ok
  // ************
  //------------------------------------Modificar ######
  s := []string{"Guardaste ", documento.Transportista}
  //--------------------------------------Modificar ######
  core.RspMsgJSON(w, req, "OK", documento.Transportista, strings.Join(s, ""), http.StatusAccepted)
  return
}

func TransportistaHabilitar(w http.ResponseWriter, req *http.Request) {
  //-------------------Modificar ######
  var documento models.Transportista
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
  estado, valor, mensaje, httpStat, documentoExistente := Transportista_X_ID(documentoID, "Buscar ID", req)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }
  //-------Modificar ###### en forma manual
  documento.Transportista = documentoExistente.Transportista
  documento.Mail = documentoExistente.Mail
  documento.Contactos = documentoExistente.Contactos
  documento.Tarifarios = documentoExistente.Tarifarios
  documento.Activo = true
  documento.Borrado = documentoExistente.Borrado

  // Modifico
  // ********
  //----------------------------------Modificar ######
  estado, valor, mensaje, httpStat = TransportistaModificar(documentoID, documento, req, audit)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }

  // Está todo Ok
  // ************
  //------------------------------------Modificar ######
  s := []string{"Habilitaste ", documento.Transportista}
  //--------------------------------------Modificar ######
  core.RspMsgJSON(w, req, "OK", documento.Transportista, strings.Join(s, ""), http.StatusAccepted)
  return
}

func TransportistaDeshabilitar(w http.ResponseWriter, req *http.Request) {
  //-------------------Modificar ######
  var documento models.Transportista
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
  estado, valor, mensaje, httpStat, documentoExistente := Transportista_X_ID(documentoID, "Buscar ID", req)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }
  //-------Modificar ###### en forma manual
  documento.Transportista = documentoExistente.Transportista
  documento.Mail = documentoExistente.Mail
  documento.Contactos = documentoExistente.Contactos
  documento.Tarifarios = documentoExistente.Tarifarios
  documento.Activo = false
  documento.Borrado = documentoExistente.Borrado

  // Modifico
  // ********
  //----------------------------------Modificar ######
  estado, valor, mensaje, httpStat = TransportistaModificar(documentoID, documento, req, audit)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }

  // Está todo Ok
  // ************
  //------------------------------------Modificar ######
  s := []string{"Deshabilitaste ", documento.Transportista}
  //--------------------------------------Modificar ######
  core.RspMsgJSON(w, req, "OK", documento.Transportista, strings.Join(s, ""), http.StatusAccepted)
  return
}

func TransportistaBorrar(w http.ResponseWriter, req *http.Request) {
  //-------------------Modificar ######
  var documento models.Transportista
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
  estado, valor, mensaje, httpStat, documentoExistente := Transportista_X_ID(documentoID, "Buscar ID", req)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }
  //-------Modificar ###### en forma manual
  documento.Transportista = documentoExistente.Transportista
  documento.Mail = documentoExistente.Mail
  documento.Contactos = documentoExistente.Contactos
  documento.Tarifarios = documentoExistente.Tarifarios
  documento.Activo = documentoExistente.Activo
  documento.Borrado = true

  // Modifico
  // ********
  //----------------------------------Modificar ######
  estado, valor, mensaje, httpStat = TransportistaModificar(documentoID, documento, req, audit)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }

  // Está todo Ok
  // ************
  //------------------------------------Modificar ######
  s := []string{"Borraste ", documento.Transportista}
  //--------------------------------------Modificar ######
  core.RspMsgJSON(w, req, "OK", documento.Transportista, strings.Join(s, ""), http.StatusAccepted)
  return
}

func TransportistaRecuperar(w http.ResponseWriter, req *http.Request) {
  //-------------------Modificar ######
  var documento models.Transportista
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
  estado, valor, mensaje, httpStat, documentoExistente := Transportista_X_ID(documentoID, "Buscar ID", req)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }
  //-------Modificar ###### en forma manual
  documento.Transportista = documentoExistente.Transportista
  documento.Mail = documentoExistente.Mail
  documento.Contactos = documentoExistente.Contactos
  documento.Tarifarios = documentoExistente.Tarifarios
  documento.Activo = documentoExistente.Activo
  documento.Borrado = false

  // Modifico
  // ********
  //----------------------------------Modificar ######
  estado, valor, mensaje, httpStat = TransportistaModificar(documentoID, documento, req, audit)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }

  // Está todo Ok
  // ************
  //------------------------------------Modificar ######
  s := []string{"Recuperaste ", documento.Transportista}
  //--------------------------------------Modificar ######
  core.RspMsgJSON(w, req, "OK", documento.Transportista, strings.Join(s, ""), http.StatusAccepted)
  return
}

// Devuelve Estado, Valor, Mensaje, HttpStat, Collection, Existía
func TransportistaModificar(documentoID bson.ObjectId, documentoModi models.Transportista, req *http.Request, audit string) (string, string, string, int) {
  //-------------------Modificar ###### las 2 variables
  camposVacios := "No podés dejar vacío el campo Transportista"
  coll := config.DB_Transportista
  empresaID := context.Get(req, "Empresa_id").(bson.ObjectId)

  // Verifico los campos obligatorios
  // ********************************
  //---------------Modificar ######
  if documentoModi.Transportista == "" {
    s := []string{camposVacios}
    return "ERROR", "Alta", strings.Join(s, ""), http.StatusNonAuthoritativeInfo
  }

  // Me fijo si ya Existe la clave única
  // ***********************************
  //------------------------------------------------------Modificar ######-------------Modificar ######
  estado, valor, mensaje, httpStat, documentoExiste, _ := TransportistaExiste(documentoModi.Transportista, req)
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
      "transportista": documentoModi.Transportista,
      "mail": documentoModi.Mail,
      "contactos": documentoModi.Contactos,
      "tarifarios": documentoModi.Tarifarios,
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
