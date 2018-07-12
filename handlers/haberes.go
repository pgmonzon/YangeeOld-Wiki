package handlers

import (
  "encoding/json"
  "net/http"
  "strings"
  "time"

  "github.com/pgmonzon/Yangee/models"
  "github.com/pgmonzon/Yangee/core"
  "github.com/pgmonzon/Yangee/config"

  "gopkg.in/mgo.v2/bson"
  "github.com/gorilla/context"
)

func HaberesCrear(w http.ResponseWriter, req *http.Request) {
	var documento models.Haberes
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
  estado, valor, mensaje, httpStat, documento := HaberesAlta(documento, req, audit)
  if httpStat != http.StatusOK {
    core.RspMsgJSON(w, req, estado, valor, mensaje, httpStat)
    return
  }

  // Está todo Ok
  // ************
  s := []string{"Agregaste haberes"}
  core.RspMsgJSON(w, req, "OK", "Haberes", strings.Join(s, ""), http.StatusCreated)
  return
}

func HaberesAlta(documentoAlta models.Haberes, req *http.Request, audit string) (string, string, string, int, models.Haberes) {
	var documento models.Haberes
  var basicoSindicato float64
  coll := config.DB_Haberes
  empresaID := context.Get(req, "Empresa_id").(bson.ObjectId)

  // Genero una nueva sesión Mongo
  // *****************************
  session, err, httpStat := core.GetMongoSession()
  if err != nil {
    return "ERROR", "GetMongoSession", err.Error(), httpStat, documento
  }
  defer session.Close()

  // Intento el alta
  // ***************
  documento = documentoAlta
  objID := bson.NewObjectId()
  documento.ID = objID
  documento.Empresa_id = empresaID
  documento.Timestamp = time.Now()
  documento.Editable = true

  collection := session.DB(config.DB_Name).C(coll)
  err = collection.Insert(documento)
  if err != nil {
    s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
    return "ERROR", audit, strings.Join(s, ""), http.StatusInternalServerError, documento
  }

  // Genero las Novedades
  // ********************
  personal := make([]models.Personal, 0)
  collPersonal := session.DB(config.DB_Name).C(config.DB_Personal)
  selector := bson.M{
    "empresa_id": empresaID,
    "propio": true,
    "activo": true,
    "borrado": false,
  }
  collPersonal.Find(selector).Select(bson.M{"empresa_id":0}).Sort("categoria_id").All(&personal)

  var nov models.Novedades
  cat_id := config.FakeID
  for _, item := range personal {
    if item.Categoria_id != cat_id {
      cat_id = item.Categoria_id
    }

    nombre := []string{item.Nombre, " ", item.Apellido}
    viajesMes := ViajesPersonalComisiones(item.ID, empresaID, documento.ComisionesDesde, documento.ComisionesHasta)
    basicoSindicato = 0
    for _, sindicatos := range documento.BasicosSindicato {
      if sindicatos.Basico_id == item.BasicoSindicato_id {
        basicoSindicato = sindicatos.Importe
      }
    }
    comisionEstimada := viajesMes * float64(item.Comision) / 100
    anticiposPendientes := RendicionSaldoAnterior(item.ID, req)
    if anticiposPendientes < 0 {
      anticiposPendientes = 0
    }

    objID = bson.NewObjectId()
    nov.ID = objID
    nov.Empresa_id = empresaID
  	nov.Haberes_id = documento.ID
    nov.Personal_id = item.ID
  	nov.Personal = strings.Join(nombre, "")
    nov.BasicoSindicato = basicoSindicato
  	nov.ViajesMes = viajesMes
  	nov.Comision = item.Comision
  	nov.ComisionEstimada = comisionEstimada
  	nov.Diferencia = comisionEstimada - basicoSindicato
  	nov.LiquidacionFinal = basicoSindicato
  	nov.ComisionReal = 0
  	nov.AnticiposPendientes = anticiposPendientes
  	nov.AnticiposAplicados = 0
  	nov.NetoPagar = basicoSindicato
  	nov.PagoBanco = 0
  	nov.PagoEfectivo = 0
  	nov.Pendiente = basicoSindicato

    HaberesNovedadesAlta(nov)
  }

  // Está todo Ok
  // ************
  core.Audit(req, coll, documento.ID, audit, documento)
  return "OK", audit, "Ok", http.StatusOK, documento
}

func HaberesNovedadesAlta(documentoAlta models.Novedades) {
	var documento models.Novedades
  coll := config.DB_Novedades

  // Genero una nueva sesión Mongo
  // *****************************
  session, _, _ := core.GetMongoSession()
  defer session.Close()

  // Intento el alta
  // ***************
  documento = documentoAlta
  objID := bson.NewObjectId()
  documento.ID = objID

  collection := session.DB(config.DB_Name).C(coll)
  collection.Insert(documento)

  // Está todo Ok
  // ************
  return
}
