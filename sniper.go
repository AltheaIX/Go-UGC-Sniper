package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"sync"
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
	client := &http.Client{Timeout: 15 * time.Second}

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

	marketplaceDetail := UnmarshalMarketplaceDetail(scanner)
	return marketplaceDetail, err
}

func Sniper(detail *MarketplaceDetail) error {
	client := &http.Client{Timeout: 15 * time.Second}

	jsonPayload := fmt.Sprintf(`{
	"collectibleItemId": "%v",
	"expectedCurrency": 1,
	"expectedPrice": %d,
	"expectedPurchaserId": "%d",
	"expectedPurchaserType": "User",
	"expectedSellerId": %d,
	"expectedSellerType": "User",
	"idempotencyKey": "%v",
	"collectibleProductId": "%v"
}`, detail.Data[0].ItemId, detail.Data[0].Price, accountId, detail.Data[0].CreatorId, uuid.New(), detail.Data[0].ProductId)
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
		err = errors.New(fmt.Sprintf("error code is %d", response.StatusCode))
		return err
	}

	scanner, _ := ResponseReader(response)
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
	fmt.Println(string(scanner))
	return nil
}

func SniperHandler() {
	workerCount := 3 // Set the number of concurrent workers
	workerSem := make(chan struct{}, workerCount)
	var wg sync.WaitGroup

	for _, data := range listFreeItem {
		now := time.Now()

		workerSem <- struct{}{}
		wg.Add(1)

		go func(data string) {
			defer func() {
				<-workerSem
				wg.Done()
			}()

			detail, err := MarketplaceDetailByCollectibleItemId(data)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Printf("Sniper - Sniping items %v\n", detail.Data[0].Name)
			for {
				defer func() {
					elapsed := time.Since(now)
					elapsedMilliseconds := int64(elapsed / time.Millisecond)

					fmt.Printf("Sniper - Sniped within %d milliseconds\n", elapsedMilliseconds)
				}()

				err = Sniper(detail)
				if err != nil && err.Error() == "sold out" {
					fmt.Printf("Sniper - %v already sold out.\n", detail.Data[0].Name)
					listFreeItem = DeleteSlice(listFreeItem, data)
					break
				}

				if err != nil {
					listFreeItem = DeleteSlice(listFreeItem, data)
					break
				}
			}
		}(data)
	}

	wg.Wait()

	time.Sleep(1 * time.Second)
	resumeGoroutines()
}
