package main

import jsoniter "github.com/json-iterator/go"

type Data struct {
	Id                int                 `json:"id"`
	Type              string              `json:"itemType"`
	Name              jsoniter.RawMessage `json:"name,omitempty"`
	Price             int                 `json:"price,omitempty"`
	Quantity          int                 `json:"totalQuantity,omitempty"`
	SaleLocationType  string              `json:"saleLocationType,omitempty"`
	Image             string              `json:"imageUrl,omitempty"`
	UnitsAvailable    int                 `json:"unitsAvailableForConsumption,omitempty"`
	PriceStatus       string              `json:"priceStatus,omitempty"`
	CollectibleItemId string              `json:"collectibleItemId,omitempty"`
}

type User struct {
	Id           int    `json:"UserId"`
	Username     string `json:"Username"`
	RobuxBalance int    `json:"RobuxBalance"`
	Premium      bool   `json:"IsPremium"`
}

type ItemDetail struct {
	Detail []Data `json:"data"`
}

type OffsaleItems struct {
	ItemType int `json:"itemType"`
	ID       int `json:"id"`
}

type OffsalePayload struct {
	Items []OffsaleItems `json:"items"`
}
