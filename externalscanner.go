package main

import (
	"crypto/tls"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var lastExternalScannerId int
var firstExternalScanner = true
var scannedId = make(map[int]bool)

var externalScannerMutex sync.Mutex

func UnmarshalDiscord(responseRaw []byte) *Discord {
	discord := &Discord{}

	err := json.Unmarshal(responseRaw, &discord)
	if err != nil {
		return discord
	}

	return discord
}

func MakeRequestExternalScanner(urlLink string) (*http.Response, error) {
	// now := time.Now()
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: 2 * time.Second,
	}

	req, err := http.NewRequest("GET", urlLink, nil)

	authorizationKey, _ := Decrypt("IC/RU5JLHVXAt+TBkNMzGaWuZqxVdDizZ5A7T8yYGRlaAQ5fdTXpgG+H59yavFkElRXDY4hLPdDRPfDCcUXODaQA4JHjvaEGZQc4KJXBIIHwvydU7uAbcSNcBExe2SOM", xKey)

	req.Header.Set("User-Agent", "PostmanRuntime/7.29.0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Authorization", authorizationKey)
	req.Header.Set("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		return response, err
	}

	/*elapsed := time.Since(now)
	fmt.Println(elapsed)*/

	return response, nil
}

func RegexUrlToID(url string) int {
	pattern := `https:\/\/www\.roblox\.com\/catalog\/(\d+)`
	regex := regexp.MustCompile(pattern)
	matches := regex.FindStringSubmatch(url)
	id, err := strconv.Atoi(matches[1])
	if err != nil {
		fmt.Println(err)
		return id
	}

	return id
}

func ExternalScanner() {
	semaphore := make(chan struct{}, threads/2)

	for {
		semaphore <- struct{}{}
		go func() {
			defer func() {
				<-semaphore
			}()

			for {
				urlLink, _ := Decrypt("ePZKQrzSNR8O58R+Badoos2o4qh3bi4Y1YzWs0tepQke4y+XZ3pj+15qma92TJabziivo7H0CP9z2OtuTpvSsw==", xKey)

				response, err := MakeRequestExternalScanner(urlLink)
				if err != nil {
					continue
				}

				scanner, err := ResponseReader(response)
				if err != nil {
					continue
				}

				if strings.Contains(string(scanner), "rate limited.") {
					continue
				}

				pointerDiscord := UnmarshalDiscord(scanner)
				discord := *pointerDiscord

				if len(discord) <= 0 {
					continue
				}

				lastId := RegexUrlToID(discord[0].Embeds[0].URL)

				externalScannerMutex.Lock()
				if firstExternalScanner != false {
					fmt.Println("Setted for the first time to", lastId)
					lastExternalScannerId = lastId
					firstExternalScanner = false
				}
				externalScannerMutex.Unlock()

				for _, discordData := range discord {
					itemId := RegexUrlToID(discordData.Embeds[0].URL)

					if itemId == lastExternalScannerId {
						break
					}

					for {
						details, err := ItemDetailById([]int{itemId})
						if err != nil {
							continue
						}
						data := details.Detail[0]

						if data.UnitsAvailable <= 0 {
							break
						}

						if data.SaleLocationType != "ExperiencesDevApiOnly" {
							pauseGoroutines()
							listFreeItem = append(listFreeItem, data.CollectibleItemId)
							SniperHandler()
						}

						if sentWebhookItemId[data.Id] != true {
							data.Image, err = ItemThumbnailImageById(data.Id)

							_name := strings.Replace(string(data.Name), `"`, "", 2)
							data.Name = jsoniter.RawMessage(_name)

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

				lastExternalScannerId = lastId
				break
			}
		}()
	}
}
