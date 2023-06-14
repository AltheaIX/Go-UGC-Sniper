package main

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"strings"
	"time"
)

var accountCookie string
var accountId int

var listFreeItem []string
var listSnipedItem []string

var sniperSemaphore = make(chan struct{}, 1)

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
	client := &http.Client{Timeout: 600 * time.Millisecond}

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

	now := time.Now()
	defer func() {
		elapsed := time.Since(now)
		fmt.Println("time taken to take detail: ", elapsed)
	}()

	response, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		err = errors.New("status code is not 200")
		return nil, err
	}

	scanner, _ := ResponseReader(response)
	fmt.Println(string(scanner))

	marketplaceDetail := UnmarshalMarketplaceDetail(scanner)
	return marketplaceDetail, err
}

func Sniper(detail *MarketplaceDetail) error {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client := &http.Client{Transport: transport, Timeout: 2000 * time.Millisecond}

	jsonPayload := fmt.Sprintf(`{
	"collectibleItemId": "%v",
	"expectedCurrency": 1,
	"expectedPrice": %d,
	"expectedPurchaserId": "%d",
	"expectedPurchaserType": "User",
	"expectedSellerId": %d,
	"expectedSellerType": "%v",
	"idempotencyKey": "%v",
	"collectibleProductId": "%v"
}`, detail.Data[0].ItemId, detail.Data[0].Price, accountId, detail.Data[0].CreatorId, detail.Data[0].CreatorType, uuid.New(), detail.Data[0].ProductId)
	dataRequest := bytes.NewBuffer([]byte(jsonPayload))

	fmt.Println(jsonPayload)

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
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("x-csrf-token", GetCsrfToken())
	req.AddCookie(cookie)

	now := time.Now()
	defer func() {
		elapsed := time.Since(now)
		fmt.Println("time taken to snipe: ", elapsed)
	}()

	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	scanner, _ := ResponseReader(response)
	fmt.Println(string(scanner))
	if response.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("error code is %d", response.StatusCode))
		return err
	}

	if strings.Contains(string(scanner), "QuantityExhausted") {
		err = errors.New("sold out")
		return err
	}

	if strings.Contains(string(scanner), "purchase requests exceeds limit") {
		err = errors.New("limit")
		return err
	}

	if strings.Contains(string(scanner), `"purchased":false`) {
		return err
	}

	go BoughtNotifier(detail.Data[0].Name)
	time.Sleep(1 * time.Second)
	return nil
}

func SniperHandler() {
	defer handlePanic()

	for _, data := range listFreeItem {
		var detail *MarketplaceDetail
		var err error

		for _, dataSniped := range listSnipedItem {
			if dataSniped == data {
				listFreeItem = DeleteSlice(listFreeItem, data)
				resumeGoroutines()
				return
			}
		}

		for {
			detail, err = MarketplaceDetailByCollectibleItemId(data)
			if err != nil {
				fmt.Println(err)
				continue
			}

			if detail.Data[0].CreatorType == "" || detail.Data[0].CreatorId == 0 || detail.Data[0].ProductId == "" {
				fmt.Println("Sniper - Failed to take detail, retrying...")
				continue
			}

			break
		}

		for {
			sniperSemaphore <- struct{}{}

			_name := strings.Replace(string(detail.Data[0].Name), `"`, "", 2)
			detail.Data[0].Name = jsoniter.RawMessage(_name)

			fmt.Printf("Sniper - Sniping items %s\n", detail.Data[0].Name)
			err = Sniper(detail)
			if err != nil && err.Error() == "sold out" {
				fmt.Printf("Sniper - %s already sold out.\n", detail.Data[0].Name)
				listFreeItem = DeleteSlice(listFreeItem, data)
			}

			if err != nil && err.Error() != "sold out" {
				listFreeItem = DeleteSlice(listFreeItem, data)
			}

			listSnipedItem = append(listSnipedItem, data)
			<-sniperSemaphore
			break
		}
	}
	resumeGoroutines()
}
