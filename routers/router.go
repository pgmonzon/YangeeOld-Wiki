package routers

import (
	"net/http"
  "log"

  "github.com/pgmonzon/Yangee/handlers"

  "github.com/gorilla/mux"
)

func InicializarRutas() {
  router := mux.NewRouter()

  // Autorización: Genera token para operar
	router.HandleFunc("/autorizar", handlers.Autorizar).Methods("POST")
	router.HandleFunc("/tokenAutorizar", handlers.TokenAutorizar).Methods("POST")

	// Usuario
	router.HandleFunc("/registrar", handlers.UsuarioRegistrar).Methods("POST")

	// ClienteAPI
	router.HandleFunc("/clienteAPI/alta", handlers.ClienteAPIAlta).Methods("POST")

	log.Fatal(http.ListenAndServe(":3113", router))
}
