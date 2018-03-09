package routers

import (
	"net/http"
  "log"

  "github.com/pgmonzon/Yangee/handlers"

  "github.com/gorilla/mux"
)

func InicializarRutas() {
  router := mux.NewRouter()

	// *** AUT *** ### Módulo de Autorización #############
  // AUT - Autorización: Genera el token para operar
	router.HandleFunc("/autorizar", handlers.ValidarMiddleware(handlers.Autorizar, "NO_VALIDAR")).Methods("POST")
	// AUT - tokenCliente: Es el token que debe generar el cliente, es a los efectos de ejemplo
	router.HandleFunc("/tokenCliente", handlers.ValidarMiddleware(handlers.TokenCliente, "NO_VALIDAR")).Methods("POST")

	// *** FILOSOFOS ***
	router.HandleFunc("/filosofo", handlers.ValidarMiddleware(handlers.CrearFilosofo, "CrearFilosofo")).Methods("POST")

	// Usuario
	//router.HandleFunc("/usuario", handlers.ValidarMiddleware(handlers.UsuarioRegistrar, "UsuarioRegistrar")).Methods("POST")
	router.HandleFunc("/usuario", handlers.ValidarMiddleware(handlers.UsuarioRegistrar, "NO_VALIDAR")).Methods("POST")
	//router.HandleFunc("/usuarioPermisos", handlers.ValidarMiddleware(handlers.UsuarioPermisos)).Methods("GET")

	// ClienteAPI
	router.HandleFunc("/clienteAPI/alta", handlers.ValidarMiddleware(handlers.ClienteAPIAlta, "NO_VALIDAR")).Methods("POST")

	// RBAC
	//router.HandleFunc("/permiso", handlers.ValidarMiddleware(handlers.PermisoAgregar, "PermisoAgregar")).Methods("POST")
	//router.HandleFunc("/rol", handlers.ValidarMiddleware(handlers.RolAgregar, "RolAgregar")).Methods("POST")
	router.HandleFunc("/permiso", handlers.ValidarMiddleware(handlers.PermisoAgregar, "NO_VALIDAR")).Methods("POST")
	router.HandleFunc("/rol", handlers.ValidarMiddleware(handlers.RolAgregar, "NO_VALIDAR")).Methods("POST")

	// Revisar para sacar
	router.HandleFunc("/testPermisos", handlers.TestPermisos).Methods("POST")

	log.Fatal(http.ListenAndServe(":3113", router))
}
