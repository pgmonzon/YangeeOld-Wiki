package main

import (
  "fmt"

  "github.com/pgmonzon/Yangee/routers"
  "github.com/pgmonzon/Yangee/config"
)

func main() {
  fmt.Println(" **************************")
  fmt.Println(" ¡¡¡ Wooow Yangee is up !!!")
  fmt.Println(" **************************")

  config.Inicializar()
  routers.InicializarRutas()
}
