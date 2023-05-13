package main

import jsoniter "github.com/json-iterator/go"

type Data struct {
	Id             int                 `json:"id"`
	Type           string              `json:"itemType"`
	Name           jsoniter.RawMessage `json:"name,omitempty"`
	Price          int                 `json:"price,omitempty"`
	Quantity       int                 `json:"totalQuantity,omitempty"`
	Image          string              `json:"imageUrl,omitempty"`
	UnitsAvailable int                 `json:"unitsAvailableForConsumption,omitempty"`
}

type ItemDetail struct {
	Detail []Data `json:"data"`
}
