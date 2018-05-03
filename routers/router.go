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
	router.HandleFunc("/filosofo/{docID}", handlers.ValidarMiddleware(handlers.FilosofoGuardar, "NO_VALIDAR")).Methods("PUT")
	router.HandleFunc("/filosofoHabilitar/{docID}", handlers.ValidarMiddleware(handlers.FilosofoHabilitar, "NO_VALIDAR")).Methods("PUT")
	router.HandleFunc("/filosofoDeshabilitar/{docID}", handlers.ValidarMiddleware(handlers.FilosofoDeshabilitar, "NO_VALIDAR")).Methods("PUT")
	router.HandleFunc("/filosofo/{docID}", handlers.ValidarMiddleware(handlers.FilosofoBorrar, "NO_VALIDAR")).Methods("DELETE")
	router.HandleFunc("/filosofoRecuperar/{docID}", handlers.ValidarMiddleware(handlers.FilosofoRecuperar, "NO_VALIDAR")).Methods("GET")
	router.HandleFunc("/filosofo/{docID}", handlers.ValidarMiddleware(handlers.FilosofoTraer, "NO_VALIDAR")).Methods("GET")
	router.HandleFunc("/filosofos/{orden}/{limite}", handlers.ValidarMiddleware(handlers.FilosofosTraer, "NO_VALIDAR")).Methods("POST")

	// TiposUnidades
	// *************
	router.HandleFunc("/tipoUnidad", handlers.ValidarMiddleware(handlers.TipoUnidadCrear, "NO_VALIDAR")).Methods("POST")
	router.HandleFunc("/tipoUnidad/{docID}", handlers.ValidarMiddleware(handlers.TipoUnidadGuardar, "NO_VALIDAR")).Methods("PUT")
	router.HandleFunc("/tipoUnidadHabilitar/{docID}", handlers.ValidarMiddleware(handlers.TipoUnidadHabilitar, "NO_VALIDAR")).Methods("PUT")
	router.HandleFunc("/tipoUnidadDeshabilitar/{docID}", handlers.ValidarMiddleware(handlers.TipoUnidadDeshabilitar, "NO_VALIDAR")).Methods("PUT")
	router.HandleFunc("/tipoUnidad/{docID}", handlers.ValidarMiddleware(handlers.TipoUnidadBorrar, "NO_VALIDAR")).Methods("DELETE")
	router.HandleFunc("/tipoUnidadRecuperar/{docID}", handlers.ValidarMiddleware(handlers.TipoUnidadRecuperar, "NO_VALIDAR")).Methods("GET")
	router.HandleFunc("/tipoUnidad/{docID}", handlers.ValidarMiddleware(handlers.TipoUnidadTraer, "NO_VALIDAR")).Methods("GET")
	router.HandleFunc("/tipoUnidades/{orden}/{limite}", handlers.ValidarMiddleware(handlers.TipoUnidadesTraer, "NO_VALIDAR")).Methods("POST")

	// Categorías
	// **********
	router.HandleFunc("/categoria", handlers.ValidarMiddleware(handlers.CategoriaCrear, "NO_VALIDAR")).Methods("POST")
	router.HandleFunc("/categoria/{docID}", handlers.ValidarMiddleware(handlers.CategoriaGuardar, "NO_VALIDAR")).Methods("PUT")
	router.HandleFunc("/categoriaHabilitar/{docID}", handlers.ValidarMiddleware(handlers.CategoriaHabilitar, "NO_VALIDAR")).Methods("PUT")
	router.HandleFunc("/categoriaDeshabilitar/{docID}", handlers.ValidarMiddleware(handlers.CategoriaDeshabilitar, "NO_VALIDAR")).Methods("PUT")
	router.HandleFunc("/categoria/{docID}", handlers.ValidarMiddleware(handlers.CategoriaBorrar, "NO_VALIDAR")).Methods("DELETE")
	router.HandleFunc("/categoriaRecuperar/{docID}", handlers.ValidarMiddleware(handlers.CategoriaRecuperar, "NO_VALIDAR")).Methods("GET")
	router.HandleFunc("/categoria/{docID}", handlers.ValidarMiddleware(handlers.CategoriaTraer, "NO_VALIDAR")).Methods("GET")
	router.HandleFunc("/categorias/{orden}/{limite}", handlers.ValidarMiddleware(handlers.CategoriasTraer, "NO_VALIDAR")).Methods("POST")

	allowedOrigins := gorillaHnd.AllowedOrigins([]string{"*"})
	allowedMethods := gorillaHnd.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT"})
	allowedHeaders := gorillaHnd.AllowedHeaders([]string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"})
	log.Fatal(http.ListenAndServe(":3113", gorillaHnd.CORS(allowedOrigins, allowedMethods, allowedHeaders)(router)))
	//log.Fatal(http.ListenAndServe(":3113", router))
}
