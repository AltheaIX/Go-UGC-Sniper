package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"time"
)

var accountCookie string
var accountId int

var listFreeItem []string

func UnmarshalMarketplaceDetail(responseRaw []byte) *MarketplaceDetail {
	marketplaceData := &[]MarketplaceData{}

	err := json.Unmarshal(responseRaw, &marketplaceData)
	if err != nil {
		fmt.Println(err)
	}

	marketplaceDetail := &MarketplaceDetail{Data: *marketplaceData}
	return marketplaceDetail
}

func MarketplaceDetailByCollectibleItemId(collectibleItemId string) (*MarketplaceDetail, error) {
	client := &http.Client{Timeout: 3 * time.Second}

	jsonPayload := fmt.Sprintf(`{"itemIds": ["%v"]}`, collectibleItemId)
	dataRequest := bytes.NewBuffer([]byte(jsonPayload))

	cookie := &http.Cookie{
		Name:    ".ROBLOSECURITY",
		Path:    "/",
		Value:   accountCookie,
		Domain:  "roblox.com",
		Expires: time.Now().Add(time.Hour * 1000),
	}

	req, err := http.NewRequest("POST", "https://apis.roblox.com/marketplace-items/v1/items/details", dataRequest)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "PostmanRuntime/7.29.0")
	req.Header.Set("Connection", "keep-alive")
	req.AddCookie(cookie)

	response, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		err = errors.New("status code is not 200")
	}

	scanner, _ := ResponseReader(response)
	fmt.Println(scanner)

	marketplaceDetail := UnmarshalMarketplaceDetail(scanner)
	return marketplaceDetail, err
}

func Sniper(detail *MarketplaceDetail) error {
	client := &http.Client{Timeout: 3 * time.Second}

	jsonPayload := fmt.Sprintf(`{
	"collectibleItemId": "%v",
	"expectedCurrency": 1,
	"expectedPrice": 0,
	"expectedPurchaserId": %d,
	"expectedPurchaserType": "User",
	"expectedSellerId": %d,
	"expectedSellerType": "User",
	"idempotencyKey": "%v",
	"collectibleProductId": "%v"
}`, detail.Data[0].ItemId, accountId, detail.Data[0].CreatorId, uuid.New(), detail.Data[0].ProductId)
	dataRequest := bytes.NewBuffer([]byte(jsonPayload))

	cookie := &http.Cookie{
		Name:    ".ROBLOSECURITY",
		Value:   accountCookie,
		Expires: time.Now().Add(time.Hour * 1000),
	}

	urlBuilder := fmt.Sprintf("https://apis.roblox.com/marketplace-sales/v1/item/%v/purchase-item", detail.Data[0].ItemId)
	req, err := http.NewRequest("POST", urlBuilder, dataRequest)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-csrf-token", GetCsrfToken())
	req.AddCookie(cookie)

	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		err = errors.New("status code is not 200")
		return err
	}

	scanner, _ := ResponseReader(response)
	fmt.Println(string(scanner))
	return err
}

func SniperHandler() {
	for _, data := range listFreeItem {
		fmt.Println("Sniper - Sniping items.")
		detail, err := MarketplaceDetailByCollectibleItemId(data)
		if err != nil {
			fmt.Println(err)
			time.Sleep(15 * time.Second)
		}

		err = Sniper(detail)
		fmt.Println(err)
	}
}
