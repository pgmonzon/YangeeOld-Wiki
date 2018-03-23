package models

import (
)

type Respuesta struct {
  Estado      string    `json:"estado"`
  Valor       string    `json:"valor"`
  Mensaje     string    `json:"mensaje"`
}
