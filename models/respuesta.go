package models

import (
)

type Respuesta struct {
  EstadoGral  string    `json:"estadoGral"`
  Mensaje     string    `json:"mensaje"`
}

type Resp struct {
  EstadoGral  string    `json:"estadoGral"`
  Mensajes    []Mensaje `json:"mensajes"`
}

type Mensaje struct {
  Valor     string    `json:"valor"`
  Estado    string    `json:"estado"`
  Mensaje   string    `json:"mensaje"`
}
