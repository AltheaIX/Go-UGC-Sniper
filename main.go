package main

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/net/proxy"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

const MAX_PRICE = 50

var json = jsoniter.ConfigCompatibleWithStandardLibrary
var listId []int
var lastItemId int = 13197718725

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
	ReadProxyFromFile("proxy_fresh")
	proxyURL, err := url.Parse("socks5://" + proxyList[rand.Intn(len(proxyList))])
	if err != nil {
		fmt.Println("error on parser.")
		panic(err)
	}

	dialer, err := proxy.FromURL(proxyURL, proxy.Direct)
	if err != nil {
		fmt.Println("error on dialer.")
		panic(err)
	}

	transport := &http.Transport{
		Dial:            dialer.Dial,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Skip certificate verification
	}

	client := &http.Client{Transport: transport, Timeout: 5 * time.Second}
	response, err := client.Get("https://catalog.roblox.com/v1/search/items?category=All&includeNotForSale=true&limit=120&salesTypeFilter=2&sortType=3")
	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	defer response.Body.Close()
	if response.StatusCode != 200 {
		return nil, errors.New("status code is not 200")
	}

	scanner, _ := io.ReadAll(response.Body)
	time.Sleep(2 * time.Second)
	return scanner, nil
}

func ItemRecentlyAddedAppend(scanner []byte, err error) (int, error) {
	if err != nil {
		return 0, err
	}

	listId = nil

	listItems := UnmarshalCatalog(scanner)
	for _, data := range listItems.Detail {
		listId = append(listId, data.Id)
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

func main() {
	AddToWatcher()
}
