package main

import (
	"bytes"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io"
	"net/http"
	"time"
)

const MAX_PRICE = 50

var json = jsoniter.ConfigCompatibleWithStandardLibrary
var listId []int
var wathcerId []int

func UnmarshalCatalog(responseRaw []byte) *ItemDetail {
	itemDetail := &ItemDetail{}

	err := json.Unmarshal(responseRaw, &itemDetail)
	if err != nil {
		fmt.Println(err)
	}
	return itemDetail
}

func GetCsrfToken() string {
	client := &http.Client{Timeout: 5 * time.Second}

	jsonRequest := fmt.Sprintf(`{"items":[{"itemType": 1, "id": %x}]}`, 13177094956)
	dataRequest := bytes.NewBuffer([]byte(jsonRequest))

	req, err := http.NewRequest("POST", "https://catalog.roblox.com/v1/catalog/items/details", dataRequest)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()

	var token string
	if response.Header.Get("x-csrf-token") != "" {
		token = response.Header.Get("x-csrf-token")
	}

	time.Sleep(0 * time.Second)
	return token
}

func ItemRecentlyAdded() ([]byte, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	response, err := client.Get("https://catalog.roblox.com/v1/search/items?category=All&includeNotForSale=true&limit=10&salesTypeFilter=2&sortType=3")
	if err != nil {
		fmt.Print(err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return nil, errors.New("status code is not 200")
	}

	scanner, _ := io.ReadAll(response.Body)
	time.Sleep(2 * time.Second)
	return scanner, nil
}

func ItemRecentlyAddedAppend() (int, error) {
	responseListItems, err := ItemRecentlyAdded()
	if err != nil {
		return 0, err
	}

	listItems := UnmarshalCatalog(responseListItems)
	for _, data := range listItems.Detail {
		listId = append(listId, data.Id)
		fmt.Println(data.Id)
	}

	return listItems.Detail[0].Id, nil
}

func ItemDetailById(assetId int) (*ItemDetail, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	itemDetail := &ItemDetail{}

	jsonPayload := fmt.Sprintf(`{"items":[{"itemType": 1, "id": %d}]}`, assetId)
	dataRequest := bytes.NewBuffer([]byte(jsonPayload))

	req, err := http.NewRequest("POST", "https://catalog.roblox.com/v1/catalog/items/details", dataRequest)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-csrf-token", GetCsrfToken())

	response, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return nil, errors.New("status code is not 200")
	}

	scanner, _ := io.ReadAll(response.Body)
	itemDetail = UnmarshalCatalog(scanner)
	time.Sleep(2 * time.Second)
	return itemDetail, nil
}

func ItemThumbnailImageById(assetId int) (string, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	urlBuilder := fmt.Sprintf("https://thumbnails.roblox.com/v1/assets?assetIds=%d&returnPolicy=PlaceHolder&size=420x420&format=Png&isCircular=false", assetId)
	data := &ItemDetail{}

	response, err := client.Get(urlBuilder)
	if err != nil {
		fmt.Print(err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return "", errors.New("status code is not 200")
	}

	scanner, _ := io.ReadAll(response.Body)
	err = json.Unmarshal(scanner, &data)
	if err != nil {
		return "", err
	}

	return data.Detail[0].Image, nil
}

func AddToWatcher() {
	for {
		lastId, err := ItemRecentlyAddedAppend()
		if err != nil {
			fmt.Println("Trying to recovery too many request...")
			time.Sleep(15 * time.Second)
			fmt.Println("Sleep done, lets see if it works...")
			continue
		}

		if listId[0] == lastId {
			continue
		}

		fmt.Println(lastId)
		time.Sleep(3 * time.Second)
	}
}

func NotifierWatcher(notifierType string, data Data) error {
	switch notifierType {
	case "free":
		client := &http.Client{Timeout: 5 * time.Second}

		webhookBuilder := fmt.Sprintf(`{
		  "content": null,
		  "embeds": [
			{
			  "title": "%s",
			  "url": "https://www.roblox.com/catalog/%d/",
			  "color": 4628704,
			  "fields": [
				{
				  "name": "Price",
				  "value": "%d",
				  "inline": true
				},
				{
				  "name": "Quantity",
				  "value": "%d",
				  "inline": true
				},
				{
				  "name": "Item Id",
				  "value": "%d"
				}
			  ],
			  "thumbnail": {
				"url": "%s"
			  }
			}
		  ],
		  "username": "Test",
		  "attachments": []
		}`, data.Name, data.Id, data.Price, data.Quantity, data.Id, data.Image)
		dataRequest := bytes.NewBuffer([]byte(webhookBuilder))

		req, err := http.NewRequest("POST", "https://discord.com/api/webhooks/1098275385801719919/r24Vz-TimcV2baMCMKbeuDsJGmR-wWCZYOb_vlkTts5jCFPiT1jU4DDqDyhSiATHfYWw", dataRequest)
		if err != nil {
			fmt.Println(err)
		}
		req.Header.Set("Content-Type", "application/json")

		response, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		defer response.Body.Close()
		if response.StatusCode != 200 {
			return errors.New("string code is not 200")
		}
		break
	case "paid":
		break
	}
	return nil
}

func main() {

}
