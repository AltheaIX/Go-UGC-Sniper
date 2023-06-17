package main

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

var transport = &http.Transport{MaxIdleConns: 100, MaxIdleConnsPerHost: 100, DisableKeepAlives: false, IdleConnTimeout: 10 * time.Second, TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
var client = &http.Client{Timeout: 6 * time.Second, Transport: transport}

func GetCsrfTokenProxied(proxyURL *url.URL) string {
	var token string

	proxiedTransport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	proxiedClient := &http.Client{Timeout: 6 * time.Second, Transport: proxiedTransport}

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
	req.Header.Set("x-csrf-token", "s/BoATPEOLnW")

	response, err := proxiedClient.Do(req)
	if err != nil {
		//fmt.Printf("Error: %v - GetCsrfToken\n", err)
		return token
	}
	defer response.Body.Close()

	if response.Header.Get("x-csrf-token") != "" {
		token = response.Header.Get("x-csrf-token")
	}
	return token
}

func GetCsrfToken() string {
	var token string

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

func ItemRecentlyAdded(urlLink string) ([]byte, *url.URL, error) {
	proxyURL, err := url.Parse(BuildProxyURL(proxyList[rand.Intn(len(proxyList)-1)]))
	if err != nil {
		fmt.Printf("Error: %v - Parser Proxy\n", err)
		panic(err)
	}

	proxiedTransport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	proxiedClient := &http.Client{Timeout: 6 * time.Second, Transport: proxiedTransport}

	req, err := http.NewRequest("GET", urlLink, nil)

	req.Header.Set("User-Agent", "PostmanRuntime/7.29.0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json")

	response, err := proxiedClient.Do(req)
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

	proxyURL, err := url.Parse(BuildProxyURL(proxyList[rand.Intn(len(proxyList)-1)]))
	if err != nil {
		fmt.Printf("Error: %v - Parser Proxy\n", err)
		panic(err)
	}

	proxiedTransport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	proxiedClient := &http.Client{Timeout: 6 * time.Second, Transport: proxiedTransport}

	var items []OffsaleItems

	for _, data := range assetId {
		items = append(items, OffsaleItems{ItemType: 1, ID: data})
	}

	payload := &OffsalePayload{Items: items}
	jsonPayload, _ := json.Marshal(payload)
	dataRequest := bytes.NewBuffer(jsonPayload)
	fmt.Println(string(jsonPayload))

	req, err := http.NewRequest("POST", "https://catalog.roblox.com/v1/catalog/items/details", dataRequest)
	if err != nil {
		return itemDetail, err
	}

	csrf := GetCsrfTokenProxied(proxyURL)
	if err != nil {
		return itemDetail, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "PostmanRuntime/7.29.0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("x-csrf-token", csrf)
	fmt.Println(csrf)

	response, err := proxiedClient.Do(req)
	if err != nil {
		return itemDetail, err
	}
	defer response.Body.Close()

	scanner, _ := ResponseReader(response)
	if response.StatusCode != 200 {
		fmt.Println(string(scanner), response.StatusCode)
		err = errors.New("status code is not 200")
		return itemDetail, err
	}

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
		fmt.Print("error on line 238")
		return itemDetail, err
	}
	defer response.Body.Close()
	scanner, _ := ResponseReader(response)

	if response.StatusCode != 200 {
		/*
			pc, _, line, _ := runtime.Caller(1)
			callingFunc := runtime.FuncForPC(pc).Name()
			fmt.Println(string(scanner), callingFunc, line)*/
		fmt.Println(string(scanner))
		fmt.Print("error on line 250")
		err = errors.New("status code is not 200")
		return itemDetail, err
	}

	itemDetail, err = UnmarshalCatalog(scanner)
	if err != nil {
		fmt.Println(string(scanner))
		return itemDetail, err
	}

	return itemDetail, nil
}

func ItemThumbnailImageById(assetId int) (string, error) {
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
