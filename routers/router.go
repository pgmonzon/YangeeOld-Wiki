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
	router.HandleFunc("/test", handlers.ValidarMiddleware(handlers.TestEndpoint, "TestEndPoint")).Methods("GET")
	router.HandleFunc("/testPermisos", handlers.TestPermisos).Methods("POST")

	// Usuario
	router.HandleFunc("/usuario", handlers.ValidarMiddleware(handlers.UsuarioRegistrar, "UsuarioRegistrar")).Methods("POST")
	//router.HandleFunc("/usuarioPermisos", handlers.ValidarMiddleware(handlers.UsuarioPermisos)).Methods("GET")

	// ClienteAPI
	router.HandleFunc("/clienteAPI/alta", handlers.ClienteAPIAlta).Methods("POST")

	// RBAC
	router.HandleFunc("/permiso", handlers.ValidarMiddleware(handlers.PermisoAgregar, "PermisoAgregar")).Methods("POST")
	router.HandleFunc("/rol", handlers.ValidarMiddleware(handlers.RolAgregar, "RolAgregar")).Methods("POST")

	log.Fatal(http.ListenAndServe(":3113", router))
}
