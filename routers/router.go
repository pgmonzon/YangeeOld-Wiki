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
	router.HandleFunc("/usuarioValidar", handlers.ValidarMiddleware(handlers.UsuarioValidar, "NO_VALIDAR")).Methods("POST")

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

	// Clientes
	// ********
	router.HandleFunc("/cliente", handlers.ValidarMiddleware(handlers.ClienteCrear, "NO_VALIDAR")).Methods("POST")
	router.HandleFunc("/cliente/{docID}", handlers.ValidarMiddleware(handlers.ClienteGuardar, "NO_VALIDAR")).Methods("PUT")
	router.HandleFunc("/clienteHabilitar/{docID}", handlers.ValidarMiddleware(handlers.ClienteHabilitar, "NO_VALIDAR")).Methods("PUT")
	router.HandleFunc("/clienteDeshabilitar/{docID}", handlers.ValidarMiddleware(handlers.ClienteDeshabilitar, "NO_VALIDAR")).Methods("PUT")
	router.HandleFunc("/cliente/{docID}", handlers.ValidarMiddleware(handlers.ClienteBorrar, "NO_VALIDAR")).Methods("DELETE")
	router.HandleFunc("/clienteRecuperar/{docID}", handlers.ValidarMiddleware(handlers.ClienteRecuperar, "NO_VALIDAR")).Methods("GET")
	router.HandleFunc("/cliente/{docID}", handlers.ValidarMiddleware(handlers.ClienteTraer, "NO_VALIDAR")).Methods("GET")
	router.HandleFunc("/clientes/{orden}/{limite}", handlers.ValidarMiddleware(handlers.ClientesTraer, "NO_VALIDAR")).Methods("POST")

	// Transportistas
	// **************
	router.HandleFunc("/transportista", handlers.ValidarMiddleware(handlers.TransportistaCrear, "NO_VALIDAR")).Methods("POST")
	router.HandleFunc("/transportista/{docID}", handlers.ValidarMiddleware(handlers.TransportistaGuardar, "NO_VALIDAR")).Methods("PUT")
	router.HandleFunc("/transportistaHabilitar/{docID}", handlers.ValidarMiddleware(handlers.TransportistaHabilitar, "NO_VALIDAR")).Methods("PUT")
	router.HandleFunc("/transportistaDeshabilitar/{docID}", handlers.ValidarMiddleware(handlers.TransportistaDeshabilitar, "NO_VALIDAR")).Methods("PUT")
	router.HandleFunc("/transportista/{docID}", handlers.ValidarMiddleware(handlers.TransportistaBorrar, "NO_VALIDAR")).Methods("DELETE")
	router.HandleFunc("/transportistaRecuperar/{docID}", handlers.ValidarMiddleware(handlers.TransportistaRecuperar, "NO_VALIDAR")).Methods("GET")
	router.HandleFunc("/transportista/{docID}", handlers.ValidarMiddleware(handlers.TransportistaTraer, "NO_VALIDAR")).Methods("GET")
	router.HandleFunc("/transportistas/{orden}/{limite}", handlers.ValidarMiddleware(handlers.TransportistasTraer, "NO_VALIDAR")).Methods("POST")

	// Viajes
	// ******
	router.HandleFunc("/viaje", handlers.ValidarMiddleware(handlers.ViajeCrear, "NO_VALIDAR")).Methods("POST")
	router.HandleFunc("/viaje/{docID}", handlers.ValidarMiddleware(handlers.ViajeGuardar, "NO_VALIDAR")).Methods("PUT")
	router.HandleFunc("/viajes", handlers.ValidarMiddleware(handlers.ViajesTraer, "NO_VALIDAR")).Methods("POST")
	router.HandleFunc("/viajeCancelar/{docID}", handlers.ValidarMiddleware(handlers.ViajeCancelar, "NO_VALIDAR")).Methods("POST")
	router.HandleFunc("/viajeRemitos/{docID}", handlers.ValidarMiddleware(handlers.ViajeRemitos, "NO_VALIDAR")).Methods("PUT")
	router.HandleFunc("/viajesFacturar/{docID}", handlers.ValidarMiddleware(handlers.ViajesFacturar, "NO_VALIDAR")).Methods("POST")
	router.HandleFunc("/viajesLiquidar/{docID}", handlers.ValidarMiddleware(handlers.ViajesLiquidar, "NO_VALIDAR")).Methods("POST")

	// Autorizaciones
	// **************
	router.HandleFunc("/viajeAutValor/{docID}", handlers.ValidarMiddleware(handlers.ViajeAutValor, "NO_VALIDAR")).Methods("POST")
	router.HandleFunc("/viajeAutCosto/{docID}", handlers.ValidarMiddleware(handlers.ViajeAutCosto, "NO_VALIDAR")).Methods("POST")
	router.HandleFunc("/autorizaciones", handlers.ValidarMiddleware(handlers.AutorizacionesTraer, "NO_VALIDAR")).Methods("POST")
	router.HandleFunc("/autorizacion/{docID}", handlers.ValidarMiddleware(handlers.AutorizacionTraer, "NO_VALIDAR")).Methods("GET")
	router.HandleFunc("/autorizacionRechazar/{docID}", handlers.ValidarMiddleware(handlers.AutorizacionRechazar, "NO_VALIDAR")).Methods("PUT")
	router.HandleFunc("/autorizacionAutorizar/{docID}", handlers.ValidarMiddleware(handlers.AutorizacionAutorizar, "NO_VALIDAR")).Methods("PUT")

	// Facturas
	// ********
	router.HandleFunc("/factura", handlers.ValidarMiddleware(handlers.FacturaCrear, "NO_VALIDAR")).Methods("POST")
	router.HandleFunc("/factura/{docID}", handlers.ValidarMiddleware(handlers.FacturaTraer, "NO_VALIDAR")).Methods("GET")
	router.HandleFunc("/facturas", handlers.ValidarMiddleware(handlers.FacturasTraer, "NO_VALIDAR")).Methods("POST")
	router.HandleFunc("/facturaViajes/{docID}", handlers.ValidarMiddleware(handlers.FacturaViajesTraer, "NO_VALIDAR")).Methods("GET")
	router.HandleFunc("/facturasTraerFechas", handlers.ValidarMiddleware(handlers.FacturasTraerFechas, "NO_VALIDAR")).Methods("POST")

	// Liquidaciones
	// *************
	router.HandleFunc("/liquidacion", handlers.ValidarMiddleware(handlers.LiquidacionCrear, "NO_VALIDAR")).Methods("POST")
	router.HandleFunc("/liquidacion/{docID}", handlers.ValidarMiddleware(handlers.LiquidacionTraer, "NO_VALIDAR")).Methods("GET")
	router.HandleFunc("/liquidaciones", handlers.ValidarMiddleware(handlers.LiquidacionesTraer, "NO_VALIDAR")).Methods("POST")
	router.HandleFunc("/liquidacionViajes/{docID}", handlers.ValidarMiddleware(handlers.LiquidacionViajesTraer, "NO_VALIDAR")).Methods("GET")
	router.HandleFunc("/liquidacionFactura/{docID}", handlers.ValidarMiddleware(handlers.LiquidacionFacturaTransportista, "NO_VALIDAR")).Methods("PUT")

	// Rendiciones
	// ***********
	router.HandleFunc("/rendicion", handlers.ValidarMiddleware(handlers.RendicionCrear, "NO_VALIDAR")).Methods("POST")
	router.HandleFunc("/rendicionPersonal/{docID}", handlers.ValidarMiddleware(handlers.RendicionPersonalTraer, "NO_VALIDAR")).Methods("GET")

	// Haberes
	// ******
	router.HandleFunc("/haberes", handlers.ValidarMiddleware(handlers.HaberesCrear, "NO_VALIDAR")).Methods("POST")

	// SbrSucursal
	// **********
	router.HandleFunc("/sbrSucursal", handlers.ValidarMiddleware(handlers.SbrSucursalCrear, "NO_VALIDAR")).Methods("POST")
	router.HandleFunc("/sbrSucursal/{docID}", handlers.ValidarMiddleware(handlers.SbrSucursalGuardar, "NO_VALIDAR")).Methods("PUT")
	router.HandleFunc("/sbrSucursalHabilitar/{docID}", handlers.ValidarMiddleware(handlers.SbrSucursalHabilitar, "NO_VALIDAR")).Methods("PUT")
	router.HandleFunc("/sbrSucursalDeshabilitar/{docID}", handlers.ValidarMiddleware(handlers.SbrSucursalDeshabilitar, "NO_VALIDAR")).Methods("PUT")
	router.HandleFunc("/sbrSucursal/{docID}", handlers.ValidarMiddleware(handlers.SbrSucursalBorrar, "NO_VALIDAR")).Methods("DELETE")
	router.HandleFunc("/sbrSucursalRecuperar/{docID}", handlers.ValidarMiddleware(handlers.SbrSucursalRecuperar, "NO_VALIDAR")).Methods("GET")
	router.HandleFunc("/sbrSucursal/{docID}", handlers.ValidarMiddleware(handlers.SbrSucursalTraer, "NO_VALIDAR")).Methods("GET")
	router.HandleFunc("/sbrSucursales/{orden}/{limite}", handlers.ValidarMiddleware(handlers.SbrSucursalesTraer, "NO_VALIDAR")).Methods("POST")

	// SbrRubro
	// **********
	router.HandleFunc("/sbrRubro", handlers.ValidarMiddleware(handlers.SbrRubroCrear, "NO_VALIDAR")).Methods("POST")
	router.HandleFunc("/sbrRubro/{docID}", handlers.ValidarMiddleware(handlers.SbrRubroGuardar, "NO_VALIDAR")).Methods("PUT")
	router.HandleFunc("/sbrRubroHabilitar/{docID}", handlers.ValidarMiddleware(handlers.SbrRubroHabilitar, "NO_VALIDAR")).Methods("PUT")
	router.HandleFunc("/sbrRubroDeshabilitar/{docID}", handlers.ValidarMiddleware(handlers.SbrRubroDeshabilitar, "NO_VALIDAR")).Methods("PUT")
	router.HandleFunc("/sbrRubro/{docID}", handlers.ValidarMiddleware(handlers.SbrRubroBorrar, "NO_VALIDAR")).Methods("DELETE")
	router.HandleFunc("/sbrRubroRecuperar/{docID}", handlers.ValidarMiddleware(handlers.SbrRubroRecuperar, "NO_VALIDAR")).Methods("GET")
	router.HandleFunc("/sbrRubro/{docID}", handlers.ValidarMiddleware(handlers.SbrRubroTraer, "NO_VALIDAR")).Methods("GET")
	router.HandleFunc("/sbrRubros/{orden}/{limite}", handlers.ValidarMiddleware(handlers.SbrRubrosTraer, "NO_VALIDAR")).Methods("POST")

	// SbrArticulo
	// ***********
	router.HandleFunc("/sbrArticulo", handlers.ValidarMiddleware(handlers.SbrArticuloCrear, "NO_VALIDAR")).Methods("POST")
	router.HandleFunc("/sbrArticulo/{docID}", handlers.ValidarMiddleware(handlers.SbrArticuloGuardar, "NO_VALIDAR")).Methods("PUT")
	router.HandleFunc("/sbrArticuloHabilitar/{docID}", handlers.ValidarMiddleware(handlers.SbrArticuloHabilitar, "NO_VALIDAR")).Methods("PUT")
	router.HandleFunc("/sbrArticuloDeshabilitar/{docID}", handlers.ValidarMiddleware(handlers.SbrArticuloDeshabilitar, "NO_VALIDAR")).Methods("PUT")
	router.HandleFunc("/sbrArticulo/{docID}", handlers.ValidarMiddleware(handlers.SbrArticuloBorrar, "NO_VALIDAR")).Methods("DELETE")
	router.HandleFunc("/sbrArticuloRecuperar/{docID}", handlers.ValidarMiddleware(handlers.SbrArticuloRecuperar, "NO_VALIDAR")).Methods("GET")
	router.HandleFunc("/sbrArticulo/{docID}", handlers.ValidarMiddleware(handlers.SbrArticuloTraer, "NO_VALIDAR")).Methods("GET")
	router.HandleFunc("/sbrArticuloCodBarras/{codBarras}", handlers.ValidarMiddleware(handlers.SbrArticuloCodBarrasTraer, "NO_VALIDAR")).Methods("GET")
	router.HandleFunc("/sbrArticulos/{orden}/{limite}", handlers.ValidarMiddleware(handlers.SbrArticulosTraer, "NO_VALIDAR")).Methods("POST")

	// SbrIngresoSucursal
	// ******************
	router.HandleFunc("/sbrIngresoSucursal", handlers.ValidarMiddleware(handlers.SbrIngresoSucursalCrear, "NO_VALIDAR")).Methods("POST")
	router.HandleFunc("/sbrIngresoSucursal/{docID}", handlers.ValidarMiddleware(handlers.SbrIngresoSucursalTraer, "NO_VALIDAR")).Methods("GET")
	router.HandleFunc("/sbrIngresosSucursales/{orden}/{limite}/{sucursal}", handlers.ValidarMiddleware(handlers.SbrIngresosSucursalesTraer, "NO_VALIDAR")).Methods("POST")

	// SbrRemitoSucursal
	// ******************
	router.HandleFunc("/sbrRemitoSucursal", handlers.ValidarMiddleware(handlers.SbrRemitoSucursalCrear, "NO_VALIDAR")).Methods("POST")
	router.HandleFunc("/sbrRemitoSucursal/{docID}", handlers.ValidarMiddleware(handlers.SbrRemitoSucursalTraer, "NO_VALIDAR")).Methods("GET")
	router.HandleFunc("/sbrRemitosDeSucursal/{orden}/{limite}/{sucursal}", handlers.ValidarMiddleware(handlers.SbrRemitosDeSucursalTraer, "NO_VALIDAR")).Methods("POST")
	router.HandleFunc("/sbrRemitosASucursal/{orden}/{limite}/{sucursal}", handlers.ValidarMiddleware(handlers.SbrRemitosASucursalTraer, "NO_VALIDAR")).Methods("POST")
	router.HandleFunc("/sbrRemitoSucursalCancelar/{docID}", handlers.ValidarMiddleware(handlers.SbrRemitoSucursalCancelar, "NO_VALIDAR")).Methods("PUT")
	router.HandleFunc("/sbrRemitoSucursalAceptar/{docID}", handlers.ValidarMiddleware(handlers.SbrRemitoSucursalAceptar, "NO_VALIDAR")).Methods("PUT")
	router.HandleFunc("/sbrRemitoSucursalRechazar/{docID}", handlers.ValidarMiddleware(handlers.SbrRemitoSucursalRechazar, "NO_VALIDAR")).Methods("PUT")

	allowedOrigins := gorillaHnd.AllowedOrigins([]string{"*"})
	allowedMethods := gorillaHnd.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT"})
	allowedHeaders := gorillaHnd.AllowedHeaders([]string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"})
	log.Fatal(http.ListenAndServe(":3113", gorillaHnd.CORS(allowedOrigins, allowedMethods, allowedHeaders)(router)))
}
