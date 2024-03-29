package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"os"
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"
)

var watcherId []int
var newItemId []int
var sentWebhookItemId = make(map[int]bool)

var endpointItem = []string{"https://catalog.roblox.com/v1/search/items?category=Accessories&includeNotForSale=true&keyword=orange+teal+cyan+red+green+topaz+yellow+purple+war&limit=120&salesTypeFilter=1&sortType=3&subcategory=Accessories"}

var maxWatcherSize = 120

var notificationMutex sync.Mutex
var watcherMutex sync.Mutex
var pauseChan = make(chan struct{})

var freeWebhookUrl string
var paidWebhookUrl string
var threads int

func resumeGoroutines() {
	globalCtx, globalCancel = context.WithCancel(context.Background())
}

func ReleaseSemaphore(semaphore chan struct{}) {
	for {
		select {
		case <-semaphore:
		default:
			return
		}
	}
}

func AddToWatcher(sig chan os.Signal) {
	defer handlePanic()

	semaphore := make(chan struct{}, threads)
	go OffsaleTrackerHandler()
	go ExternalScanner()
	for {
		select {
		case <-sig:
			fmt.Println("Termination signal received.")
			return
		case <-globalCtx.Done():
			return
		default:
			semaphore <- struct{}{}
			go func(semaphore chan struct{}, newItemId []int) {
				defer func() {
					<-semaphore
				}()

				lastIdFromArray, proxyURL, err := ItemRecentlyAddedAppend(ItemRecentlyAdded("https://catalog.roblox.com/v1/search/items?category=Accessories&includeNotForSale=true&keyword=orange+teal+cyan+red+green+topaz+yellow+purple+war&limit=120&salesTypeFilter=1&sortType=3&subcategory=Accessories"))

				if err != nil {
					if strings.Contains(err.Error(), "context deadline exceeded") {
						// fmt.Printf("%v - Proxy timeout.\n", proxyURL)
						return
					}

					if strings.Contains(err.Error(), "status code is not 200") {
						// fmt.Printf("%v - Rate limited.\n", proxyURL)
						return
					}

					if strings.ContainsAny(err.Error(), "An existing connection was forcibly closed by the remote host.") {
						// fmt.Printf("%v - Proxy issues\n", proxyURL)
						return
					}

					fmt.Printf("%v - %v\n", proxyURL, err.Error())
					return
				}

				if lastItemId == lastIdFromArray {
					// fmt.Printf("%v - No news items yet.\n", proxyURL)
					return
				}

				var idsToAdd []int
				for _, data := range listId {
					if data <= lastItemId {
						break
					}

					idsToAdd = append(idsToAdd, data)
				}

				if len(idsToAdd) > 0 {
					notificationMutex.Lock()
					defer notificationMutex.Unlock()
					go NotifierWatcherHandler(idsToAdd)
				}

				lastItemId = lastIdFromArray
			}(semaphore, newItemId[:])
			newItemId = nil
		}
	}
}

func NotifierWatcher(notifierType string, data Data) error {
	defer handlePanic()

	if sentWebhookItemId[data.Id] != false {
		return errors.New("webhook tried to send more than once")
	}

	sentWebhookItemId[data.Id] = true

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
				  "value": "%d"
				},
				{
				  "name": "Available Copy",
				  "value": "%d"
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
		  "username": "Free Item Notifier",
		  "attachments": []
		}`, data.Name, data.Id, data.Price, data.Quantity, data.UnitsAvailable, data.Id, data.Image)
		dataRequest := bytes.NewBuffer([]byte(webhookBuilder))

		req, err := http.NewRequest("POST", freeWebhookUrl, dataRequest)
		if err != nil {
			panic(err)
		}
		req.Header.Set("Content-Type", "application/json")

		response, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer response.Body.Close()

		if response.StatusCode != 200 {
			return errors.New("webhook string code is not 200")
		}
		break
	case "paid":
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
				  "value": "%d"
				},
				{
				  "name": "Unit Available",
				  "value": "%d"
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
		  "username": "Paid Item Notifier",
		  "attachments": []
		}`, data.Name, data.Id, data.Price, data.Quantity, data.UnitsAvailable, data.Id, data.Image)
		dataRequest := bytes.NewBuffer([]byte(webhookBuilder))

		req, err := http.NewRequest("POST", paidWebhookUrl, dataRequest)
		if err != nil {
			panic(err)
		}
		req.Header.Set("Content-Type", "application/json")

		response, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer response.Body.Close()
		if response.StatusCode != 200 {
			return errors.New("webhook string code is not 200")
		}

		fmt.Println(err)
		break
	}
	return nil
}

