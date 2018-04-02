package main

import (
  "fmt"

  "github.com/pgmonzon/Yangee/routers"
  "github.com/pgmonzon/Yangee/config"

  "google.golang.org/appengine"
)

func main() {
  appengine.Main()
  
  fmt.Println(" ***********************")
  fmt.Println(" ¡¡¡ Wooow is Yangee !!!")
  fmt.Println(" ***********************")

  config.Inicializar()
  routers.InicializarRutas()
}
