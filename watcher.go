package main

import (
	"bytes"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"strings"
	"time"
)

var watcherId []int
var newItemId []int

func AddToWatcher() {
	for {
		lastIdFromArray, err := ItemRecentlyAddedAppend()
		newItemId = nil

		if err != nil {
			fmt.Println("Trying to recovery too many request...")
			time.Sleep(15 * time.Second)
			fmt.Println("Sleep done, lets see if it works...")
			continue
		}

		if lastItemId == lastIdFromArray {
			fmt.Println("No news items yet...")
			fmt.Println(newItemId)
			continue
		}

		for _, data := range listId {
			if data == lastItemId {
				break
			}

			watcherId = append(watcherId, data)
			newItemId = append(newItemId, data)
		}
		fmt.Println(watcherId)

		if newItemId != nil {
			go NotifierWatcherHandle(newItemId)
		}

		lastItemId = lastIdFromArray
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
		  "username": "Free Item Notifier",
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
		  "username": "Paid Item Notifier",
		  "attachments": []
		}`, data.Name, data.Id, data.Price, data.Quantity, data.Id, data.Image)
		dataRequest := bytes.NewBuffer([]byte(webhookBuilder))

		req, err := http.NewRequest("POST", "https://discord.com/api/webhooks/1098275294277816341/obz7AD8ju43-89LRoMVI5YzvN_YL03v1DJvTBCOSX-RbxtBg2Qrg5ZZLddzev1U-ZnBh", dataRequest)
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

func NotifierWatcherHandle(newItemId []int) {
	for _, data := range newItemId {
		for {
			detail, err := ItemDetailById(data)
			if err != nil {
				fmt.Println(err)
				time.Sleep(15 * time.Second)
				continue
			}

			detail.Detail[0].Image, err = ItemThumbnailImageById(data)
			_name := strings.Replace(string(detail.Detail[0].Name), `"`, "", 2)
			detail.Detail[0].Name = jsoniter.RawMessage(_name)
			if err != nil {
				fmt.Println(err)
				continue
			}

			if detail.Detail[0].Price != 0 {
				go NotifierWatcher("paid", detail.Detail[0])
				fmt.Printf("Webhook sent to for %d", data)
				break
			}

			go NotifierWatcher("free", detail.Detail[0])
			fmt.Printf("Webhook sent to for %d", data)
			break
		}
	}
}
