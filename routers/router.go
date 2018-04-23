package routers

import (
	"net/http"
  "log"

  "github.com/pgmonzon/Yangee/handlers"

  "github.com/gorilla/mux"
)

func InicializarRutas() {
  router := mux.NewRouter()

	// ClientesAPI
	// ***********
	router.HandleFunc("/clienteAPI", handlers.ValidarMiddleware(handlers.ClienteAPICrear, "NO_VALIDAR")).Methods("POST")

	// Autorización
	// ************
	router.HandleFunc("/modulo", handlers.ValidarMiddleware(handlers.ModuloCrear, "NO_VALIDAR")).Methods("POST")
	router.HandleFunc("/autorizar", handlers.ValidarMiddleware(handlers.Autorizar, "NO_VALIDAR")).Methods("POST")
	router.HandleFunc("/tokenCliente", handlers.ValidarMiddleware(handlers.TokenCliente, "NO_VALIDAR")).Methods("POST")
	//router.HandleFunc("/invitacionEmpresa", handlers.ValidarMiddleware(handlers.InvitacionEmpresa, "NO_VALIDAR")).Methods("POST")

	// Filosofos
	// *********
	router.HandleFunc("/filosofo", handlers.ValidarMiddleware(handlers.FilosofoCrear, "FilosofoCrear")).Methods("POST")
	router.HandleFunc("/filosofo/{filosofoID}", handlers.ValidarMiddleware(handlers.FilosofoModificar, "FilosofoModificar")).Methods("PUT")
	router.HandleFunc("/filosofo/{filosofoID}", handlers.ValidarMiddleware(handlers.FilosofoBorrar, "FilosofoBorrar")).Methods("DELETE")
	router.HandleFunc("/filosofo/{filosofoID}", handlers.ValidarMiddleware(handlers.FilosofoTraer, "FilosofoTraer")).Methods("GET")
	router.HandleFunc("/filosofos/{orden}/{limite}", handlers.ValidarMiddleware(handlers.FilosofosTraer, "FilosofosTraer")).Methods("POST")
	router.HandleFunc("/filosofosSiguiente/{orden}/{limite}/{ultimo_campo_orden}", handlers.ValidarMiddleware(handlers.FilosofosTraerSiguiente, "FilosofosTraerSiguiente")).Methods("POST")
	router.HandleFunc("/filosofosAnterior/{orden}/{limite}/{primer_campo_orden}", handlers.ValidarMiddleware(handlers.FilosofosTraerAnterior, "FilosofosTraerAnterior")).Methods("POST")

	// Usuario
	// *******
	//router.HandleFunc("/usuario", handlers.ValidarMiddleware(handlers.UsuarioRegistrar, "UsuarioRegistrar")).Methods("POST")
	//router.HandleFunc("/usuario", handlers.ValidarMiddleware(handlers.UsuarioRegistrar, "NO_VALIDAR")).Methods("POST")
	//router.HandleFunc("/usuarioPermisos", handlers.ValidarMiddleware(handlers.UsuarioPermisos)).Methods("GET")

	// RBAC
	//router.HandleFunc("/permiso", handlers.ValidarMiddleware(handlers.PermisoAgregar, "PermisoAgregar")).Methods("POST")
	//router.HandleFunc("/rol", handlers.ValidarMiddleware(handlers.RolAgregar, "RolAgregar")).Methods("POST")
	//router.HandleFunc("/permiso", handlers.ValidarMiddleware(handlers.PermisoAgregar, "NO_VALIDAR")).Methods("POST")
	//router.HandleFunc("/rol", handlers.ValidarMiddleware(handlers.RolAgregar, "NO_VALIDAR")).Methods("POST")

	log.Fatal(http.ListenAndServe(":3113", router))
}
