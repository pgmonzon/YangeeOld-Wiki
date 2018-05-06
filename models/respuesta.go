package models

import (
)

type Respuesta struct {
  Estado      string    `json:"estado"`
  Valor       string    `json:"valor"`
  Mensaje     string    `json:"mensaje"`
}
// *****SACAR ******
type Resp struct {
  EstadoGral  string    `json:"estadoGral"`
  Mensajes    []Mensaje `json:"mensajes"`
}

type Mensaje struct {
  Valor     string    `json:"valor"`
  Estado    string    `json:"estado"`
  Mensaje   string    `json:"mensaje"`
}
