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

	// Cuentas Gastos
	// **************
	router.HandleFunc("/cuentaGasto", handlers.ValidarMiddleware(handlers.CuentaGastoCrear, "NO_VALIDAR")).Methods("POST")
	router.HandleFunc("/cuentaGasto/{docID}", handlers.ValidarMiddleware(handlers.CuentaGastoGuardar, "NO_VALIDAR")).Methods("PUT")
	router.HandleFunc("/cuentaGastoHabilitar/{docID}", handlers.ValidarMiddleware(handlers.CuentaGastoHabilitar, "NO_VALIDAR")).Methods("PUT")
	router.HandleFunc("/cuentaGastoDeshabilitar/{docID}", handlers.ValidarMiddleware(handlers.CuentaGastoDeshabilitar, "NO_VALIDAR")).Methods("PUT")
	router.HandleFunc("/cuentaGasto/{docID}", handlers.ValidarMiddleware(handlers.CuentaGastoBorrar, "NO_VALIDAR")).Methods("DELETE")
	router.HandleFunc("/cuentaGastoRecuperar/{docID}", handlers.ValidarMiddleware(handlers.CuentaGastoRecuperar, "NO_VALIDAR")).Methods("GET")
	router.HandleFunc("/cuentaGasto/{docID}", handlers.ValidarMiddleware(handlers.CuentaGastoTraer, "NO_VALIDAR")).Methods("GET")
	router.HandleFunc("/cuentaGastos/{orden}/{limite}", handlers.ValidarMiddleware(handlers.CuentaGastosTraer, "NO_VALIDAR")).Methods("POST")

	// Basico Sindicatos
	// *****************
	router.HandleFunc("/basicoSindicato", handlers.ValidarMiddleware(handlers.BasicoSindicatoCrear, "NO_VALIDAR")).Methods("POST")
	router.HandleFunc("/basicoSindicato/{docID}", handlers.ValidarMiddleware(handlers.BasicoSindicatoGuardar, "NO_VALIDAR")).Methods("PUT")
	router.HandleFunc("/basicoSindicatoHabilitar/{docID}", handlers.ValidarMiddleware(handlers.BasicoSindicatoHabilitar, "NO_VALIDAR")).Methods("PUT")
	router.HandleFunc("/basicoSindicatoDeshabilitar/{docID}", handlers.ValidarMiddleware(handlers.BasicoSindicatoDeshabilitar, "NO_VALIDAR")).Methods("PUT")
	router.HandleFunc("/basicoSindicato/{docID}", handlers.ValidarMiddleware(handlers.BasicoSindicatoBorrar, "NO_VALIDAR")).Methods("DELETE")
	router.HandleFunc("/basicoSindicatoRecuperar/{docID}", handlers.ValidarMiddleware(handlers.BasicoSindicatoRecuperar, "NO_VALIDAR")).Methods("GET")
	router.HandleFunc("/basicoSindicato/{docID}", handlers.ValidarMiddleware(handlers.BasicoSindicatoTraer, "NO_VALIDAR")).Methods("GET")
	router.HandleFunc("/basicoSindicatos/{orden}/{limite}", handlers.ValidarMiddleware(handlers.BasicoSindicatosTraer, "NO_VALIDAR")).Methods("POST")

	// Unidades
	// ********
	router.HandleFunc("/unidad", handlers.ValidarMiddleware(handlers.UnidadCrear, "NO_VALIDAR")).Methods("POST")
	router.HandleFunc("/unidad/{docID}", handlers.ValidarMiddleware(handlers.UnidadGuardar, "NO_VALIDAR")).Methods("PUT")
	router.HandleFunc("/unidadHabilitar/{docID}", handlers.ValidarMiddleware(handlers.UnidadHabilitar, "NO_VALIDAR")).Methods("PUT")
	router.HandleFunc("/unidadDeshabilitar/{docID}", handlers.ValidarMiddleware(handlers.UnidadDeshabilitar, "NO_VALIDAR")).Methods("PUT")
	router.HandleFunc("/unidad/{docID}", handlers.ValidarMiddleware(handlers.UnidadBorrar, "NO_VALIDAR")).Methods("DELETE")
	router.HandleFunc("/unidadRecuperar/{docID}", handlers.ValidarMiddleware(handlers.UnidadRecuperar, "NO_VALIDAR")).Methods("GET")
	router.HandleFunc("/unidad/{docID}", handlers.ValidarMiddleware(handlers.UnidadTraer, "NO_VALIDAR")).Methods("GET")
	router.HandleFunc("/unidades/{orden}/{limite}", handlers.ValidarMiddleware(handlers.UnidadesTraer, "NO_VALIDAR")).Methods("POST")

	// Personal
	// **********
	router.HandleFunc("/personal", handlers.ValidarMiddleware(handlers.PersonalCrear, "NO_VALIDAR")).Methods("POST")
	router.HandleFunc("/personal/{docID}", handlers.ValidarMiddleware(handlers.PersonalGuardar, "NO_VALIDAR")).Methods("PUT")
	router.HandleFunc("/personalHabilitar/{docID}", handlers.ValidarMiddleware(handlers.PersonalHabilitar, "NO_VALIDAR")).Methods("PUT")
	router.HandleFunc("/personalDeshabilitar/{docID}", handlers.ValidarMiddleware(handlers.PersonalDeshabilitar, "NO_VALIDAR")).Methods("PUT")
	router.HandleFunc("/personal/{docID}", handlers.ValidarMiddleware(handlers.PersonalBorrar, "NO_VALIDAR")).Methods("DELETE")
	router.HandleFunc("/personalRecuperar/{docID}", handlers.ValidarMiddleware(handlers.PersonalRecuperar, "NO_VALIDAR")).Methods("GET")
	router.HandleFunc("/personal/{docID}", handlers.ValidarMiddleware(handlers.PersonalTraer, "NO_VALIDAR")).Methods("GET")
	router.HandleFunc("/personales/{orden}/{limite}", handlers.ValidarMiddleware(handlers.PersonalesTraer, "NO_VALIDAR")).Methods("POST")

	// Locaciones
	// **********
	router.HandleFunc("/locacion", handlers.ValidarMiddleware(handlers.LocacionCrear, "NO_VALIDAR")).Methods("POST")
	router.HandleFunc("/locacion/{docID}", handlers.ValidarMiddleware(handlers.LocacionGuardar, "NO_VALIDAR")).Methods("PUT")
	router.HandleFunc("/locacionHabilitar/{docID}", handlers.ValidarMiddleware(handlers.LocacionHabilitar, "NO_VALIDAR")).Methods("PUT")
	router.HandleFunc("/locacionDeshabilitar/{docID}", handlers.ValidarMiddleware(handlers.LocacionDeshabilitar, "NO_VALIDAR")).Methods("PUT")
	router.HandleFunc("/locacion/{docID}", handlers.ValidarMiddleware(handlers.LocacionBorrar, "NO_VALIDAR")).Methods("DELETE")
	router.HandleFunc("/locacionRecuperar/{docID}", handlers.ValidarMiddleware(handlers.LocacionRecuperar, "NO_VALIDAR")).Methods("GET")
	router.HandleFunc("/locacion/{docID}", handlers.ValidarMiddleware(handlers.LocacionTraer, "NO_VALIDAR")).Methods("GET")
	router.HandleFunc("/locaciones/{orden}/{limite}", handlers.ValidarMiddleware(handlers.LocacionesTraer, "NO_VALIDAR")).Methods("POST")

	allowedOrigins := gorillaHnd.AllowedOrigins([]string{"*"})
	allowedMethods := gorillaHnd.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT"})
	allowedHeaders := gorillaHnd.AllowedHeaders([]string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"})
	log.Fatal(http.ListenAndServe(":3113", gorillaHnd.CORS(allowedOrigins, allowedMethods, allowedHeaders)(router)))
	//log.Fatal(http.ListenAndServe(":3113", router))
}
