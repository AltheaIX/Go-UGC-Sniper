package main

import (
	"bufio"
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
var lastItemId int

func ResponseReader(response *http.Response) ([]byte, error) {
	var body []byte
	var err error

	reader := bufio.NewReaderSize(response.Body, 4096*2)
	for {
		chunk, err := reader.ReadSlice('\n')
		if err != nil && err != io.EOF {
			fmt.Printf("Error: %v - Response Reader\n", err)
		}
		body = append(body, chunk...)
		if err == io.EOF {
			break
		}
	}
	return body, err
}

func GetCsrfToken() string {
	client := &http.Client{}

	jsonRequest := fmt.Sprintf(`{"items":[{"itemType": 1, "id": %x}]}`, 13177094956)
	dataRequest := bytes.NewBuffer([]byte(jsonRequest))

	req, err := http.NewRequest("POST", "https://catalog.roblox.com/v1/catalog/items/details", dataRequest)
	if err != nil {
		fmt.Printf("Error: %v - GetCsrfToken\n", err)
	}
	req.Header.Set("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error: %v - GetCsrfToken\n", err)
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
		fmt.Printf("Error: %v - Parser Proxy\n", err)
		panic(err)
	}

	dialer, err := proxy.FromURL(proxyURL, proxy.Direct)
	if err != nil {
		fmt.Printf("Error: %v - Dialer Proxy\n", err)
		panic(err)
	}

	transport := &http.Transport{
		Dial:            dialer.Dial,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Skip certificate verification
	}

	client := &http.Client{Transport: transport, Timeout: 3 * time.Second}
	response, err := client.Get("https://catalog.roblox.com/v1/search/items?category=All&includeNotForSale=true&limit=120&salesTypeFilter=2&sortType=3")
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	if response.StatusCode != 200 {
		return nil, errors.New("status code is not 200")
	}

	scanner, _ := ResponseReader(response)

	if string(scanner) == "" {
		return nil, errors.New("empty body")
	}

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
	itemDetail := &ItemDetail{}
	var err error
	client := &http.Client{Timeout: 3 * time.Second}

	jsonPayload := fmt.Sprintf(`{"items":[{"itemType": 1, "id": %d}]}`, assetId)
	dataRequest := bytes.NewBuffer([]byte(jsonPayload))

	req, err := http.NewRequest("POST", "https://catalog.roblox.com/v1/catalog/items/details", dataRequest)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-csrf-token", GetCsrfToken())

	for {

		response, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer response.Body.Close()
		if response.StatusCode != 200 {
			err = errors.New("status code is not 200")
		}

		scanner, _ := ResponseReader(response)
		itemDetail = UnmarshalCatalog(scanner)
		break
	}
	return itemDetail, err
}

func UnmarshalCatalog(responseRaw []byte) *ItemDetail {
	itemDetail := &ItemDetail{}

	err := json.Unmarshal(responseRaw, &itemDetail)
	if err != nil {
		fmt.Println(err)
	}
	return itemDetail
}

func ItemThumbnailImageById(assetId int) (string, error) {
	client := &http.Client{Timeout: 3 * time.Second}
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

	scanner, _ := ResponseReader(response)
	err = json.Unmarshal(scanner, &data)
	if err != nil {
		return "", err
	}

	return data.Detail[0].Image, nil
}

func main() {
	LoadConfig()
	ProxyTester()
	AddToWatcher()
}
