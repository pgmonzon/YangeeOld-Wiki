package handlers

import (
  "fmt"
  "net/http"
  "encoding/json"

  "github.com/pgmonzon/Yangee/models"
)

// Valida las credenciales del usuario
func Autorizar(w http.ResponseWriter, req *http.Request) {
  fmt.Println("Autorizar")

  // Verifico BAD_REQUEST
  var aut models.Autorizar
  err := json.NewDecoder(req.Body).Decode(&aut)
  if err != nil || aut.Usuario == "" || aut.Clave == "" {
    fmt.Println("mal request")
  }
}

/**
// Valida las credenciales del usuario
func UsuarioLogin(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

  // Genero una nueva sesión Mongo
	session := core.GetMongoSession()
	defer session.Close()

  // Intento traer el Usuario
	collection := session.DB("yangee").C("usuario")
	collection.Find(bson.M{"usuario": usuarioLogin.Usuario, "clave": core.HashSha512(usuarioLogin.Clave)}).One(&usuario)
	if usuario.ID == "" {
		core.ErrorJSON(w, r, start, "Acceso denegado", http.StatusUnauthorized)
	} else {
    token, err := core.CrearToken(usuario)
    if err != nil {
      core.ErrorJSON(w, r, start, token, http.StatusInternalServerError)
    }
    response, err := json.Marshal(models.Token{token})
		core.FatalErr(err)
		core.RespuestaJSON(w, r, start, response, http.StatusOK)
	}
}
**/
