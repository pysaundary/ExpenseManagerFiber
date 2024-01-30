package models

type Currency struct {
	Id     uint   `json :"id"`
	Name   string `json :"name"`
	Symbol string `json:"sysmbol"`
}
