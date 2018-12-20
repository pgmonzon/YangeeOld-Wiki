package models

import (
	"gopkg.in/mgo.v2/bson"
)

type Opcion struct {
	Path				string					`json:"path"`
	Type				string					`json:"type"`
	Title				string					`json:"title"`
	Icontype		string					`json:"icontype"`
	Collapse		string					`json:"collapse, omitempty"`
	Children		[]Sub						`json:"children,omitempty"`
}

type Sub struct {
	Path				string					`json:"path"`
	Title				string					`json:"title"`
	Ab					string					`json:"ab"`
}

type Menu struct {
	ID							          bson.ObjectId	 	`bson:"_id" json:"id,omitempty"`
	Empresa_id			          bson.ObjectId		`bson:"empresa_id" json:"empresa_id,omitempty"`
  Empresa                   string          `bson:"empresa" json:"empresa"`
  Menu			               	[]OpcionEstado	`json:"menu, omitempty"`
}

type OpcionEstado struct {
	Path				string					`json:"path"`
	Type				string					`json:"type"`
	Title				string					`json:"title"`
	Icontype		string					`json:"icontype"`
	Collapse		string					`json:"collapse, omitempty"`
  Estado      bool            `json:"estado"`
	Children		[]SubEstado			`json:"children,omitempty"`
}

type SubEstado struct {
	Path				string					`json:"path"`
	Title				string					`json:"title"`
	Ab					string					`json:"ab"`
  Estado      bool            `json:"estado"`
}
