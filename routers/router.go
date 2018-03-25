package routers

import (
	"net/http"
  "log"

  "github.com/pgmonzon/Yangee/handlers"

  "github.com/gorilla/mux"
)

func InicializarRutas() {
  router := mux.NewRouter()

	// Autorización
	// ************
	router.HandleFunc("/autorizar", handlers.ValidarMiddleware(handlers.Autorizar, "NO_VALIDAR")).Methods("POST")
	//router.HandleFunc("/tokenCliente", handlers.ValidarMiddleware(handlers.TokenCliente, "NO_VALIDAR")).Methods("POST")

	// Filosofos
	// *********
	router.HandleFunc("/filosofo", handlers.ValidarMiddleware(handlers.FilosofoCrear, "FilosofoCrear")).Methods("POST")
	router.HandleFunc("/filosofo/{filosofoID}", handlers.ValidarMiddleware(handlers.FilosofoModificar, "FilosofoModificar")).Methods("PUT")
	router.HandleFunc("/filosofo/{filosofoID}", handlers.ValidarMiddleware(handlers.FilosofoBorrar, "FilosofoBorrar")).Methods("DELETE")

	// Usuario
	// *******
	//router.HandleFunc("/usuario", handlers.ValidarMiddleware(handlers.UsuarioRegistrar, "UsuarioRegistrar")).Methods("POST")
	//router.HandleFunc("/usuario", handlers.ValidarMiddleware(handlers.UsuarioRegistrar, "NO_VALIDAR")).Methods("POST")
	//router.HandleFunc("/usuarioPermisos", handlers.ValidarMiddleware(handlers.UsuarioPermisos)).Methods("GET")

	// ClientesAPI
	// ***********
	router.HandleFunc("/clienteAPI", handlers.ValidarMiddleware(handlers.CrearClienteAPI, "NO_VALIDAR")).Methods("POST")

	// RBAC
	//router.HandleFunc("/permiso", handlers.ValidarMiddleware(handlers.PermisoAgregar, "PermisoAgregar")).Methods("POST")
	//router.HandleFunc("/rol", handlers.ValidarMiddleware(handlers.RolAgregar, "RolAgregar")).Methods("POST")
	//router.HandleFunc("/permiso", handlers.ValidarMiddleware(handlers.PermisoAgregar, "NO_VALIDAR")).Methods("POST")
	//router.HandleFunc("/rol", handlers.ValidarMiddleware(handlers.RolAgregar, "NO_VALIDAR")).Methods("POST")

	// Revisar para sacar
	//router.HandleFunc("/testPermisos", handlers.TestPermisos).Methods("POST")

	log.Fatal(http.ListenAndServe(":3113", router))
}
