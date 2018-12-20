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

func RolCrear(w http.ResponseWriter, req *http.Request) {
	var documento models.RolEstado
  var docCreado models.Rol
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
  estado, valor, mensaje, httpStat, docCreado, existia := RolAlta(documento, req, audit)
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
  s := []string{"Agregaste ", docCreado.Rol}
  core.RspMsgJSON(w, req, "OK", docCreado.Rol, strings.Join(s, ""), http.StatusCreated)
  return
}

func RolAlta(documentoAlta models.RolEstado, req *http.Request, audit string) (string, string, string, int, models.Rol, bool) {
	var documento models.Rol
  camposVacios := "No podés dejar vacío el campo Rol"
  coll := config.DB_Rol
  empresaID := context.Get(req, "Empresa_id").(bson.ObjectId)

  // Verifico los campos obligatorios
  // ********************************
  if documentoAlta.Rol == "" {
    s := []string{camposVacios}
    return "ERROR", "Alta", strings.Join(s, ""), http.StatusNonAuthoritativeInfo, documento, false
  }

  // Me fijo si ya Existe
  // ********************
  estado, valor, mensaje, httpStat, documento, existia := RolExiste(documentoAlta.Rol, req)
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
  objID := bson.NewObjectId()
  documento.ID = objID
  documento.Empresa_id = empresaID
  documento.Rol = documentoAlta.Rol
  documento.Activo = documentoAlta.Activo
  documento.Timestamp = time.Now()
  documento.Borrado = false
  var menu []models.Opcion
  for _, item := range documentoAlta.Menu {
    if item.Estado == true {
      var subMenu []models.Sub
      for _, subItem := range item.Children {
        if subItem.Estado == true {
          subMenu = append(subMenu, models.Sub{Path: subItem.Path, Title: subItem.Title, Ab: subItem.Ab})
        }
      }
      menu = append(menu, models.Opcion{Path: item.Path, Type: item.Type, Title: item.Title, Icontype: item.Icontype, Collapse: item.Collapse, Children: subMenu})
    }
  }
  documento.Menu = menu

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

func RolExiste(documentoExiste string, req *http.Request) (string, string, string, int, models.Rol, bool) {
  var documento models.Rol
  indice := []string{"empresa_id", "rol"}
  coll := config.DB_Rol
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
  collection.Find(bson.M{"empresa_id": empresaID, "rol": documentoExiste}).One(&documento)
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

func RolesTraer(w http.ResponseWriter, req *http.Request) {
  var documento models.Rol
  var documentos []models.Rol
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
  estado, valor, mensaje, httpStat, documentos := RolesBuscar(documento, orden, limiteInt, false, "Buscar", req)
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

func RolesBuscar(documento models.Rol, orden string, limiteInt int, borrados bool, audit string, req *http.Request) (string, string, string, int, []models.Rol) {
  var documentos []models.Rol
  coll := config.DB_Rol
  empresaID := context.Get(req, "Empresa_id").(bson.ObjectId)

  // Verifico que el campo orden sea Unique
  // **************************************
  if orden != "rol" && orden != "-rol" {
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
  selector := bson.M{
    "empresa_id": empresaID,
    "rol": bson.M{"$regex": bson.RegEx{documento.Rol, "i"}},
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

func RolTraer(w http.ResponseWriter, req *http.Request) {
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
  estado, valor, mensaje, httpStat, documento := Rol_X_ID(documentoID, "Buscar ID", req)
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

func Rol_X_ID(documentoID bson.ObjectId, audit string, req *http.Request) (string, string, string, int, models.RolEstado) {
  var documento models.Rol
  var documentoEstado models.RolEstado
  coll := config.DB_Rol
  empresaID := context.Get(req, "Empresa_id").(bson.ObjectId)

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    return "ERROR", "GetMongoSession", err.Error(), httpStat, documentoEstado
  }
  defer session.Close()

  // Trato de traerlo
  // ****************
  collection := session.DB(config.DB_Name).C(coll)
  collection.Find(bson.M{"_id": documentoID, "empresa_id": empresaID}).Select(bson.M{"empresa_id":0}).One(&documento)
  // No existe
  if documento.ID == "" {
    s := []string{"No encuentro el documento"}
    return "ERROR", audit, strings.Join(s, ""), http.StatusNonAuthoritativeInfo, documentoEstado
  }

  // Busco el menu de la empresa
  // ***************************
  estado, valor, mensaje, httpStat, docMenu := Menu_X_Empresa(empresaID, "Buscar Empresa", req)
  if httpStat != http.StatusOK {
    return estado, valor, mensaje, httpStat, documentoEstado
  }

  // Preparo el documento a devolver
  // *******************************
  documentoEstado.ID = documento.ID
  documentoEstado.Rol = documento.Rol
  documentoEstado.Activo = documento.Activo
  documentoEstado.Borrado = documento.Borrado
  documentoEstado.Timestamp = documento.Timestamp

  // Pongo los estados de los ítems del menú
  // ***************************************
  estadoOpcion := false
  estadoSub := false
  var menu []models.OpcionEstado
  for _, item := range docMenu.Menu {
    var subMenu []models.SubEstado
    for _, subItem := range item.Children {
      estadoOpcion = false
      estadoSub = false
      for _, itemEstado := range documento.Menu {
        if item.Path == itemEstado.Path {
          estadoOpcion = true
          for _, subItemEstado := range itemEstado.Children {
            if subItem.Path == subItemEstado.Path {
              estadoSub = true
            }
          }
        }
      }
      subMenu = append(subMenu, models.SubEstado{Path: subItem.Path, Title: subItem.Title, Ab: subItem.Ab, Estado: estadoSub})
    }
    menu = append(menu, models.OpcionEstado{Path: item.Path, Type: item.Type, Title: item.Title, Icontype: item.Icontype, Collapse: item.Collapse, Children: subMenu, Estado: estadoOpcion})
  }
  documentoEstado.Menu = menu

  // Existe
  return "OK", audit, "Ok", http.StatusOK, documentoEstado
}

func Rol_X_ID_SinEstado(documentoID bson.ObjectId, audit string, req *http.Request) (string, string, string, int, models.Rol) {
  var documento models.Rol
  coll := config.DB_Rol
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

func RolGuardar(w http.ResponseWriter, req *http.Request) {
  var documento models.RolEstado
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
  estado, valor, mensaje, httpStat, documentoExistente := Rol_X_ID(documentoID, "Buscar ID", req)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }
  documento.Borrado = documentoExistente.Borrado

  // Modifico
  // ********
  estado, valor, mensaje, httpStat = RolModificar(documentoID, documento, req, audit)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }

  // Está todo Ok
  // ************
  s := []string{"Guardaste ", documento.Rol}
  core.RspMsgJSON(w, req, "OK", documento.Rol, strings.Join(s, ""), http.StatusAccepted)
  return
}

func RolHabilitar(w http.ResponseWriter, req *http.Request) {
  var documento models.RolEstado
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
  estado, valor, mensaje, httpStat, documentoExistente := Rol_X_ID(documentoID, "Buscar ID", req)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }
  documento.Rol = documentoExistente.Rol
  documento.Menu = documentoExistente.Menu
  documento.Activo = true
  documento.Borrado = documentoExistente.Borrado

  // Modifico
  // ********
  estado, valor, mensaje, httpStat = RolModificar(documentoID, documento, req, audit)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }

  // Está todo Ok
  // ************
  s := []string{"Habilitaste ", documento.Rol}
  core.RspMsgJSON(w, req, "OK", documento.Rol, strings.Join(s, ""), http.StatusAccepted)
  return
}

func RolDeshabilitar(w http.ResponseWriter, req *http.Request) {
  var documento models.RolEstado
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
  estado, valor, mensaje, httpStat, documentoExistente := Rol_X_ID(documentoID, "Buscar ID", req)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }
  documento.Rol = documentoExistente.Rol
  documento.Menu = documentoExistente.Menu
  documento.Activo = false
  documento.Borrado = documentoExistente.Borrado

  // Modifico
  // ********
  estado, valor, mensaje, httpStat = RolModificar(documentoID, documento, req, audit)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }

  // Está todo Ok
  // ************
  s := []string{"Deshabilitaste ", documento.Rol}
  core.RspMsgJSON(w, req, "OK", documento.Rol, strings.Join(s, ""), http.StatusAccepted)
  return
}

func RolBorrar(w http.ResponseWriter, req *http.Request) {
  var documento models.RolEstado
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
  estado, valor, mensaje, httpStat, documentoExistente := Rol_X_ID(documentoID, "Buscar ID", req)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }
  documento.Rol = documentoExistente.Rol
  documento.Menu = documentoExistente.Menu
  documento.Activo = documentoExistente.Activo
  documento.Borrado = true

  // Modifico
  // ********
  estado, valor, mensaje, httpStat = RolModificar(documentoID, documento, req, audit)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }

  // Está todo Ok
  // ************
  s := []string{"Borraste ", documento.Rol}
  core.RspMsgJSON(w, req, "OK", documento.Rol, strings.Join(s, ""), http.StatusAccepted)
  return
}

