package routers

import (
	"net/http"
  "log"

  "github.com/pgmonzon/Yangee/handlers"

	gorillaHnd "github.com/gorilla/handlers"
  "github.com/gorilla/mux"
)

func InicializarRutas() {
  router := mux.NewRouter()

	// ClientesAPI
	// ***********
	router.HandleFunc("/clienteAPI", handlers.ValidarMiddleware(handlers.ClienteAPICrear, "NO_VALIDAR")).Methods("POST")

	// Autorización
	// ************
	router.HandleFunc("/testPostBody", handlers.ValidarMiddleware(handlers.TestPostBody, "AUTH")).Methods("POST")
	router.HandleFunc("/testOptionsBody", handlers.ValidarMiddleware(handlers.TestOptionsBody, "AUTH")).Methods("OPTIONS")
	router.HandleFunc("/testPostHeader", handlers.ValidarMiddleware(handlers.TestPostHeader, "AUTH")).Methods("POST")
	router.HandleFunc("/testOptionsHeader", handlers.ValidarMiddleware(handlers.TestOptionsHeader, "AUTH")).Methods("OPTIONS")

	router.HandleFunc("/modulo", handlers.ValidarMiddleware(handlers.ModuloCrear, "NO_VALIDAR")).Methods("POST")
	router.HandleFunc("/empresa", handlers.ValidarMiddleware(handlers.EmpresaCrear, "NO_VALIDAR")).Methods("POST")
	router.HandleFunc("/autorizar", handlers.ValidarMiddleware(handlers.Autorizar, "AUTH")).Methods("POST")
	router.HandleFunc("/tokenCliente", handlers.ValidarMiddleware(handlers.TokenCliente, "AUTH")).Methods("POST")
	router.HandleFunc("/empresaInvitacion", handlers.ValidarMiddleware(handlers.EmpresaInvitar, "NO_VALIDAR")).Methods("POST")

	// RBAC
	// ****
	router.HandleFunc("/permiso", handlers.ValidarMiddleware(handlers.PermisoCrear, "NO_VALIDAR")).Methods("POST")
	router.HandleFunc("/rol", handlers.ValidarMiddleware(handlers.RolCrear, "NO_VALIDAR")).Methods("POST")
	router.HandleFunc("/usuario", handlers.ValidarMiddleware(handlers.UsuarioCrear, "NO_VALIDAR")).Methods("POST")

	// Filosofos
	// *********
	router.HandleFunc("/filosofo", handlers.ValidarMiddleware(handlers.FilosofoCrear, "NO_VALIDAR")).Methods("POST")
	router.HandleFunc("/filosofo/{filosofoID}", handlers.ValidarMiddleware(handlers.FilosofoGuardar, "NO_VALIDAR")).Methods("PUT")
	router.HandleFunc("/filosofoHabilitar/{filosofoID}", handlers.ValidarMiddleware(handlers.FilosofoHabilitar, "NO_VALIDAR")).Methods("PUT")
	router.HandleFunc("/filosofoDeshabilitar/{filosofoID}", handlers.ValidarMiddleware(handlers.FilosofoDeshabilitar, "NO_VALIDAR")).Methods("PUT")
	router.HandleFunc("/filosofo/{filosofoID}", handlers.ValidarMiddleware(handlers.FilosofoBorrar, "NO_VALIDAR")).Methods("DELETE")
	router.HandleFunc("/filosofoRecuperar/{filosofoID}", handlers.ValidarMiddleware(handlers.FilosofoRecuperar, "NO_VALIDAR")).Methods("GET")
	router.HandleFunc("/filosofo/{filosofoID}", handlers.ValidarMiddleware(handlers.FilosofoTraer, "NO_VALIDAR")).Methods("GET")
	router.HandleFunc("/filosofos/{orden}/{limite}", handlers.ValidarMiddleware(handlers.FilosofosTraer, "NO_VALIDAR")).Methods("POST")

	allowedOrigins := gorillaHnd.AllowedOrigins([]string{"*"})
	allowedMethods := gorillaHnd.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT"})
	allowedHeaders := gorillaHnd.AllowedHeaders([]string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"})
	log.Fatal(http.ListenAndServe(":3113", gorillaHnd.CORS(allowedOrigins, allowedMethods, allowedHeaders)(router)))
	//log.Fatal(http.ListenAndServe(":3113", router))
}
