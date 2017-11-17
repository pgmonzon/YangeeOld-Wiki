package routers

import (
	"net/http"
  "log"

  "github.com/pgmonzon/Yangee/handlers"

  "github.com/gorilla/mux"
)

func InicializarRutas() {
  router := mux.NewRouter()

  // Autorización: Genera token
	router.HandleFunc("/autorizar", handlers.Autorizar).Methods("POST")

	log.Fatal(http.ListenAndServe(":3113", router))
}
