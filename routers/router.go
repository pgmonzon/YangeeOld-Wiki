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
	router.HandleFunc("/empresa", handlers.ValidarMiddleware(handlers.EmpresaCrear, "NO_VALIDAR")).Methods("POST")
	router.HandleFunc("/autorizar", handlers.ValidarMiddleware(handlers.Autorizar, "NO_VALIDAR")).Methods("OPTIONS")
	router.HandleFunc("/tokenCliente", handlers.ValidarMiddleware(handlers.TokenCliente, "NO_VALIDAR")).Methods("POST")
	router.HandleFunc("/empresaInvitacion", handlers.ValidarMiddleware(handlers.EmpresaInvitar, "NO_VALIDAR")).Methods("POST")

	// RBAC
	// ****
	router.HandleFunc("/permiso", handlers.ValidarMiddleware(handlers.PermisoCrear, "NO_VALIDAR")).Methods("POST")
	router.HandleFunc("/rol", handlers.ValidarMiddleware(handlers.RolCrear, "NO_VALIDAR")).Methods("POST")
	router.HandleFunc("/usuario", handlers.ValidarMiddleware(handlers.UsuarioCrear, "NO_VALIDAR")).Methods("POST")

	// Filosofos
	// *********
	router.HandleFunc("/filosofo", handlers.ValidarMiddleware(handlers.FilosofoCrear, "FilosofoCrear")).Methods("POST")
	router.HandleFunc("/filosofo/{filosofoID}", handlers.ValidarMiddleware(handlers.FilosofoModificar, "FilosofoModificar")).Methods("PUT")
	router.HandleFunc("/filosofo/{filosofoID}", handlers.ValidarMiddleware(handlers.FilosofoBorrar, "FilosofoBorrar")).Methods("DELETE")
	router.HandleFunc("/filosofo/{filosofoID}", handlers.ValidarMiddleware(handlers.FilosofoTraer, "FilosofoTraer")).Methods("GET")
	router.HandleFunc("/filosofos/{orden}/{limite}", handlers.ValidarMiddleware(handlers.FilosofosTraer, "FilosofosTraer")).Methods("POST")
	router.HandleFunc("/filosofosSiguiente/{orden}/{limite}/{ultimo_campo_orden}", handlers.ValidarMiddleware(handlers.FilosofosTraerSiguiente, "FilosofosTraerSiguiente")).Methods("POST")
	router.HandleFunc("/filosofosAnterior/{orden}/{limite}/{primer_campo_orden}", handlers.ValidarMiddleware(handlers.FilosofosTraerAnterior, "FilosofosTraerAnterior")).Methods("POST")

	log.Fatal(http.ListenAndServe(":3113", router))
}