func RolRecuperar(w http.ResponseWriter, req *http.Request) {
  var documento models.RolEstado
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
  estado, valor, mensaje, httpStat, documentoExistente := Rol_X_ID(documentoID, "Buscar ID", req)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }
  documento.Rol = documentoExistente.Rol
  documento.Menu = documentoExistente.Menu
  documento.Activo = documentoExistente.Activo
  documento.Borrado = false

  // Modifico
  // ********
  estado, valor, mensaje, httpStat = RolModificar(documentoID, documento, req, audit)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }

  // Está todo Ok
  // ************
  s := []string{"Recuperaste ", documento.Rol}
  core.RspMsgJSON(w, req, "OK", documento.Rol, strings.Join(s, ""), http.StatusAccepted)
  return
}

func RolModificar(documentoID bson.ObjectId, documentoModi models.RolEstado, req *http.Request, audit string) (string, string, string, int) {
  camposVacios := "No podés dejar vacío el campo Rol"
  coll := config.DB_Rol
  empresaID := context.Get(req, "Empresa_id").(bson.ObjectId)

  // Verifico los campos obligatorios
  // ********************************
  if documentoModi.Rol == "" {
    s := []string{camposVacios}
    return "ERROR", "Alta", strings.Join(s, ""), http.StatusNonAuthoritativeInfo
  }

  // Me fijo si ya Existe la clave única
  // ***********************************
  estado, valor, mensaje, httpStat, documentoExiste, _ := RolExiste(documentoModi.Rol, req)
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
  documentoModi.ID = documentoID
  documentoModi.Empresa_id = empresaID
  var menu []models.Opcion
  for _, item := range documentoModi.Menu {
    if item.Estado == true {
      var subMenu []models.Sub
      for _, subItem := range item.Children {
        if subItem.Estado == true {
          subMenu = append(subMenu, models.Sub{Path: subItem.Path, Title: subItem.Title, Ab: subItem.Ab})
        }
      }
      menu = append(menu, models.Opcion{Path: item.Path, Type: item.Type, Title: item.Title, Icontype: item.Icontype, Collapse: item.Collapse, Children: subMenu})
    }
  }

  collection := session.DB(config.DB_Name).C(coll)
  selector := bson.M{"_id": documentoID, "empresa_id": empresaID}
  updator := bson.M{
    "$set": bson.M{
      "rol": documentoModi.Rol,
      "menu": menu,
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

func RolMenuEmpresa(w http.ResponseWriter, req *http.Request) {
  empresaID := context.Get(req, "Empresa_id").(bson.ObjectId)

  // Busco el menu de la empresa
  // ***************************
  estado, valor, mensaje, httpStat, docMenu := Menu_X_Empresa(empresaID, "Buscar Menú Empresa", req)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }

  // Está todo Ok
  // ************
  respuesta, error := json.Marshal(docMenu)
  core.FatalErr(error)
  core.RspJSON(w, req, respuesta, http.StatusOK)
  return
}
