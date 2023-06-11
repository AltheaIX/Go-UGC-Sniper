package main

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func GetCsrfTokenProxied(proxyURL *url.URL) string {
	var token string

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: 6 * time.Second,
	}

	jsonRequest := fmt.Sprintf(`{"items":[{"itemType": 1, "id": %x}]}`, 13177094956)
	dataRequest := bytes.NewBuffer([]byte(jsonRequest))

	req, err := http.NewRequest("POST", "https://catalog.roblox.com/v1/catalog/items/details", dataRequest)
	if err != nil {
		// fmt.Printf("Error: %v - GetCsrfToken\n", err)
		return token
	}

	req.Header.Set("User-Agent", "PostmanRuntime/7.29.0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-csrf-token", "xcsrf")

	response, err := client.Do(req)
	if err != nil {
		//fmt.Printf("Error: %v - GetCsrfToken\n", err)
		return token
	}
	defer response.Body.Close()

	if response.Header.Get("x-csrf-token") != "" {
		token = response.Header.Get("x-csrf-token")
	}

	time.Sleep(0 * time.Second)
	return token
}

func GetCsrfToken() string {
	var token string
	client := &http.Client{Timeout: 3 * time.Second}

	cookie := &http.Cookie{
		Name:    ".ROBLOSECURITY",
		Value:   accountCookie,
		Path:    "/",
		Domain:  "roblox.com",
		Expires: time.Now().Add(time.Hour * 1000),
	}

	jsonRequest := fmt.Sprintf(`{"items":[{"itemType": 1, "id": %x}]}`, 13177094956)
	dataRequest := bytes.NewBuffer([]byte(jsonRequest))

	req, err := http.NewRequest("POST", "https://catalog.roblox.com/v1/catalog/items/details", dataRequest)
	if err != nil {
		// fmt.Printf("Error: %v - GetCsrfToken\n", err)
		return token
	}

	req.AddCookie(cookie)
	req.Header.Set("User-Agent", "PostmanRuntime/7.29.0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-csrf-token", "xcsrf")

	response, err := client.Do(req)
	if err != nil {
		//fmt.Printf("Error: %v - GetCsrfToken\n", err)
		return token
	}
	defer response.Body.Close()

	if response.Header.Get("x-csrf-token") != "" {
		token = response.Header.Get("x-csrf-token")
	}

	time.Sleep(0 * time.Second)
	return token
}

func ItemRecentlyAdded() ([]byte, *url.URL, error) {
	proxyURL, err := url.Parse(strings.TrimRight("http://"+proxyList[rand.Intn(len(proxyList)-1)], "\x00"))
	if err != nil {
		fmt.Printf("Error: %v - Parser Proxy\n", err)
		panic(err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: 6 * time.Second,
	}

	req, err := http.NewRequest("GET", "https://catalog.roblox.com/v1/search/items?category=Accessories&includeNotForSale=true&limit=120&salesTypeFilter=1&sortType=3&subcategory=Accessories", nil)

	req.Header.Set("User-Agent", "PostmanRuntime/7.29.0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		return nil, proxyURL, err
	}

	defer response.Body.Close()
	if response.StatusCode != 200 {
		return nil, proxyURL, errors.New("status code is not 200")
	}

	scanner, _ := ResponseReader(response)

	if string(scanner) == "" {
		return nil, proxyURL, errors.New("empty body")
	}

	return scanner, proxyURL, nil
}

func ItemDetailByIdProxied(assetId []int) (*ItemDetail, error) {
	itemDetail := &ItemDetail{}
	var err error

	proxyURL, err := url.Parse(strings.TrimRight("http://"+proxyList[rand.Intn(len(proxyList)-1)], "\x00"))
	if err != nil {
		fmt.Printf("Error: %v - Parser Proxy\n", err)
		panic(err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: 6 * time.Second,
	}

	var items []OffsaleItems

	for _, data := range assetId {
		items = append(items, OffsaleItems{ItemType: 1, ID: data})
	}

	payload := &OffsalePayload{Items: items}
	jsonPayload, _ := json.Marshal(payload)
	dataRequest := bytes.NewBuffer(jsonPayload)

	req, err := http.NewRequest("POST", "https://catalog.roblox.com/v1/catalog/items/details", dataRequest)
	if err != nil {
		return itemDetail, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "PostmanRuntime/7.29.0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("x-csrf-token", GetCsrfTokenProxied(proxyURL))

	response, err := client.Do(req)
	if err != nil {
		return itemDetail, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		err = errors.New("status code is not 200")
		return itemDetail, err
	}

	scanner, _ := ResponseReader(response)
	itemDetail, err = UnmarshalCatalog(scanner)
	if err != nil {
		return itemDetail, err
	}

	if string(scanner) == "" {
		return itemDetail, errors.New("empty body")
	}

	return itemDetail, nil
}

func ItemDetailById(assetId []int) (*ItemDetail, error) {
	itemDetail := &ItemDetail{}
	var err error

	client := &http.Client{Timeout: 3 * time.Second}

	cookie := &http.Cookie{
		Name:    ".ROBLOSECURITY",
		Path:    "/",
		Value:   accountCookie,
		Domain:  "roblox.com",
		Expires: time.Now().Add(time.Hour * 1000),
	}

	var items []OffsaleItems

	for _, data := range assetId {
		items = append(items, OffsaleItems{ItemType: 1, ID: data})
	}

	payload := &OffsalePayload{Items: items}
	jsonPayload, _ := json.Marshal(payload)
	dataRequest := bytes.NewBuffer(jsonPayload)

	req, err := http.NewRequest("POST", "https://catalog.roblox.com/v1/catalog/items/details", dataRequest)
	if err != nil {
		panic(err)
	}
	req.AddCookie(cookie)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "PostmanRuntime/7.29.0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("x-csrf-token", GetCsrfToken())

	response, err := client.Do(req)
	if err != nil {
		return itemDetail, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		err = errors.New("status code is not 200")
		fmt.Println("ItemDetail - Rate limit, Item notifier maybe delayed!")
		return itemDetail, err
	}

	scanner, _ := ResponseReader(response)

	itemDetail, err = UnmarshalCatalog(scanner)
	if err != nil {
		fmt.Println(string(scanner))
		return itemDetail, err
	}

	return itemDetail, nil
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

func MakeRequest(url string) (*http.Response, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: 6 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)

	req.Header.Set("User-Agent", "PostmanRuntime/7.29.0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		return response, err
	}

	return response, nil
}
