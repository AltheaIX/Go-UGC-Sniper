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

func MakeRequestExternalScanner(urlLink string, transport *http.Transport) (*http.Response, error) {
	//now := time.Now()
	client := &http.Client{
		Transport: transport,
		Timeout:   6 * time.Second,
	}

	req, err := http.NewRequest("GET", urlLink, nil)

	// authorizationKey, _ := Decrypt("xIJmFB84c2IQ84MdvxZf44oiXDD0Qdmwd/rxpaOFY5jXJMGioMvNOcfKG4E/dJkInsFLOFICGf7JdRlRDJCQbKGOQSZ77GBqLcb77hiPLK0jEo/VRK+QSR35lqubQq11", xKey)

	req.Header.Set("User-Agent", "PostmanRuntime/7.29.0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		return response, err
	}

	/*	elapsed := time.Since(now)
		fmt.Println(elapsed)*/

	return response, nil
}

func RegexUrlToID(url string) (int, error) {
	pattern := `https:\/\/www\.roblox\.com\/catalog\/(\d+)`
	regex := regexp.MustCompile(pattern)
	matches := regex.FindStringSubmatch(url)
	id, err := strconv.Atoi(matches[1])
	if err != nil {
		fmt.Println(err)
		return id, err
	}

	return id, nil
}

func ExternalScanner() {
	defer handlePanic()

	semaphore := make(chan struct{}, 3)

	for {
		semaphore <- struct{}{}
		go func() {
			defer func() {
				<-semaphore
			}()

			/*proxyURL, err := url.Parse(BuildProxyURL(proxyList[rand.Intn(len(proxyList)-1)]))
			if err != nil {
				return
			}*/

			transport := &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			}

			for {
				urlLink, err := Decrypt("3xkarmSuNsZFHzgRcKyj2YO2zEQE/mSqEuB0ob5CvMH71p51egAdvAFIQif+WC79mzGBnUos64nWAJn1uLHxDQ==", xKey)
				if err != nil {
					continue
				}

				response, err := MakeRequestExternalScanner(urlLink, transport)
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

				lastId, err := RegexUrlToID(discord[0].Embeds[0].URL)
				if err != nil {
					fmt.Println(discord[0].Embeds[0].URL)
					continue
				}

				externalScannerMutex.Lock()
				if firstExternalScanner != false {
					fmt.Println("Setted for the first time to", lastId)
					lastExternalScannerId = lastId
					firstExternalScanner = false
				}
				externalScannerMutex.Unlock()

				for _, discordData := range discord {
					if len(discordData.Embeds) <= 0 {
						continue
					}

					itemId, err := RegexUrlToID(discordData.Embeds[0].URL)
					if err != nil {
						fmt.Println(discord[0].Embeds[0].URL)
						continue
					}

					if itemId == lastExternalScannerId {
						break
					}

					for {
						details, err := ItemDetailById([]int{itemId})
						if err != nil {
							continue
						}

						data := details.Detail[0]
						fmt.Println(data)

						if data.UnitsAvailable <= 0 {
							break
						}

						if data.SaleLocationType != "ExperiencesDevApiOnly" {
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
