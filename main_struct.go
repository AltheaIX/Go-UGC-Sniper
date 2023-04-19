package main

type Data struct {
	Id       int    `json:"id"`
	Type     string `json:"itemType"`
	Name     string `json:"name,omitempty"`
	Price    int    `json:"price,omitempty"`
	Quantity int    `json:"totalQuantity,omitempty"`
	Image    string `json:"imageUrl,omitempty"`
}

type ItemDetail struct {
	Detail []Data `json:"data"`
}
