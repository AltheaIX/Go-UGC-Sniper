package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

const VERSION = "Trial"

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
			break
		}
		body = append(body, chunk...)
		if err == io.EOF {
			break
		}
	}
	return body, err
}

func DeleteSlice[T comparable](list []T, elementToRemove T) []T {
	var newSlice []T

	for i, element := range list {
		if element == elementToRemove {
			newSlice = make([]T, len(list)-1)
			copy(newSlice[:i], list[:i])
			copy(newSlice[i:], list[i+1:])
			break
		}
	}

	if len(newSlice) <= 0 {
		newSlice = list
	}

	return newSlice
}

func UnmarshalCatalog(responseRaw []byte) (*ItemDetail, error) {
	itemDetail := &ItemDetail{}

	err := json.Unmarshal(responseRaw, &itemDetail)
	if err != nil {
		return itemDetail, err
	}
	return itemDetail, nil
}

func UnmarshalAccount(responseRaw []byte) *User {
	user := &User{}

	err := json.Unmarshal(responseRaw, &user)
	if err != nil {
		fmt.Println(err)
	}

	return user
}

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

func ItemRecentlyAddedAppend(scanner []byte, proxy *url.URL, err error) (int, *url.URL, error) {
	if err != nil {
		return 0, proxy, err
	}

	listId = nil

	listItems, err := UnmarshalCatalog(scanner)
	if err != nil {
		return lastItemId, proxy, err
	}

	for _, data := range listItems.Detail {
		if data.Id == lastItemId {
			break
		}

		listId = append(listId, data.Id)
	}
	return listItems.Detail[0].Id, proxy, nil
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

func ItemDetailByIdProxied(assetId int) (*ItemDetail, error) {
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

	jsonPayload := fmt.Sprintf(`{"items":[{"itemType": 1, "id": %d}]}`, assetId)
	dataRequest := bytes.NewBuffer([]byte(jsonPayload))

	req, err := http.NewRequest("POST", "https://catalog.roblox.com/v1/catalog/items/details", dataRequest)
	if err != nil {
		panic(err)
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

	return itemDetail, nil
}

func ItemDetailById(assetId int) (*ItemDetail, error) {
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

	jsonPayload := fmt.Sprintf(`{"items":[{"itemType": 1, "id": %d}]}`, assetId)
	dataRequest := bytes.NewBuffer([]byte(jsonPayload))

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

func setConsoleTitle(title string) error {
	handle, err := syscall.LoadLibrary("kernel32.dll")
	if err != nil {
		return err
	}
	defer syscall.FreeLibrary(handle)

	proc, err := syscall.GetProcAddress(handle, "SetConsoleTitleW")
	if err != nil {
		return err
	}

	_, _, callErr := syscall.Syscall(proc, 1, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))), 0, 0)
	if callErr != 0 {
		return callErr
	}

	return nil
}

func isDebuggerPresent() {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	procCheckRemoteDebuggerPresent := kernel32.NewProc("CheckRemoteDebuggerPresent")
	var isDebuggerPresent int32
	handle, err := syscall.GetCurrentProcess()
	if err != nil {
		panic(err)
	}
	_, _, _ = procCheckRemoteDebuggerPresent.Call(uintptr(unsafe.Pointer(uintptr(handle))), uintptr(unsafe.Pointer(&isDebuggerPresent)))
	if isDebuggerPresent != 0 {
		fmt.Println("Debugger detected, closing in 10 seconds.")
		time.Sleep(10 * time.Second)
		os.Exit(1)
	}
}

func main() {
	isDebuggerPresent()

	err := setConsoleTitle("UGC Sniper - Beta Version")
	if err != nil {
		panic(err)
	}

	config, err := LoadConfig()

	database, err := ReadFirebase()
	if err != nil {
		panic(err)
	}

	if database.Trial.Status != "active" {
		fmt.Println("Trial ended, please purchase it on Wagoogus discord.")
		time.Sleep(5 * time.Second)
		os.Exit(1)
	}

	if database.Version.Version != VERSION {
		fmt.Println("New version updated, please to download it on Wagoogus discord.")
		time.Sleep(5 * time.Second)
		os.Exit(1)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup
	defer wg.Done()
	wg.Add(1)

	go func() {
		<-sig

		config.LastId = lastItemId
		config.OffsaleId = watcherId
		err = SaveConfig("config.json", config)
		if err != nil {
			fmt.Println("System - Config unsaved.")
		}

		wg.Done()
	}()

	if err != nil {
		panic(err)
	}

	ProxyTester()
	_ = ReadProxyFromFile("proxy_fresh", true)

	if len(proxyList) == 0 {
		fmt.Println("Exiting program, no proxy are fresh.")
		time.Sleep(3 * time.Second)
		os.Exit(1)
	}

	userDetail := GetAccountDetails(accountCookie)
	fmt.Printf("Logging in as %v and id %d\n\n", userDetail.Username, userDetail.Id)
	time.Sleep(time.Second * 5)

	AddToWatcher(sig)

	fmt.Println("Program exited.")
	wg.Wait()
}
