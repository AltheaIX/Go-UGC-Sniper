package main

import (
	"bytes"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"
)

var watcherId []int
var newItemId []int
var sentWebhookItemId = make(map[int]bool)

var notificationMutex sync.Mutex
var watcherMutex sync.Mutex
var pauseChan = make(chan struct{})
var pauseFlag bool

var freeWebhookUrl string
var paidWebhookUrl string
var threads int

func pauseGoroutines() {
	pauseChan <- struct{}{}
}

func resumeGoroutines() {
	watcherMutex.Lock()
	pauseFlag = false
	watcherMutex.Unlock()
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
	semaphore := make(chan struct{}, threads)
	go OffsaleTrackerHandler()
	for {
		select {
		case <-sig:
			fmt.Println("Termination signal received.")
			return
		case <-pauseChan:
			watcherMutex.Lock()
			isPaused := pauseFlag
			watcherMutex.Unlock()
			if isPaused {
				ReleaseSemaphore(semaphore)
				<-pauseChan // Wait for resume signal before continuing
				continue
			}
		default:
			semaphore <- struct{}{}
			go func(semaphore chan struct{}, newItemId []int) {
				defer func() {
					<-semaphore
				}()

				watcherMutex.Lock()
				if pauseFlag {
					watcherMutex.Unlock()
					return
				}
				watcherMutex.Unlock()

				lastIdFromArray, proxyURL, err := ItemRecentlyAddedAppend(ItemRecentlyAdded())

				if err != nil {
					if strings.Contains(err.Error(), "context deadline exceeded") {
						fmt.Printf("%v - Proxy timeout.\n", proxyURL)
						return
					}

					if strings.Contains(err.Error(), "status code is not 200") {
						fmt.Printf("%v - Rate limited.\n", proxyURL)
						return
					}

					if strings.ContainsAny(err.Error(), "An existing connection was forcibly closed by the remote host.") {
						fmt.Printf("%v - Proxy issues\n", proxyURL)
						return
					}

					fmt.Printf("%v - %v\n", proxyURL, err.Error())
					return
				}

				if lastItemId == lastIdFromArray {
					fmt.Printf("%v - No news items yet.\n", proxyURL)
					return
				}

				var idsToAdd []int
				for _, data := range listId {
					if data == lastItemId {
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
			fmt.Println(err)
		}
		req.Header.Set("Content-Type", "application/json")

		response, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
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
			fmt.Println(err)
		}
		req.Header.Set("Content-Type", "application/json")

		response, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
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

func OffsaleTracker(offsaleId []int, wg *sync.WaitGroup, semaphore chan struct{}) {
	defer wg.Done()

	for i := 0; i < len(offsaleId); i++ {
		select {
		case <-pauseChan:
			watcherMutex.Lock()
			isPaused := pauseFlag
			watcherMutex.Unlock()
			if isPaused {
				ReleaseSemaphore(semaphore)
				<-pauseChan // Wait for resume signal before continuing
				continue
			}
		default:
			semaphore <- struct{}{}

			go func(data []int) {
				defer func() {
					<-semaphore
				}()

				defer func() {
					if err := recover(); err != nil {
						fmt.Printf("Recovered from panic: %v\n", err)
					}
				}()

				watcherMutex.Lock()
				if pauseFlag {
					return
				}
				watcherMutex.Unlock()

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
							fmt.Printf("Watcher - %d still on offsale.\n", data.Id)
							continue
						}

						if data.UnitsAvailable == 0 || data.Quantity == 0 {
							fmt.Printf("Watcher - %d will be removed from watcher list.\n", data.Id)
							watcherMutex.Lock()
							watcherId = DeleteSlice(watcherId, data.Id)
							watcherMutex.Unlock()
							break
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

							break
						}

						listFreeItem = append(listFreeItem, data.CollectibleItemId)
						watcherMutex.Unlock()

						pauseGoroutines()
						go SniperHandler()

						if sentWebhookItemId[data.Id] != true {
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
			}(offsaleId)
		}
	}
}

func OffsaleTrackerHandler() {
	fmt.Println("System - Offsale Tracker activated.")
	wg := sync.WaitGroup{}
	semaphore := make(chan struct{}, 10)

	for {
		wg.Add(1)

		go OffsaleTracker(watcherId, &wg, semaphore)

		wg.Wait()
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
	for _, data := range newItemId {
		for {
			defer func(data int) {
				if err := recover(); err != nil {
					fmt.Printf("Recovered from panic: %v\n", err)
				}
			}(data)

			if sentWebhookItemId[data] != false {
				break
			}

			detail, err := ItemDetailById(data)
			if err != nil {
				continue
			}

			if detail.Detail[0].PriceStatus == "Off Sale" && detail.Detail[0].Quantity == 0 {
				watcherId = append(watcherId, data)
				break
			}

			if detail.Detail[0].Quantity == 0 {
				break
			}

			_name := strings.Replace(string(detail.Detail[0].Name), `"`, "", 2)
			detail.Detail[0].Name = jsoniter.RawMessage(_name)

			if detail.Detail[0].Price != 0 {
				detail.Detail[0].Image, err = ItemThumbnailImageById(data)
				if err != nil {
					fmt.Println(err)
					continue
				}

				go NotifierWatcher("paid", detail.Detail[0])
				fmt.Printf("Notifier - Webhook sent to for %d \n", data)
				break
			}

			pauseGoroutines()
			listFreeItem = append(listFreeItem, detail.Detail[0].CollectibleItemId)
			go SniperHandler()

			detail.Detail[0].Image, err = ItemThumbnailImageById(data)
			if err != nil {
				fmt.Println(err)
				continue
			}

			go NotifierWatcher("free", detail.Detail[0])
			fmt.Printf("Notifier - Webhook sent to for %d \n", data)
			break
		}
	}
}

func BoughtNotifier(name string) {
	client := &http.Client{Timeout: 5 * time.Second}

	webhookBuilder := fmt.Sprintf(`{
	  "content": "Bought **%v** on **%v**",
	  "embeds": null,
	  "attachments": []
	}`, name, username)
	dataRequest := bytes.NewBuffer([]byte(webhookBuilder))

	req, err := http.NewRequest("POST", "https://discord.com/api/webhooks/1110976926283214958/JJg0SEhpMT2xpt_g4LSfjPgiAqhYx2iiA88MlZ8t7aQuSnxELnaulhjDxvEJoV1w0o95", dataRequest)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		fmt.Println("Bought notifier error")
		return
	}
}
