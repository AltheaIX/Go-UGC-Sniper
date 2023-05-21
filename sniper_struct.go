package main

type MarketplaceDetail struct {
	Data []MarketplaceData
}

type MarketplaceData struct {
	ItemId    string `json:"collectibleItemId"`
	ProductId string `json:"collectibleProductId"`
	CreatorId int    `json:"creatorId"`
}
