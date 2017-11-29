package main

import (
  "fmt"

  "github.com/pgmonzon/Yangee/routers"
  "github.com/pgmonzon/Yangee/config"
)

func main() {
  fmt.Println(" ***********************")
  fmt.Println(" ¡¡¡ Wooow is Yangee !!!")
  fmt.Println(" ***********************")

  config.Inicializar()
  routers.InicializarRutas()
}
