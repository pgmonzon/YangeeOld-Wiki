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
	router.HandleFunc("/tokenCliente", handlers.TokenCliente).Methods("POST")
	router.HandleFunc("/test", handlers.ValidarMiddleware(handlers.TestEndpoint)).Methods("GET")

	// Usuario
	router.HandleFunc("/usuario", handlers.ValidarMiddleware(handlers.UsuarioRegistrar)).Methods("POST")

	// ClienteAPI
	router.HandleFunc("/clienteAPI/alta", handlers.ClienteAPIAlta).Methods("POST")

	// RBAC
	router.HandleFunc("/permiso", handlers.ValidarMiddleware(handlers.PermisoAgregar)).Methods("POST")
	router.HandleFunc("/rol", handlers.ValidarMiddleware(handlers.RolAgregar)).Methods("POST")

	log.Fatal(http.ListenAndServe(":3113", router))
}
