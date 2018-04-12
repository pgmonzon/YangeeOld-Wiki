package models

import (
)

type Opcion struct {
	Opcion string		`json:"opcion"`
	Sub    []Sub		`json:"sub"`
}

type Sub struct {
  Opcion  string  `json:"opcion"`
}