func OffsaleTracker(offsaleId []int, semaphore chan struct{}) {
	defer handlePanic()

	select {
	case <-globalCtx.Done():
		return
	default:
		go func(data []int) {
			defer func() {
				<-semaphore
			}()

			defer func() {
				if err := recover(); err != nil {
					fmt.Printf("Recovered from panic: %v\n", err)
				}
			}()

			for {
				details, err := ItemDetailByIdProxied(offsaleId)
				if err != nil {
					continue
				}

				for _, data := range details.Detail {
					watcherMutex.Lock()

					if sentWebhookItemId[data.Id] != false {
						watcherMutex.Unlock()
						continue
					}

					watcherMutex.Unlock()

					if data.PriceStatus == "Off Sale" && data.Quantity == 0 {
						//fmt.Printf("Watcher - %d is still on offsale.\n", data.Id)
						continue
					}

					if data.PriceStatus != "Off Sale" && (data.UnitsAvailable == 0 || data.Quantity == 0) {
						//fmt.Printf("Watcher - %d will be removed from watcher list.\n", data.Id)
						watcherMutex.Lock()
						watcherId = DeleteSlice(watcherId, data.Id)
						watcherMutex.Unlock()
						continue
					}

					_name := strings.Replace(string(data.Name), `"`, "", 2)
					data.Name = jsoniter.RawMessage(_name)

					watcherMutex.Lock()
					watcherId = DeleteSlice(watcherId, data.Id)

					if data.Price != 0 {
						watcherMutex.Unlock()

						if sentWebhookItemId[data.Id] != true {
							data.Image, err = ItemThumbnailImageById(data.Id)
							if err != nil {
								fmt.Println(err)
								continue
							}

							go NotifierWatcher("paid", data)
							fmt.Printf("Notifier - Webhook sent to for %d \n", data.Id)
						}

						continue
					}
					watcherMutex.Unlock()

					if data.SaleLocationType != "ExperiencesDevApiOnly" {
						listFreeItem = append(listFreeItem, data.CollectibleItemId)
						SniperHandler("Tracker")
					}

					if sentWebhookItemId[data.Id] != true {
						data.Image, err = ItemThumbnailImageById(data.Id)
						if err != nil {
							fmt.Println(err)
							continue
						}
						go NotifierWatcher("free", data)
						fmt.Printf("Notifier - Webhook sent to for %d \n", data.Id)
					}

					continue
				}
				break
			}
		}(offsaleId)
	}
}

func OffsaleTrackerHandler() {
	fmt.Println("System - Offsale Tracker activated.")
	semaphore := make(chan struct{}, threads/2)

	for {
		size := len(watcherId) / maxWatcherSize

		for i := 0; i < size; i++ {
			semaphore <- struct{}{}
			offsaleId := watcherId[(i * maxWatcherSize) : (i+1)*maxWatcherSize]
			go OffsaleTracker(offsaleId, semaphore)
		}

		if len(watcherId) > maxWatcherSize*size {
			offsaleId := watcherId[maxWatcherSize*size:]
			semaphore <- struct{}{}
			go OffsaleTracker(offsaleId, semaphore)
		}
	}
}

func IsFieldSet(s interface{}, fieldName string) bool {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		return false
	}

	zeroValue := reflect.Zero(field.Type())

	return !reflect.DeepEqual(field.Interface(), zeroValue.Interface())
}

func NotifierWatcherHandler(newItemId []int) {
	defer handlePanic()

	for {
		details, err := ItemDetailById(newItemId)
		if err != nil {
			fmt.Println(err, "Notifier")
			continue
		}

		for _, data := range details.Detail {
			if sentWebhookItemId[data.Id] != false {
				continue
			}

			if data.PriceStatus == "Off Sale" && data.Quantity == 0 {
				if IsExist(watcherId, data.Id) {
					continue
				}

				if len(watcherId) == 120 {
					watcherId = watcherId[:118]
				}

				watcherId = append(watcherId, data.Id)
				sort.Sort(sort.Reverse(sort.IntSlice(watcherId)))
				continue
			}

			if data.Quantity == 0 {
				continue
			}

			_name := strings.Replace(string(data.Name), `"`, "", 2)
			data.Name = jsoniter.RawMessage(_name)

			if data.Price != 0 {
				data.Image, err = ItemThumbnailImageById(data.Id)
				if err != nil {
					fmt.Println(err)
					continue
				}

				go NotifierWatcher("paid", data)
				fmt.Printf("Notifier - Webhook sent to for %d \n", data.Id)
				continue
			}

			if data.SaleLocationType != "ExperiencesDevApiOnly" {
				listFreeItem = append(listFreeItem, data.CollectibleItemId)
				SniperHandler("Watcher")
			}

			data.Image, err = ItemThumbnailImageById(data.Id)
			if err != nil {
				fmt.Println(err)
				continue
			}

			go NotifierWatcher("free", data)
			fmt.Printf("Notifier - Webhook sent to for %d \n", data.Id)
		}
		break
	}
}

func BoughtNotifier(name jsoniter.RawMessage, source string) {
	client := &http.Client{Timeout: 5 * time.Second}

	webhookBuilder := fmt.Sprintf(`{
	  "content": "Bought **%s** on **%v** - From %v",
	  "embeds": null,
	  "attachments": []
	}`, name, username, source)
	dataRequest := bytes.NewBuffer([]byte(webhookBuilder))

	fmt.Println(webhookBuilder)

	webhookURL, _ := Decrypt("i/LOatue4KyPz9MRDB61XW9BIez/ZMyRD2/EbR0oOPWt7dVA1Jg5R5UKy02vEJotBbb4p6ohzEVjf0AD+SFhrS4RWldSzpH3dlABnVzKpBNtDpvCPKl/4/fTP2sKlyOFTEUUV74vgaab8FjJsKwXeV4PJOhSIoJFreB3hLSIQZRNBE75mM1oLvGTsWrm8Ll9", xKey)

	req, err := http.NewRequest("POST", webhookURL, dataRequest)
	if err != nil {
		fmt.Println("line 490", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		fmt.Println("line 497", err)
		return
	}

	scanner, _ := ResponseReader(response)
	fmt.Println(string(scanner))

	defer response.Body.Close()
}
