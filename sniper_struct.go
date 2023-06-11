package main

import jsoniter "github.com/json-iterator/go"

type MarketplaceDetail struct {
	Data []MarketplaceData
}

type MarketplaceData struct {
	ItemId      string              `json:"collectibleItemId"`
	ProductId   string              `json:"collectibleProductId"`
	CreatorId   int                 `json:"creatorId"`
	CreatorType string              `json:"creatorType"`
	Price       int                 `json:"price"`
	Name        jsoniter.RawMessage `json:"name"`
}
